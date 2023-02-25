package sysclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ricoschulte/go-myapps/connection"
)

type Sysclient struct {
	Identity           Identity
	Url                string
	Timeout            time.Duration
	InsecureSkipVerify bool
	Context            context.Context
	Conn               *websocket.Conn
	Tunnels            map[int32]*SysclientTunnel // map of active tunnels indexed by the sessionid
	ServeMux           *http.ServeMux             // the instance of the Http Server Mux to handle Http Requests

	FileSysclientPassword      string // filename to store
	FileAdministrativePassword string // filename to store
	SecretKey                  []byte // key to encrypt the local files as []bytes
}

func NewSysclient(identity Identity, url string, timeout time.Duration, insecureSkipVerify bool, mux *http.ServeMux, fileSysclientPassword string, fileAdministrativePassword string, secretkey string) (*Sysclient, error) {
	if fileSysclientPassword == "" {
		return nil, errors.New("fileSysclientPassword cant be empty")
	}
	if fileAdministrativePassword == "" {
		return nil, errors.New("fileAdministrativePassword cant be empty")
	}

	sysclient := &Sysclient{
		Identity:           identity,
		Url:                url,
		Timeout:            timeout,
		InsecureSkipVerify: insecureSkipVerify,
		Tunnels:            map[int32]*SysclientTunnel{},
		ServeMux:           mux,

		FileSysclientPassword:      fileSysclientPassword,
		FileAdministrativePassword: fileAdministrativePassword,
		SecretKey:                  []byte(secretkey),
	}

	return sysclient, nil

}
func (sc *Sysclient) Println(v ...any) {
	log.Printf("{%s}\n", v...)

}

func (sc *Sysclient) Printf(format string, a ...interface{}) {
	log.Printf("{%s} %s\n", "sysclient", fmt.Sprintf(format, a...))

}

func (sc *Sysclient) Connect() error {

	for {
		sc.Printf("connecting to %s as %v/%v", sc.Url, sc.Identity.Id, sc.Identity.Product)

		// Dialer configuration
		dialer := websocket.DefaultDialer

		// if `config.InsecureSkipVerify` is set to true, the TLS/SSL certificate is not checked
		dialer.TLSClientConfig = &tls.Config{InsecureSkipVerify: sc.InsecureSkipVerify}

		// Connect to the WebSocket
		ctx, cancel := context.WithTimeout(context.Background(), connection.ReconnectTimeout)
		conn, _, err := dialer.DialContext(ctx, sc.Url, http.Header{})
		if err != nil {
			sc.Printf("error while connecting to url '%s': %+v", sc.Url, err)
			time.Sleep(connection.ReconnectTimeout) // wait before trying to reconnect, avoid hammering
			cancel()                                // call cancel function here, It's used to stop the context's timer.
			continue
		}

		sc.Context = ctx
		sc.Conn = conn

		// Add onDisconnect function
		conn.SetCloseHandler(func(code int, text string) error {
			fmt.Printf("connection closed %v %v", code, text)
			sc.onDisconnect(code, text)
			cancel()
			conn.Close()
			return nil
		})

		err_handler := sc.onConnect()
		if err_handler != nil {
			sc.Printf("Error in onConnect: %v", err_handler)
		}

		sc.Println("WebSocket disconnected")
		time.Sleep(connection.ReconnectTimeout) // wait before trying to reconnect
	}

}

func (sc *Sysclient) onConnect() error {
	sc.Printf("WebSocket to '%s' connected", sc.Url)
	mess := []byte{MessageTypeAdmin}
	mess = append(mess, AdminSendIdentify...)

	msgb, err := sc.Identity.ToBytes()
	if err != nil {
		return err
	}
	mess = append(mess, msgb...)

	identify, err_message := NewAdminMessage(mess)
	if err_message != nil {
		return err_message
	}
	err_send := sc.Send(identify.AsBytes())
	if err_send != nil {
		return err_send
	}
	return sc.ReadWriteLoop()

}

