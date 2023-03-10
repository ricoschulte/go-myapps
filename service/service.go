package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	logger "github.com/chi-middleware/logrus-logger"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
)

type AppService struct {
	ListenIp      string
	ListenPort    int
	ListenPortTls int
	TlsCertString string
	TlsKeyString  string
	Domain        string
	Instance      string
	Name          string
	Password      string

	Fs         http.FileSystem
	ApiHandler []PbxApiInterface

	Connections       []*AppServicePbxConnection // list of current connected websocket connections
	ConnectionsMutext sync.Mutex
}

func NewAppService(ip string, port int, portTls int, tlsCert string, tlsCertKey string, domain, name, instance, password string, fS http.FileSystem) (*AppService, error) {
	if portTls != 0 && (tlsCert == "" || tlsCertKey == "") {
		return nil, fmt.Errorf("when Tls Port set, a Certificate has to be set too. generate one with 'openssl req -new -newkey rsa:2048 -days 365 -nodes -x509 -keyout server.key -out server.crt'")
	}

	return &AppService{
		ListenIp:      ip,
		ListenPort:    port,
		ListenPortTls: portTls,
		TlsCertString: tlsCert,    // X509 Certificate
		TlsKeyString:  tlsCertKey, // X509 Private Key
		Domain:        domain,
		Instance:      instance,
		Name:          name,
		Password:      password,
		Fs:            fS,
	}, nil
}

func (s *AppService) Start() error {
	rootpath := fmt.Sprintf("/%s/%s/%s/", s.Domain, s.Name, s.Instance)
	log.Infof("register service with path: %s", rootpath)

	mux := http.NewServeMux()

	router := chi.NewRouter()
	logs := log.New()

	router.Use(logger.Logger("router", logs))

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.SetHeader("Server", s.Name))

	mux.HandleFunc(fmt.Sprintf("/%s/%s/%s", s.Domain, s.Name, s.Instance), func(w http.ResponseWriter, r *http.Request) {
		handleConnectionForHttpOrWebsocket(s, w, r)
	})

	// the router for the Webpage
	router_api := chi.NewRouter()
	//GetHttpRoutes(s, router_api)
	router_api.Get("/*",
		http.StripPrefix(fmt.Sprintf("/%s/%s/%s/", s.Domain, s.Name, s.Instance),
			http.FileServer(s.Fs),
		).ServeHTTP,
	)

	router.Mount(rootpath, router_api)
	mux.Handle(rootpath, router)

	mux.HandleFunc("/manager/fixcert.htm", s.HandleManagerFixCertResponse)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.
			WithField("path", r.URL.Path).
			Info("404 not found")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	})

	go func() {
		if s.ListenPort > 0 {
			log.Infof("starting Http %s:%v", s.ListenIp, s.ListenPort)
			if err := http.ListenAndServe(fmt.Sprintf("%s:%v", s.ListenIp, s.ListenPort), mux); err != nil {
				log.Errorf("Error starting Http server: '%s:%d' %v ", s.ListenIp, s.ListenPort, err)
			}
		}
	}()

	go func() {
		if s.ListenPortTls > 0 {
			log.Infof("starting Http/Tls %s:%v", s.ListenIp, s.ListenPortTls)
			err := s.StartTls(mux)
			if err != nil {
				log.Errorf("Error starting Http/Tls server: '%s:%d' %v ", s.ListenIp, s.ListenPortTls, err)
			}
		}
	}()

	return nil
}

/*
openssl req -new -newkey rsa:2048 -days 365 -nodes -x509 -keyout server.key -out server.crt
*/
func (s *AppService) StartTls(mux *http.ServeMux) error {

	cert, err := tls.X509KeyPair([]byte(s.TlsCertString), []byte(s.TlsKeyString))
	if err != nil {
		log.Errorf("Error loading certificate:", err)
		return err
	}

	server := &http.Server{
		Addr: fmt.Sprintf("%s:%v", s.ListenIp, s.ListenPortTls),
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
		Handler: mux,
	}

	return server.ListenAndServeTLS("", "")
}

func (s *AppService) RegisterHandler(handler PbxApiInterface) error {
	s.ApiHandler = append(s.ApiHandler, handler)
	return nil
}

