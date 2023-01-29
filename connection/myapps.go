package connection

import (
	"context"
	"crypto/rc4"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const session_length_usr = 32
const session_length_pwd = 23

var ReconnectTimeout = time.Second * 2

// Message struct
type Message struct {
	Mt string `json:"mt"`
}

// the Values for one Account
type Config struct {
	Host               string                 `yaml:"host"`               // IP or Hostname to initialy connect to. could be a DNS name of the pbx like 'pbx.company.com' or a IP address with a Port '192.168.33.11:433'
	InsecureSkipVerify bool                   `yaml:"insecureskipverify"` // disabled the check for valid SSL/TLS certificate on the outgoing websocket connection
	Username           string                 `yaml:"username"`           // Username of the pbx
	Password           string                 `yaml:"password"`           // Password to the Username
	SessionFilePath    string                 `yaml:"sessionfilepath"`    // Filename to a local JSON file to store the session. Will be created if it not exists
	UserAgent          string                 `yaml:"useragent"`          // the User Agnent shown in the list of current sessions in the user profile
	Handler            MessageHandlerRegister // list of message handler on the ssesion
	RedirectHost       string                 // is set, when the user is located not in the master and should open a connection to the secondary pbx
	Debug              bool                   `yaml:"debug"` // set to true to print log messages of the connection
}

func (config *Config) Println(v ...any) {
	if config.Debug {
		log.Println(fmt.Sprintf("[%s] (%s) ", config.Host, config.Username), v)
	}
}

func (config *Config) Printf(format string, a ...interface{}) {
	if config.Debug {
		log.Println(fmt.Sprintf("[%s] (%s) ", config.Host, config.Username) + fmt.Sprintf(format, a...))
	}
}

func (config *Config) StartSession(wg *sync.WaitGroup) {
	wg.Add(1)       // add goroutines to the wait group
	defer wg.Done() // mark the goroutine as done

	// start the websocket in a loop to reconnect if it failes/disconnects
	for {
		url := ""

		if config.RedirectHost != "" {
			url = fmt.Sprintf("wss://%s/PBX0/APPCLIENT/websocket", config.RedirectHost)
		} else {
			url = fmt.Sprintf("wss://%s/PBX0/APPCLIENT/websocket", config.Host)
		}
		config.Printf("connecting to %s", url)

		// Dialer configuration
		dialer := websocket.DefaultDialer

		// if `config.InsecureSkipVerify` is set to true, the TLS/SSL certificate is not checked
		dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: config.InsecureSkipVerify}

		// Connect to the WebSocket
		ctx, cancel := context.WithTimeout(context.Background(), ReconnectTimeout)
		conn, _, err := dialer.DialContext(ctx, url, http.Header{})
		if err != nil {
			config.Printf("connecting to url '%s' failed: %s", url, err)
			time.Sleep(ReconnectTimeout) // wait before trying to reconnect, avoid hammering
			cancel()                     // call cancel function here, It's used to stop the context's timer.
			continue
		}

		myappsSession := MyAppsConnection{
			Context: ctx,
			Conn:    conn,

			Config: config,
			Nonce:  GetRandomHexString(16),
		}

		// Add onDisconnect function
		conn.SetCloseHandler(func(code int, text string) error {
			onDisconnect(&myappsSession)
			conn.Close()
			return nil
		})

		err_handler := onConnect(&myappsSession)
		if err_handler != nil {
			config.Printf("Error in onConnect: %v", err_handler)
		}
		config.Println("WebSocket disconnected")
		time.Sleep(ReconnectTimeout) // wait before trying to reconnect

	}
}

// loads session keys from file
func (myappconfig *Config) GetSessionKeys() (string, string, error) {
	// check if file exists
	if _, err := os.Stat(myappconfig.SessionFilePath); os.IsNotExist(err) {
		return "", "", err
	}

	// read file
	file, err := ioutil.ReadFile(myappconfig.SessionFilePath)
	if err != nil {
		return "", "", err
	}

	// parse json
	var data struct {
		Usr string
		Pwd string
	}
	if err := json.Unmarshal(file, &data); err != nil {
		return "", "", err
	}
	if len(data.Usr) != session_length_usr {
		return "", "", fmt.Errorf("the Usr key of the stored session has not a length of %v", session_length_usr)
	}
	if len(data.Pwd) != session_length_pwd {
		return "", "", fmt.Errorf("the Pwd key of the stored session has not a length of %v", session_length_pwd)
	}
	return data.Usr, data.Pwd, nil
}

// writes session keys to file
func (myappconfig *Config) SaveSessionKeys(usr, pwd string) error {
	if len(usr) != session_length_usr {
		return fmt.Errorf("the Usr key of the stored session has not a length of %v", session_length_usr)
	}

	if len(pwd) != session_length_pwd {
		return fmt.Errorf("the Pwd key of the stored session has not a length of %v", session_length_pwd)
	}

	// create data
	data := struct {
		Usr string
		Pwd string
	}{
		Usr: usr,
		Pwd: pwd,
	}

	// marshal json
	file, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	// write file
	if err := ioutil.WriteFile(myappconfig.SessionFilePath, file, 0600); err != nil {
		return err
	}

	return nil
}

