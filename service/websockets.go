package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var challenges = make(map[*websocket.Conn]string)

type BaseMessage struct {
	Api string `json:"api"`
	Mt  string `json:"mt"`
	Src string `json:"src,omitempty"`
}

type AppChallengeResult struct {
	Mt        string `json:"mt"`
	Challenge string `json:"challenge"`
}

type AppLogin struct {
	BaseMessage
	Sip    string `json:"sip"`
	Guid   string `json:"guid"`
	Dn     string `json:"dn"`
	Digest string `json:"digest"`
	Domain string `json:"domain"`
	App    string `json:"app"`
	//Info   AppLoginInfo `json:"info"`
	Info   json.RawMessage `json:"info"`
	PbxObj string          `json:"pbxObj,omitempty"`
}

func (al AppLogin) GetInfoObject() (*AppLoginInfo, error) {
	obj := AppLoginInfo{}
	err := json.Unmarshal(al.Info, &obj)
	if err != nil {
		return nil, err
	}
	return &obj, err
}

func (al AppLogin) InfoAsPbxobjectDigestString() (string, error) {
	msg, err := json.Marshal(al.Info)
	return string(msg), err
}

func (al AppLogin) InfoAsUserDigestString() (string, error) {
	msg, err := json.Marshal(al.Info)
	return string(msg), err
}

type AppLoginInfo struct {
	Appobj     string   `json:"appobj"`
	Appdn      string   `json:"appdn"`
	Appurl     string   `json:"appurl"`
	Pbx        string   `json:"pbx"`
	Cn         string   `json:"cn"`
	Unlicensed bool     `json:"unlicensed,omitempty"`
	Testmode   bool     `json:"testmode,omitempty"`
	Groups     []string `json:"groups"`
	Apps       []string `json:"apps"`
}

func (a *AppLoginInfo) MarshalJSON() ([]byte, error) {
	type Alias AppLoginInfo
	if len(a.Groups) == 0 {
		a.Groups = []string{}
	}
	return json.Marshal(&struct {
		*Alias
	}{
		Alias: (*Alias)(a),
	})
}

type AppLoginResult struct {
	BaseMessage
	App    string `json:"app"`
	Domain string `json:"domain"`
	Ok     bool   `json:"ok"`
}

type AppInfo struct {
	BaseMessage
	App string `json:"app"`
}

type AppInfoInfo struct {
	Hidden bool `json:"hidden"`
	Apis   struct {
		ComInnovaphoneSearch map[string]interface{} `json:"com.innovaphone.search"`
	} `json:"apis"`
}

type AppInfoServiceInfo struct {
	Apis struct {
		ComInnovaphoneReplicator map[string]interface{} `json:"com.innovaphone.replicator"`
	} `json:"apis"`
}

type AppInfoResult struct {
	BaseMessage
	App         string             `json:"app"`
	Info        AppInfoInfo        `json:"info"`
	Serviceinfo AppInfoServiceInfo `json:"sericeinfo"`
}

type PbxInfo struct {
	BaseMessage
	Domain  string   `json:"domain"`
	Pbx     string   `json:"pbx"`
	PbxDns  string   `json:"pbxDns"`
	Apis    []string `json:"apis"`
	Version string   `json:"version"`
	Build   string   `json:"build"`
}

type AppServicePbxConnection struct {
	AppService    *AppService
	conn          *websocket.Conn
	PbxInfo       PbxInfo
	AppLogin      AppLogin
	Info          *AppLoginInfo
	Authenticated bool
	WriteMutext   sync.Mutex
}

func NewAppServicePbxConnection(appservice *AppService, conn *websocket.Conn) *AppServicePbxConnection {
	return &AppServicePbxConnection{
		AppService: appservice,
		conn:       conn,
	}
}
func (connection *AppServicePbxConnection) log() *log.Entry {
	return log.
		WithField("Domain", connection.AppService.Domain).
		WithField("Instance", connection.AppService.Instance).
		WithField("name", connection.AppService.Name).
		WithField("pbx", connection.PbxInfo.Pbx).
		WithField("pbxDns", connection.PbxInfo.PbxDns)
}

