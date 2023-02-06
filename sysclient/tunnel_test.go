package sysclient_test

import (
	"testing"

	"github.com/ricoschulte/go-myapps/sysclient"
	"gotest.tools/assert"
)

func TestTunnelMessage(t *testing.T) {
	tests := []struct {
		name string
		raw  []byte

		Type           byte
		SessionId      []byte
		SessionIdInt32 int32
		EventType      []byte
		Data           []byte

		Error     bool
		ErrorText string
	}{
		{
			name: "Dummy",
			raw:  []byte{0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0x00, 0xff},

			Type:           0x02,
			SessionId:      []byte{0x00, 0x00, 0x00, 0x00},
			SessionIdInt32: 0,
			EventType:      []byte{0x00, 0x00, 0x00, 0x00},
			Data:           []byte{0xff, 0x00, 0xff},

			Error:     false,
			ErrorText: "",
		},
		{
			name: "b",
			raw:  []byte{0x02, 0x12, 0x34, 0x56, 0x78, 0xa1, 0xa2, 0xa3, 0xa4, 0xf1, 0xf2, 0xf3},

			Type:           0x02,
			SessionId:      []byte{0x12, 0x34, 0x56, 0x78},
			SessionIdInt32: 305419896,
			EventType:      []byte{0xa1, 0xa2, 0xa3, 0xa4},
			Data:           []byte{0xf1, 0xf2, 0xf3},

			Error:     false,
			ErrorText: "",
		}, {
			name: "TunnelSend",
			raw: []byte{0x02, 0x12, 0x34, 0x56, 0x78,
				0x00, 0x00, 0x00, 0x00,
				0xf1, 0xf2, 0xf3},

			Type:           0x02,
			SessionId:      []byte{0x12, 0x34, 0x56, 0x78},
			SessionIdInt32: 305419896,
			EventType:      sysclient.TunnelSend,
			Data:           []byte{0xf1, 0xf2, 0xf3},

			Error:     false,
			ErrorText: "",
		}, {

			name: "TunnelReceive",
			raw: []byte{0x02, 0x12, 0x34, 0x56, 0x78,
				0x00, 0x00, 0x00, 0x01,
				0xf1, 0xf2, 0xf3},

			Type:           0x02,
			SessionId:      []byte{0x12, 0x34, 0x56, 0x78},
			SessionIdInt32: 305419896,
			EventType:      sysclient.TunnelReceive,
			Data:           []byte{0xf1, 0xf2, 0xf3},

			Error:     false,
			ErrorText: "",
		}, {
			name: "TunnelReceiveResult",
			raw: []byte{0x02,
				0x12, 0x34, 0x56, 0x78,
				0x00, 0x00, 0x00, 0x02,
				0xf1, 0xf2, 0xf3},

			Type:           0x02,
			SessionId:      []byte{0x12, 0x34, 0x56, 0x78},
			SessionIdInt32: 305419896,
			EventType:      sysclient.TunnelReceiveResult,
			Data:           []byte{0xf1, 0xf2, 0xf3},

			Error:     false,
			ErrorText: "",
		}, {
			name: "TunnelSendResult",
			raw: []byte{0x02, 0x12, 0x34, 0x56, 0x78,
				0x00, 0x00, 0x00, 0x03,
				0xf1, 0xf2, 0xf3},

			Type:           0x02,
			SessionId:      []byte{0x12, 0x34, 0x56, 0x78},
			SessionIdInt32: 305419896,
			EventType:      sysclient.TunnelSendResult,
			Data:           []byte{0xf1, 0xf2, 0xf3},

			Error:     false,
			ErrorText: "",
		}, {
			name: "TunnelShutdown",
			raw: []byte{0x02, 0xff, 0xff, 0xff, 0xff,
				0x00, 0x00, 0x00, 0x04,
				0xf1, 0xf2, 0xf3},

			Type:           0x02,
			SessionId:      []byte{0xff, 0xff, 0xff, 0xff},
			SessionIdInt32: -1,
			EventType:      sysclient.TunnelShutdown,
			Data:           []byte{0xf1, 0xf2, 0xf3},

			Error:     false,
			ErrorText: "",
		}, {
			name: "Not a Tunnel Message",
			raw: []byte{0x01, // is a Admin message
				0x12, 0x34, 0x56, 0x78, 0x00, 0x00, 0x00, 0x03, 0xf1, 0xf2, 0xf3},

			Type:           0x02,
			SessionId:      []byte{0x12, 0x34, 0x56, 0x78},
			SessionIdInt32: 305419896,
			EventType:      sysclient.TunnelShutdown,
			Data:           []byte{0xf1, 0xf2, 0xf3},

			Error:     true,
			ErrorText: "message is not a TunnelMessage",
		}, {
			name: "Message with invalid length",
			raw:  []byte{0x02, 0x12, 0x34, 0x56, 0x78, 0x00, 0x00, 0x00 /* , one missing*/},

			Type:           0x00,
			SessionId:      nil,
			SessionIdInt32: 0,
			EventType:      nil,
			Data:           nil,

			Error:     true,
			ErrorText: "invalid length of message",
			// }, {
			// 	name: "Message requesting Response but has no data",
			// 	raw:  []byte{0x02, 0x12, 0x34, 0x56, 0x78, 0x00, 0x00, 0x00, 0x00 /* , one missing*/},

			// 	Type:      0x00,
			// 	SessionId: nil,
			// 	EventType: nil,
			// 	Data:      nil,

			// 	Error:     true,
			// 	ErrorText: "invalid length of data",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := sysclient.NewTunnelMessage(test.raw)
			if test.Error {
				// error expected
				if err == nil {
					t.Fatal("returned no error")
				}
				assert.Equal(t, test.ErrorText, err.Error(), "invalid ErrorText")
			} else {
				// no error expected
				if err != nil {
					t.Fatalf("returned error: %v", err)
				}

				assert.Equal(t, test.Type, m.Type)

				assert.Equal(t, 4, len(m.SessionId), "SessionId has invalid length")
				assert.DeepEqual(t, test.SessionId, m.SessionId)

				IdInt32, err_session_int32 := m.GetSessionId()
				assert.NilError(t, err_session_int32, "returned a error while returning SessionId as in64")
				assert.Equal(t, test.SessionIdInt32, IdInt32, "invalid int32 SessionId")

				assert.Equal(t, 4, len(m.EventType), "EventType has invalid length")
				assert.DeepEqual(t, test.EventType, m.EventType)

				assert.Equal(t, len(test.Data), len(m.Data), "Data has invalid length")
				assert.DeepEqual(t, test.Data, m.Data)

				assert.DeepEqual(t, test.raw, m.AsBytes())
			}

		})
	}

}