// Deletes the JSON file with stored sessions
func (myappconfig *Config) DeleteSessionKeys() error {
	// check if file exists
	if _, err := os.Stat(myappconfig.SessionFilePath); os.IsNotExist(err) {
		return nil
	}

	// remove file
	return os.Remove(myappconfig.SessionFilePath)
}

type MyAppsConnection struct {
	Context context.Context
	Conn    *websocket.Conn

	Config *Config

	Nonce string

	LoggedIn bool           // indicates if a user is logged in
	User     *MyAppUserInfo // the info object of the logged in user
}

func (myappsConnection *MyAppsConnection) send(message []byte) error {
	err := myappsConnection.Conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		myappsConnection.Config.Println("Error sending message:", err)
		return err
	}
	return nil
}

func (myappsConnection *MyAppsConnection) sendLogin() error {
	_, _, err_read_session_from_file := myappsConnection.Config.GetSessionKeys()
	if err_read_session_from_file != nil {
		myappsConnection.Config.Printf("reading session from file failed: %v", err_read_session_from_file)
	}
	if err_read_session_from_file != nil {
		// reading session from file failed, so we do a 'user' login with username/password
		myappsConnection.Config.Println("using username/password to login")
		msg, _ := json.Marshal(Login{"Login", "user", myappsConnection.Config.UserAgent, "", "", "", ""})
		return myappsConnection.send(msg)
	} else {
		// using the session keys to do a 'session' login
		myappsConnection.Config.Printf("using session to login")
		msg, _ := json.Marshal(Login{"Login", "session", myappsConnection.Config.UserAgent, "", "", "", ""})
		return myappsConnection.send(msg)
	}
}

func (myappsConnection *MyAppsConnection) received(message []byte) error {
	// Unmarshal message
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		myappsConnection.Config.Println("server: error unmarshalling message:", err)
	}

	switch msg["mt"] {
	default:
		myappsConnection.Config.Handler.HandleMessage(myappsConnection, fmt.Sprintf("%s", msg["mt"]), message)
	case "CheckBuildResult":
		var checkbuildresult CheckBuildResult
		err := json.Unmarshal(message, &checkbuildresult)
		if err != nil {
			myappsConnection.Config.Println("server: error unmarshalling message:", err)
			return err
		}
		myappsConnection.send([]byte(`{"mt":"SubscribeRegister"}`))

	case "UpdateRegister":
		var updateRegister UpdateRegister
		err := json.Unmarshal(message, &updateRegister)
		if err != nil {
			myappsConnection.Config.Println("server: error unmarshalling message:", err)
			return err
		}
		// msg, _ := json.Marshal(LoginInfo{"LoginInfo"})
		// myapps.send(msg)
		myappsConnection.sendLogin()

	case "LoginInfoResult":
		var loginInfoResult LoginInfoResult
		err := json.Unmarshal(message, &loginInfoResult)
		if err != nil {
			myappsConnection.Config.Println("server: error unmarshalling message:", err)
			return err
		}
		myappsConnection.sendLogin()

	case "Authenticate":
		var authenticate Authenticate
		err := json.Unmarshal(message, &authenticate)
		if err != nil {
			myappsConnection.Config.Println("server: error unmarshalling message:", err)
			return err
		}
		switch authenticate.Type {
		case "session":
			usr, pwd, err_read_session_from_file := myappsConnection.Config.GetSessionKeys()
			if err_read_session_from_file != nil {
				return err_read_session_from_file
			}
			response := GetLoginDigestDigest(authenticate.Type, authenticate.Domain, usr, pwd, myappsConnection.Nonce, authenticate.Challenge)
			msg, _ := json.Marshal(Login{
				"Login", authenticate.Type, myappsConnection.Config.UserAgent, "digest",
				usr, myappsConnection.Nonce, response,
			})
			myappsConnection.send(msg)

		case "user":

			response := GetLoginDigestDigest(authenticate.Type, authenticate.Domain, myappsConnection.Config.Username, myappsConnection.Config.Password, myappsConnection.Nonce, authenticate.Challenge)
			msg, _ := json.Marshal(Login{"Login", authenticate.Type, myappsConnection.Config.UserAgent, "digest", myappsConnection.Config.Username, myappsConnection.Nonce, response})
			myappsConnection.send(msg)
		default:
			myappsConnection.Config.Printf("unknown Authenticate.Type '%s'", authenticate.Type)
		}

	case "Authorize":
		var authorize Authorize
		err := json.Unmarshal(message, &authorize)
		if err != nil {
			myappsConnection.Config.Println("server: error unmarshalling message:", err)
			return err
		}
		myappsConnection.Config.Printf("Login needs a 2FA login verification. the code is: %v", authorize.Code)

	case "LoginResult":
		var loginResult LoginResult
		err := json.Unmarshal(message, &loginResult)
		if err != nil {
			myappsConnection.Config.Println("server: error unmarshalling message:", err)
			return err
		}
		if loginResult.Error == LoginResultSessionExpired {
			myappsConnection.Config.Printf("Login failed because '%s', deleting the stored session", loginResult.ErrorText)
			myappsConnection.Config.DeleteSessionKeys()

			// send a new login request, this time with type:user
			msg, _ := json.Marshal(LoginInfo{"LoginInfo"})
			myappsConnection.send(msg)
		}
		usr, _ := DecryptRc4(fmt.Sprintf("innovaphoneAppClient:usr:%v:%v", myappsConnection.Nonce, myappsConnection.Config.Password), loginResult.Info.Session.Usr)
		pwd, _ := DecryptRc4(fmt.Sprintf("innovaphoneAppClient:pwd:%v:%v", myappsConnection.Nonce, myappsConnection.Config.Password), loginResult.Info.Session.Pwd)

		myappsConnection.Config.SaveSessionKeys(usr, pwd)

		myappsConnection.OnUserLoggedIn(&loginResult.Info.User)
		myappsConnection.send([]byte(`{"mt":"SubscribeApps"}`))
		myappsConnection.send([]byte(`{"mt":"SubscribePresence","sip":"chat"}`))

	case "Redirect":
		var redirect Redirect
		err := json.Unmarshal(message, &redirect)
		if err != nil {
			myappsConnection.Config.Println("server: error unmarshalling message:", err)
			return err
		}

		usr, _ := DecryptRc4(fmt.Sprintf("innovaphoneAppClient:usr:%v:%v", myappsConnection.Nonce, myappsConnection.Config.Password), redirect.Info.Session.Usr)
		pwd, _ := DecryptRc4(fmt.Sprintf("innovaphoneAppClient:pwd:%v:%v", myappsConnection.Nonce, myappsConnection.Config.Password), redirect.Info.Session.Pwd)

		myappsConnection.Config.SaveSessionKeys(usr, pwd)
		myappsConnection.Config.Printf("login successful, but we need to redirect to '%s'", redirect.Info.Host)
		myappsConnection.Config.RedirectHost = redirect.Info.Host

		myappsConnection.Conn.Close()
		myappsConnection.Context.Done()
	}
	return nil
}