func (connection *AppServicePbxConnection) Loop() {

	for {
		// Read message
		_, message, err := connection.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				connection.log().Errorf("websocket error: %s", err)
				return
			}
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				connection.log().Debugf("websocket closed: %s", err)
				return
			}
			break
		}

		connection.log().Tracef("received message: %s", message)

		// Unmarshal message
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			connection.log().Errorf("server: error unmarshalling message: %v", err)
			continue
		}

		// Check message type
		mt, ok := msg["mt"].(string)
		if !ok {
			connection.log().Warning("received a message without 'mt'")
			continue
		}
		var response []byte
		var hErr error

		switch mt {
		case "AppChallenge":
			response, hErr = connection.handleAppChallenge(connection.conn, message)
		case "AppLogin":
			response, hErr = connection.handleAppLogin(challenges[connection.conn], message)
		case "AppInfo":
			response, hErr = connection.handleAppInfo(connection.conn, message)
		case "PbxInfo":
			var pbxInfo PbxInfo

			if err := json.Unmarshal([]byte(message), &pbxInfo); err != nil {
				connection.log().Errorf("unmarshal PbxInfo failed: %v", err)
			}
			connection.PbxInfo = pbxInfo
			connection.AppService.HandleApiConnected(connection, message)

		default:
			if _, ok := msg["api"]; ok {
				// Key "api" exists in the map
				if !connection.Authenticated {
					log.Warn("message for a api received but connection isnt authenticated. Closing connection.")
					connection.conn.Close()

				} else {
					// look for a api handler in app service
					go func() {
						connection.AppService.HandleApiMessage(connection, msg["api"].(string), message)
					}()
				}
				continue
			} else {
				connection.log().Warnf("unknown mt: %s", string(message))
				response = []byte("{\"mt\":\"Error\",\"text\":\"unknown mt '" + mt + "'\"}")
				connection.WriteMessage(response)
				continue
			}

		}

		// check if handler function returns an error
		if hErr != nil {
			connection.log().Errorf("error handling message: %v", hErr)
			continue
		}

		if string(response) == "" {
			response = []byte("{\"mt\":\"Error\",\"text\":\"the server has no content to send\"}")
			connection.WriteMessage(response)
		} else {
			connection.log().Tracef("sending Response: %s", response)
			// send response
			connection.WriteMessage(response)
		}
	}
}

func (connection *AppServicePbxConnection) WriteMessage(message []byte) error {
	connection.WriteMutext.Lock()
	defer connection.WriteMutext.Unlock()
	err := connection.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Errorf("Error writing message:", err)
		return err
	}
	return nil
}

func (connection *AppServicePbxConnection) handleAppChallenge(conn *websocket.Conn, msg []byte) ([]byte, error) {
	response := &AppChallengeResult{
		Mt:        "AppChallengeResult",
		Challenge: challenges[conn],
	}
	return json.Marshal(response)
}

