package sysclient

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// split Tunnel responses into chunks of this size of []byte
var chunkSize = 8000

// Sysclient message header
var (
	MessageTypeTunnel byte = 0x02

	TunnelSend          []byte = []byte{0x00, 0x00, 0x00, 0x00} // socket send
	TunnelReceive       []byte = []byte{0x00, 0x00, 0x00, 0x01} // socket receive
	TunnelReceiveResult []byte = []byte{0x00, 0x00, 0x00, 0x02} // receive result
	TunnelSendResult    []byte = []byte{0x00, 0x00, 0x00, 0x03} // socket send result
	TunnelShutdown      []byte = []byte{0x00, 0x00, 0x00, 0x04} // socket shutdown
)

// SMT | Byte offset
//
//	   |  0      1      2      3      4      5      6      7      8      9     ...
//	   +------+------+------+------+------+------+------+------+------+------+------+
//	2  |      |        Session-ID         |         Event-Type        |  Data.....
//	   +------+------+------+------+------+------+------+------+------+------+------+
type TunnelMessage struct {
	Type      byte
	SessionId []byte // 4 bytes
	EventType []byte // 4 bytes
	Data      []byte // rest of the bytes
}

func NewTunnelMessage(message []byte) (*TunnelMessage, error) {
	if len(message) < 9 {
		return nil, errors.New("invalid length of message")
	}

	m := &TunnelMessage{
		Type:      message[0],
		SessionId: message[1:5],
		EventType: message[5:9],
		Data:      message[9:],
	}

	if m.Type != MessageTypeTunnel {
		return nil, errors.New("message is not a TunnelMessage")
	}

	return m, nil
}

func (tm *TunnelMessage) GetSessionId() (int32, error) {

	if len(tm.SessionId) != 4 {
		return int32(0), errors.New("invalid length of SessionId to convert to int32")
	}
	return int32(binary.BigEndian.Uint32(tm.SessionId)), nil
}

func (tm *TunnelMessage) AsBytes() []byte {
	message_as_bytes := []byte{tm.Type}
	message_as_bytes = append(message_as_bytes, tm.SessionId...)
	message_as_bytes = append(message_as_bytes, tm.EventType...)
	message_as_bytes = append(message_as_bytes, tm.Data...)
	return message_as_bytes
}

type SysclientTunnel struct {
	ServeMux       *http.ServeMux
	SessionId      []byte
	Response       []byte
	ResponseChunks [][]byte
}

func NewSysclientTunnel(session []byte, servermux *http.ServeMux) (*SysclientTunnel, error) {

	tunnel := &SysclientTunnel{}
	tunnel.ServeMux = servermux
	tunnel.SessionId = session
	tunnel.Response = []byte{}
	tunnel.ResponseChunks = [][]byte{}

	return tunnel, nil
}

func (tunnel *SysclientTunnel) HandleRequest(message *TunnelMessage) (*TunnelMessage, error) {

	switch true {
	default:
		return nil, fmt.Errorf("unknown MessageType: %v", message.EventType)
	case bytes.Equal(message.EventType, TunnelSend):
		resp, err_handle := tunnel.HandleTunnelSend(message)
		if err_handle != nil {
			return nil, err_handle
		}
		return resp, nil
	case bytes.Equal(message.EventType, TunnelReceive):
		resp, err_handle := tunnel.HandleTunnelReceive(message)
		if err_handle != nil {
			return nil, err_handle
		}
		return resp, nil
	case bytes.Equal(message.EventType, TunnelShutdown):
		response_bytes := []byte{MessageTypeTunnel}
		response_bytes = append(response_bytes, tunnel.SessionId...)
		response_bytes = append(response_bytes, TunnelShutdown...)
		resp, err_responsemessage := NewTunnelMessage(response_bytes)
		if err_responsemessage != nil {
			return nil, err_responsemessage
		}
		return resp, nil

	}

}
func (tunnel *SysclientTunnel) GetNextChunk() ([]byte, error) {
	if len(tunnel.ResponseChunks) > 0 {
		return tunnel.ResponseChunks[0], nil
	} else {
		return nil, errors.New("no chunks left to respond")
	}
}

