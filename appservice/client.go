package appservice

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ricoschulte/go-myapps/connection"
)

type AppServiceClient struct {
	MyAppsConnection        *connection.MyAppsConnection
	Context                 context.Context
	Conn                    *websocket.Conn
	AppInfo                 *connection.App
	LoggedIn                bool
	MessageHandlerRegister  *AppServiceMessageHandlerRegister  // list of message handler on the session
	CallbackHandlerRegister *AppServiceCallbackHandlerRegister // hold a list of callback handler that are registered for a message with src attribute

}

func NewAppServiceClient() *AppServiceClient {
	asclient := &AppServiceClient{}
	asclient.CallbackHandlerRegister = NewAppServiceCallbackHandlerRegister()

	return asclient

}
func (ac *AppServiceClient) Println(v ...any) {
	if ac.MyAppsConnection.Config.Debug {
		ac.MyAppsConnection.Config.Printf(fmt.Sprintf("{%s} ", ac.AppInfo.Name), v)
	}
}

func (ac *AppServiceClient) Printf(format string, a ...interface{}) {
	if ac.MyAppsConnection.Config.Debug {
		ac.MyAppsConnection.Config.Printf(fmt.Sprintf("{%s} ", ac.AppInfo.Name) + fmt.Sprintf(format, a...))
	}
}

func (ac *AppServiceClient) Connect() error {
	url := ""

	switch true {
	default:
		url = ac.AppInfo.Url
	case strings.HasPrefix(ac.AppInfo.Url, "http://"):
		url = strings.Replace(ac.AppInfo.Url, "http", "ws", 1)
	case strings.HasPrefix(ac.AppInfo.Url, "https://"):
		url = strings.Replace(ac.AppInfo.Url, "https", "wss", 1)
	case strings.HasPrefix(ac.AppInfo.Url, "../../"):
		// #TODO handle
		// ../../APPS/chat/chat
		return fmt.Errorf("error, currently apps hosted on a pbx are not supported: %+v", ac.AppInfo)
	}

	for {
		ac.Printf("connecting to %s\n", url)

		// Dialer configuration
		dialer := websocket.DefaultDialer

		// if `config.InsecureSkipVerify` is set to true, the TLS/SSL certificate is not checked
		dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: ac.MyAppsConnection.Config.InsecureSkipVerify}

		// Connect to the WebSocket
		ctx, cancel := context.WithTimeout(context.Background(), connection.ReconnectTimeout)
		conn, _, err := dialer.DialContext(ctx, url, http.Header{})
		if err != nil {
			ac.Printf("error while connecting to url '%s': %+v", url, err)
			time.Sleep(connection.ReconnectTimeout) // wait before trying to reconnect, avoid hammering
			cancel()                                // call cancel function here, It's used to stop the context's timer.
			continue
		}

		ac.Context = ctx
		ac.Conn = conn

		// Add onDisconnect function
		conn.SetCloseHandler(func(code int, text string) error {
			ac.onDisconnect(code, text)
			cancel()
			conn.Close()
			return nil
		})

		err_handler := ac.onConnect()
		if err_handler != nil {
			ac.Printf("Error in onConnect: %v", err_handler)
		}

		ac.Println("WebSocket disconnected")
		time.Sleep(connection.ReconnectTimeout) // wait before trying to reconnect
	}

}

func (ac *AppServiceClient) onConnect() error {
	ac.Printf("WebSocket to '%s' connected\n", ac.AppInfo.Name)
	ac.Send([]byte(`{"mt":"AppChallenge"}`))
	ac.ReadWriteLoop()
	return nil
}

func (ac *AppServiceClient) onDisconnect(code int, text string) error {
	ac.Printf("WebSocket to '%s' disconnected: %v %v\n", ac.AppInfo.Name, code, text)
	return nil
}

func (ac *AppServiceClient) Send(message []byte) error {
	ac.Printf("sending message: %s", string(message))
	err := ac.Conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return err
	}
	return nil
}

func (ac *AppServiceClient) ReadWriteLoop() error {
	for {
		// Read message
		_, message, err := ac.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				ac.Printf("error: %v", err)
				return err
			}
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				ac.Println("server: websocket closed")
				return err
			}
			break
		}

		ac.received(message)
	}
	return nil
}

// called when we received a message from the appservice
func (ac *AppServiceClient) received(message []byte) error {
	var msg connection.Message
	if err := json.Unmarshal(message, &msg); err != nil {
		ac.Println("error unmarshalling message from appservice:", err)
	}

	// var data map[string]interface{}
	// json.Unmarshal(message, &data)
	// jsonData, _ := json.MarshalIndent(data, "", "  ")
	// ac.Println(string(jsonData))

	switch msg.Mt {
	// default:
	// 	ac.Printf("received unknown message from the appservice: %v", string(message))

	case "AppChallengeResult":
		var msgl AppChallengeResult
		if err := json.Unmarshal(message, &msgl); err != nil {
			ac.Println("server: error unmarshalling AppChallengeResult:", err)
		}
		msg := NewAppGetLogin(ac.AppInfo.Name, msgl.Challenge)
		msgs, errm := json.Marshal(msg)
		if errm != nil {
			ac.Println("server: error marshalling AppGetLogin:", errm)
		}
		ac.MyAppsConnection.SendWithResult(msgs, msg.Src, 1, ac)
	case "AppLoginResult":
		var msgl AppLoginResult
		if err := json.Unmarshal(message, &msgl); err != nil {
			ac.Println("server: error unmarshalling AppChallengeResult:", err)
		}

		if !msgl.Ok {
			ac.Printf("error in response of the login to the appservice with the infos from the pbx: %+v", msgl)
			ac.LoggedIn = false
		} else {
			ac.Printf("login successful")

			ac.LoggedIn = true
		}

	}

	err := ac.MessageHandlerRegister.HandleMessage(ac, msg.Mt, message)
	if err != nil {
		fmt.Println("erro", err)
	}
	return nil
}

// handles the callback received from the pbx in this case
func (ac *AppServiceClient) HandleCallbackMessage(myappsconnection *connection.MyAppsConnection, message []byte) error {
	var msg connection.Message
	if err := json.Unmarshal(message, &msg); err != nil {
		ac.Println("error unmarshalling message from pbx:", err)
	}
	switch msg.Mt {
	default:
		ac.Printf("received unknown message from pbx: %v", string(message))
	case "AppGetLoginResult":
		var msgl AppGetLoginResult
		if err := json.Unmarshal(message, &msgl); err != nil {
			ac.Println("error unmarshalling AppChallengeResult:", err)
		}

		msg := NewAppLogin(msgl)
		msgs, errm := json.Marshal(msg)
		if errm != nil {
			ac.Println("error marshalling AppLogin:", errm)
		}
		ac.Send(msgs)
	}
	return nil
}
