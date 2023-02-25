package sysclient_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/ricoschulte/go-myapps/encryption"
	"github.com/ricoschulte/go-myapps/sysclient"
	"github.com/stretchr/testify/assert"
)

func TestResponseTypesToAdminMessages(t *testing.T) {
	sysclientpassword := "2jjH!u3ucXscEzHq8X!l83BX3!U8TPwA"
	secretkey := "to encrypt the local files"
	// create a dummy password file
	sysclientpassword_file, err := os.CreateTemp("", "sysclientpassword*.txt")
	if err != nil {
		panic(err)
	}
	defer os.Remove(sysclientpassword_file.Name()) // remove the file after the test is done

	// create a dummy admin password file
	administrativepassword_file, err := os.CreateTemp("", "sysclient_administrativepassword*.txt")
	if err != nil {
		panic(err)
	}
	defer os.Remove(administrativepassword_file.Name()) // remove the file after the test is done

	// reset the content after the test
	err = encryption.EncryptFileSha256AES256([]byte(secretkey), []byte(sysclientpassword), sysclientpassword_file.Name(), 0644)
	if err != nil {
		panic(err)
	}

	identity := sysclient.Identity{
		Id:      "f19033480af9",
		Product: "IP232",
		Version: "13r2 dvl [13.4250/131286/1300]",
		FwBuild: "134250",
		BcBuild: "131286",
		Major:   "13r2",
		Fw:      "ip222.bin",
		Bc:      "boot222.bin",
		Mini:    false,
		Platform: sysclient.Platform{
			Type: "PHONE",
		},
		EthIfs: []sysclient.EthIf{
			{
				If:   "ETH0",
				Ipv4: "172.16.4.141",
				Ipv6: "2002:91fd:9d07:0:290:33ff:fe46:af2",
			},
		},
	}

	tests := []struct {
		name string

		// message in
		TypeIn    byte
		CommandIn []byte
		DataIn    []byte

		// response message

		Type           byte
		Command        []byte
		Data           []byte
		expectNoRepose bool
		expectError    bool
	}{
		{
			name: "AdminReceiveChallenge to AdminSendIdentify with digest",

			TypeIn:    sysclient.MessageTypeAdmin,
			CommandIn: sysclient.AdminReceiveChallenge,
			DataIn:    []byte(`{"challenge":"1446931255"}`),

			Type:    sysclient.MessageTypeAdmin,
			Command: sysclient.AdminSendIdentify,
			Data:    []byte(string(`{"id":"f19033480af9","product":"IP232","version":"13r2 dvl [13.4250/131286/1300]","fwBuild":"134250","bcBuild":"131286","major":"13r2","fw":"ip222.bin","bc":"boot222.bin","mini":false,"pbxActive":false,"other":false,"platform":{"type":"PHONE","fxs":false},"digest":"13cc86b3f0691e2ce1a626a9963a5aefeefcb78f7e5b23a8ede10e8204247c97","ethIfs":[{"if":"ETH0","ipv4":"172.16.4.141","ipv6":"2002:91fd:9d07:0:290:33ff:fe46:af2"}]}`)),

			expectNoRepose: false,
			expectError:    false,
		},
		{
			name: "AdminReceiveSysclientPassword",

			TypeIn:    sysclient.MessageTypeAdmin,
			CommandIn: sysclient.AdminReceiveSysclientPassword,
			DataIn:    []byte(`{"password":"01234567890123456789012345678901"}`),

			Type:    0x00,
			Command: []byte{},
			Data:    []byte{},

			expectNoRepose: true,
			expectError:    false,
		},
		{
			name: "AdminReceiveNewAdminPassword",

			TypeIn:    sysclient.MessageTypeAdmin,
			CommandIn: sysclient.AdminReceiveNewAdminPassword,
			DataIn:    []byte(`{"seed":"0123456789012345", "pwd": "dummycontent"}`),

			Type:    0x00,
			Command: []byte{},
			Data:    []byte{},

			expectNoRepose: true,
			expectError:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// create the input / request message
			bam := []byte{test.TypeIn}
			bam = append(bam, test.CommandIn...)
			bam = append(bam, test.DataIn...)
			am, err_am := sysclient.NewAdminMessage(bam)
			assert.NoErrorf(t, err_am, "Error creating testmessage: %v", err_am)

			// create the sysclient
			sc, err_creating_client := sysclient.NewSysclient(
				identity,
				"wss://nourl",
				time.Duration(2*time.Second),
				true,
				http.NewServeMux(),
				sysclientpassword_file.Name(),
				administrativepassword_file.Name(),
				secretkey,
			)
			if err_creating_client != nil {
				t.Fatalf("Error creating client: %v", err_creating_client)
			}

			// send the test message to the client
			resp, err_handle := sc.HandleAdminMessage(am)
			if test.expectError {
				if err_handle == nil {
					t.Fatalf("a error is expected, but there was none")
				}
			} else {
				assert.NoError(t, err_handle, "error while handling the request")
				// testing the response
				if test.expectNoRepose {
					assert.Nil(t, resp, "the input message should not produce a response, but has.")
				} else {
					assert.NotNil(t, resp, "the response is nil")
					assert.NoErrorf(t, err_am, "Error handling request: %v", err_handle)
					assert.Equal(t, test.Type, resp.Type)
					assert.Equal(t, test.Command, resp.Command)
					assert.Equal(t, test.Data, resp.Data)

				}

			}
		})

	}

}