func (connection *AppServicePbxConnection) handleAppLogin(challenge string, msg []byte) ([]byte, error) {
	var msgin AppLogin
	if err := json.Unmarshal([]byte(msg), &msgin); err != nil {
		return nil, err
	}
	info, err := msgin.GetInfoObject()
	if err != nil {
		return nil, err
	}
	connection.AppLogin = msgin
	connection.Info = info
	mu := &MyAppsUtils{}
	calculated_digest, err := mu.GetDigestForAppLoginFromJson(string(msg), connection.AppService.Password, challenge)
	if err != nil {
		return nil, err
	}

	if msgin.Digest == calculated_digest {
		log.
			WithField("Domain", connection.AppService.Domain).
			WithField("Instance", connection.AppService.Instance).
			WithField("name", connection.AppService.Name).
			WithField("app", msgin.App).
			WithField("sip", msgin.Sip).
			Debug("appservice login successful")
		response := AppLoginResult{
			BaseMessage: BaseMessage{
				Mt: "AppLoginResult",
			},
			App:    msgin.App,
			Domain: msgin.Domain,
			Ok:     true,
		}
		connection.Authenticated = true
		go func() {
			app := strings.TrimSuffix(msgin.App, ".htm")
			switch app {
			case "user":

				connection.AppService.HandleUserConnected(connection, msg)
			case "admin":
				connection.AppService.HandleAdminConnected(connection, msg)
			default:
				log.Warnf("unknown App '%s'", app)
			}
		}()
		return json.Marshal(response)
	} else {
		log.
			WithField("Domain", connection.AppService.Domain).
			WithField("Instance", connection.AppService.Instance).
			WithField("name", connection.AppService.Name).
			WithField("app", msgin.App).
			WithField("sip", msgin.Sip).
			Warn("appservice login failed: digest not correct")

		response := AppLoginResult{
			BaseMessage: BaseMessage{
				Mt: "AppLoginResult",
			},
			App:    msgin.App,
			Domain: msgin.Domain,
			Ok:     false,
		}
		connection.Authenticated = false
		return json.Marshal(response)
	}

}

func (connection *AppServicePbxConnection) handleAppInfo(conn *websocket.Conn, msg []byte) ([]byte, error) {
	log.
		WithField("msg", string(msg)).
		Trace("handleAppInfo")
	var msgin AppInfo
	if err := json.Unmarshal([]byte(msg), &msgin); err != nil {
		log.
			WithField("msg", string(msg)).
			Errorf("could not unmashal AppInfo: %v", err)
		return nil, err
	}
	resp := `{
		"mt": "AppInfoResult",
		"app": "%s",
		"info": {
			"hidden": %t,
			"apis": {
				"com.innovaphone.search": {}
			}
		},
		"serviceInfo": {
			"apis": {
				"com.innovaphone.replicator": {}
			}
		}
	}`
	var jsn string
	switch msgin.App {
	default:
		jsn = fmt.Sprintf(resp, msgin.App, false)

	case "searchapi":
		jsn = fmt.Sprintf(resp, msgin.App, true)

	}

	// pasre from json
	var respmsg AppInfoResult
	if err := json.Unmarshal([]byte(jsn), &respmsg); err != nil {
		log.
			WithField("msg", string(msg)).
			Errorf("could not unmashal AppInfoResult: %v", err)
		return nil, err
	}

	// to json and return
	return json.Marshal(respmsg)
}

func HandleWebsocket(appservice *AppService, w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	mlog := log.
		WithField("Domain", appservice.Domain).
		WithField("Instance", appservice.Instance).
		WithField("name", appservice.Name)

	defer func() {
		mlog.Debug("server: websocket handler ended")
		// handling panic  "repeated read on failed websocket connection"
		// This will recover from the panic and log the error, then close the connection by sending a 500 Internal Server Error response to the client.
		// This will prevent the panic from propagating and crashing your server.
		if r := recover(); r != nil {
			log.
				WithField("webserver", handleConnectionForHttpOrWebsocket).
				Debugf("Recovered from panic: %v", r)
			w.WriteHeader(http.StatusInternalServerError)

		}
	}()
	mlog.Debug("server: websocket handler started")

	conn, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		mlog.Errorf("upgrading websocket failed: %s", err)
		return
	}
	connection := NewAppServicePbxConnection(appservice, conn)
	appservice.AddConnection(connection)
	defer func() {
		conn.Close()
		appservice.DeleteConnection(connection)
		appservice.HandleApiDisConnected(connection)
	}()

	// Generate challenge for the connection
	mu := &MyAppsUtils{}
	challenge := mu.GetRandomHexString(16)
	challenges[conn] = challenge

	select {
	case <-ctx.Done():
		err := ctx.Err()

		mlog.Errorf("websocket error: %s", err)
		internalError := http.StatusInternalServerError
		http.Error(w, err.Error(), internalError)
		delete(challenges, conn)
		return
	default:
		connection.Loop()
	}
}