func (sc *Sysclient) onDisconnect(code int, text string) error {
	sc.Printf("WebSocket to '%s' disconnected: %v %v", sc.Url, code, text)
	return nil
}

func (sc *Sysclient) Send(message []byte) error {
	//sc.Printf("+++++++++++++++++++ Sysclient Send: %v\n\n", message)
	err := sc.Conn.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		fmt.Println("Error sending message:", err)
		return err
	}
	return nil
}

func (sc *Sysclient) ReadWriteLoop() error {
	for {
		// Read message
		messagetype, message, err := sc.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				sc.Printf("error: %v", err)
				return err
			}
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				sc.Println("server: websocket closed")
				return err
			}
			break
		}
		if messagetype == websocket.BinaryMessage {
			sc.received(message)
		}

	}
	return nil
}

// called when we received a message from the appservice
func (sc *Sysclient) received(message []byte) error {
	if len(message) < 3 {
		return fmt.Errorf("invalid length of message: %v", len(message))
	}
	//	fmt.Printf("\n\n+++++++++++++++++++ Sysclient received: %v\n", message)
	switch message[0] {
	default:
		return fmt.Errorf("received unknown message from the sysclient server: %v", string(message[0]))
	case MessageTypeAdmin:
		msg, err_message := NewAdminMessage(message)
		if err_message != nil {
			return fmt.Errorf("invalid message: %v", err_message)
		}

		response, err_response := sc.HandleAdminMessage(msg)
		if err_response != nil {
			return err_response
		}

		// if we need to send a response
		if response != nil {
			err_send := sc.Send(response.AsBytes())
			if err_send != nil {
				return err_send
			}
		}
	case MessageTypeTunnel:

		msg_received, err_message := NewTunnelMessage(message)
		if err_message != nil {
			return fmt.Errorf("invalid message: %v", err_message)
		}
		msg_received_tunnel_id, err_tunnelId := msg_received.GetSessionId()
		if err_tunnelId != nil {
			return fmt.Errorf("invalid message: %v", err_tunnelId)
		}

		tunnel := sc.Tunnels[msg_received_tunnel_id]
		if tunnel == nil {
			// create a new tunnel for that sessionId
			tunnel_new, err_create_tunnel := NewSysclientTunnel(msg_received.SessionId, sc.ServeMux)
			if err_create_tunnel != nil {
				return fmt.Errorf("error while creating new tunnelsession: %v", err_create_tunnel)
			}
			sc.Tunnels[msg_received_tunnel_id] = tunnel_new
			tunnel = tunnel_new
		}

		response, err_handle := tunnel.HandleRequest(msg_received)
		if err_handle != nil {
			return err_handle
		}

		// remove closed tunnel from list
		if bytes.Equal(response.EventType, TunnelShutdown) {
			sc.Tunnels[msg_received_tunnel_id] = nil
		}

		err_send := sc.Send(response.AsBytes())
		if err_send != nil {
			return err_send
		}

	}

	return nil
}

func (sc *Sysclient) HandleAdminMessage(messageIn *AdminMessage) (*AdminMessage, error) {
	switch true {
	default:
		return nil, fmt.Errorf("unknown Admin Message of Type %v", messageIn.Type)

	case bytes.Equal(messageIn.Command, AdminReceiveSysclientPassword):
		err := messageIn.HandleAdminReceiveSysclientPassword(sc.SecretKey, sc.FileSysclientPassword)
		if err != nil {
			return nil, err
		}
		return nil, nil
	case bytes.Equal(messageIn.Command, AdminReceiveChallenge):
		return messageIn.HandleAdminReceiveChallenge(sc.SecretKey, &sc.Identity, sc.FileSysclientPassword)

	case bytes.Equal(messageIn.Command, AdminReceiveNewAdminPassword):
		err := messageIn.handleAdminReceiveNewAdministrativePassword(sc.SecretKey, sc.FileSysclientPassword, sc.FileAdministrativePassword)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}