func (s *AppService) HandleApiConnected(connection *AppServicePbxConnection, msg []byte) {
	for _, apiName := range connection.PbxInfo.Apis {
		for _, handler := range s.ApiHandler {
			if handler.GetApiName() == apiName {
				handler.OnConnect(connection)
			}
		}
	}
}
func (s *AppService) HandleUserConnected(connection *AppServicePbxConnection, msg []byte) {
	for _, handler := range s.ApiHandler {
		if handler.GetApiName() == "user" {
			handler.OnConnect(connection)
		}
	}
}
func (s *AppService) HandleAdminConnected(connection *AppServicePbxConnection, msg []byte) {
	for _, handler := range s.ApiHandler {
		if handler.GetApiName() == "admin" {
			handler.OnConnect(connection)
		}
	}
}

func (s *AppService) HandleApiDisConnected(connection *AppServicePbxConnection) {
	for _, apiName := range connection.PbxInfo.Apis {
		for _, handler := range s.ApiHandler {
			if handler.GetApiName() == apiName {
				handler.OnDisconnect(connection)
			}
		}
	}

	log.Debug("on disconnect of ", connection.AppLogin.App)
	if strings.HasSuffix(connection.AppLogin.App, ".htm") {
		app := strings.TrimSuffix(connection.AppLogin.App, ".htm")
		for _, handler := range s.ApiHandler {
			if handler.GetApiName() == app {
				handler.OnDisconnect(connection)
			}
		}
	}
}

func (s *AppService) HandleApiMessage(connection *AppServicePbxConnection, apiName string, message []byte) {
	msg := BaseMessage{}
	if err := json.Unmarshal(message, &msg); err != nil {
		connection.log().Errorf("server: error unmarshalling message: %v", err)
	}
	for _, handler := range s.ApiHandler {
		if handler.GetApiName() == apiName {
			handler.HandleMessage(connection, &msg, message)
			return
		}
	}
	log.Warnf("no handler for api '%s'", apiName)
}
func (s *AppService) HandleManagerFixCertResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/html")

	body := `<!DOCTYPE html>
<html>
<head>
	<title></title>
	<script type="text/javascript">
		function loaded() {
			var caller = window.opener || window.parent;
			if (caller) caller.postMessage(JSON.stringify({ mt: "HostReachable", host: window.location.host }), "*");
			window.close();
		}
	</script>
</head>
<body onload="loaded()">
	You can close this page now.
</body>
</html>`
	w.Write([]byte(body))
}

func (s *AppService) AddConnection(connection *AppServicePbxConnection) {
	s.ConnectionsMutext.Lock()
	defer s.ConnectionsMutext.Unlock()
	s.Connections = append(s.Connections, connection)
}

func (s *AppService) DeleteConnection(connection *AppServicePbxConnection) {
	s.ConnectionsMutext.Lock()
	defer s.ConnectionsMutext.Unlock()

	for i, value := range s.Connections {
		if value == connection {
			// Remove the item from the slice by slicing it
			s.Connections = append(s.Connections[:i], s.Connections[i+1:]...)
			break
		}
	}
}

/*
send the message to all connected Users or Admins
*/
func (s *AppService) SendToAll(message []byte) {
	log.Debug("SendToAll")
	s.SendToAllAdmins(message)
	s.SendToAllUsers(message)
}

/*
send the message to all connected Users
*/
func (s *AppService) SendToAllUsers(message []byte) {
	log.Debug("SendToAllUsers")
	for _, connection := range s.Connections {
		log.Tracef("SendToAllUsers %s", connection.AppLogin.App)
		if connection.AppLogin.App == "user.htm" {
			log.Tracef("SendToAllUsers %s", string(message))
			connection.WriteMessage(message)
		}
	}
}

/*
send the message to all connected Admins
*/
func (s *AppService) SendToAllAdmins(message []byte) {
	log.Debug("SendToAllAdmins")
	for _, connection := range s.Connections {
		log.Tracef("SendToAllAdmins %s", connection.AppLogin.App)
		if connection.AppLogin.App == "admin.htm" {
			log.Tracef("SendToAllAdmins %s", string(message))
			connection.WriteMessage(message)
		}
	}
}

/*
send the message to all connected clients of a specific user

sip is a string of the sip/h323 name of the user
*/
func (s *AppService) SendToAllConnectionsOfSip(message []byte, sip string) {
	log.Debug("SendToAllConnectionsOfSip")

	for _, connection := range s.Connections {
		if connection.AppLogin.Sip == sip {
			log.Tracef("SendToAllConnectionsOfSip ", sip, string(message))
			connection.WriteMessage(message)
		}
	}
}

type PbxApiInterface interface {
	GetApiName() string
	OnConnect(connection *AppServicePbxConnection)
	OnDisconnect(connection *AppServicePbxConnection)
	HandleMessage(connection *AppServicePbxConnection, msg *BaseMessage, message []byte)
}

type Guid string