/*
the sysclient server is requesting data with http data

	in:
		TunnelSend
		02 0000000d 00000000 474554202f7374617469632f6a732f706167652e6a7320485454502f312e310d0a486f73743a20617070732e667269747a2e626f780d0a…]
	out:
		TunnelSendResult
		02 0000000d 00000003
*/
func (tunnel *SysclientTunnel) HandleTunnelSend(message *TunnelMessage) (*TunnelMessage, error) {
	// read Data of input http
	request, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(message.Data)))
	if err != nil {
		return nil, fmt.Errorf("error reading request: %v", err)
	}

	// Use the custom request handler to handle the custom request
	w := &HttpResponseWriter{}
	tunnel.ServeMux.ServeHTTP(w, request)

	tunnel.Response = w.GetBytes()

	// split http response to chunks if needed
	for i := 0; i < len(tunnel.Response); i += chunkSize {
		end := i + chunkSize
		if end > len(tunnel.Response) {
			end = len(tunnel.Response)
		}
		tunnel.ResponseChunks = append(tunnel.ResponseChunks, tunnel.Response[i:end])
	}

	// return with a TunnelSendResult without result
	response_bytes := []byte{MessageTypeTunnel}
	response_bytes = append(response_bytes, tunnel.SessionId...)
	response_bytes = append(response_bytes, TunnelSendResult...)

	response_message, err_responsemessage := NewTunnelMessage(response_bytes)
	if err_responsemessage != nil {
		return nil, err_responsemessage
	}

	return response_message, nil
}

/*
in:

	TunnelReceive
	02 0000000d 00000001 00000000

out:

	TunnelReceiveResult
	02 0000000d 00000002 485454502f312e3120323030204f4b0d0a436f6e74656e742d547970653a206170706c69636174696f6e2f6a6176617363726970740d0a…
*/
func (tunnel *SysclientTunnel) HandleTunnelReceive(message *TunnelMessage) (*TunnelMessage, error) {

	response_bytes := []byte{MessageTypeTunnel}
	response_bytes = append(response_bytes, tunnel.SessionId...)
	if len(tunnel.ResponseChunks) > 0 {
		// return with a TunnelReceiveResult
		data := tunnel.ResponseChunks[0]
		tunnel.ResponseChunks = tunnel.ResponseChunks[1:]

		if len(tunnel.ResponseChunks) > 0 {
			// the are more chunks
			response_bytes = append(response_bytes, []byte{0x00, 0x00, 0x00, 0x02}...)
		} else {
			// this is the last chunk
			response_bytes = append(response_bytes, []byte{0x00, 0x00, 0x00, 0x02}...)

		}
		response_bytes = append(response_bytes, data...)
	} else {
		// return with a TunnelShutdown
		return nil, errors.New("no chunks left to send")
	}

	response_message, err_responsemessage := NewTunnelMessage(response_bytes)
	if err_responsemessage != nil {
		return nil, err_responsemessage
	}

	return response_message, nil
}

// Define a HttpResponseWriter ResponseWriter
type HttpResponseWriter struct {
	response string
	status   int
	header   http.Header
}

func (w *HttpResponseWriter) Write(p []byte) (int, error) {
	w.response += string(p)
	return len(w.response), nil
}

func (w *HttpResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = http.Header{}
	}
	return w.header
}

func (w *HttpResponseWriter) WriteHeader(status int) {
	w.status = status
}
func (w *HttpResponseWriter) GetBytes() []byte {
	bu := fmt.Sprintf("HTTP/1.1 %v\r\n", http.StatusText(w.status))
	bu += fmt.Sprintf("%s: %s\r\n", "Content-Length", strconv.Itoa(len(w.response)))
	bu += fmt.Sprintf("%s: %s\r\n", "Content-Type", http.DetectContentType([]byte(w.response)))
	bu += "\r\n"
	bu += w.response
	return []byte(bu)
}