func (myappsConnection *MyAppsConnection) OnUserLoggedIn(user *MyAppUserInfo) {
	myappsConnection.LoggedIn = true
	myappsConnection.User = user
	myappsConnection.Config.Printf("login successful, user sip='%s',dn='%s',guid='%s' is now authenticated", myappsConnection.User.Sip, myappsConnection.User.Dn, myappsConnection.User.Guid)
}

func GetRandomHexString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func Hex2bin(inputstring string) []int {
	data, _ := hex.DecodeString(inputstring)
	list_of_int := make([]int, len(data))
	for i, b := range data {
		list_of_int[i] = int(b)
	}
	return list_of_int
}

func DecryptRc4(key string, data string) (string, error) {
	keyBytes := []byte(key)
	ciphertext, _ := hex.DecodeString(data)
	plaintext := make([]byte, len(ciphertext))
	c, _ := rc4.NewCipher(keyBytes)
	c.XORKeyStream(plaintext, ciphertext)
	return string(plaintext), nil
}

func GetLoginDigestDigest(type_ string, domain string, username string, password string, nonce string, challenge string) string {
	txt := fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s", "innovaphoneAppClient", type_, domain, username, password, nonce, challenge)
	m := sha256.New()
	m.Write([]byte(txt))
	return fmt.Sprintf("%x", m.Sum(nil))
}

func onDisconnect(myappsConnection *MyAppsConnection) {
	myappsConnection.Config.Println("WebSocket disconnected")
}

func onConnect(myappsConnection *MyAppsConnection) error {
	myappsConnection.Config.Printf("connected to myApps socket at %s", myappsConnection.Conn.RemoteAddr())

	// Send JSON-encoded message on connect
	message := Message{Mt: "CheckBuild"} //, Url: fmt.Sprintf("https://%s/PBX0/APPCLIENT/appclient.htm", myappsConnection.Conn.RemoteAddr()), Always: false}
	messageJSON, _ := json.Marshal(message)
	myappsConnection.send(messageJSON)

	// Send and receive messages
	// ...
	for {
		// Read message
		_, message, err := myappsConnection.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				myappsConnection.Config.Printf("error: %v", err)
				return err
			}
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				myappsConnection.Config.Println("server: websocket closed")
				return err
			}
			break
		}

		myappsConnection.received(message)
	}

	return nil
}
