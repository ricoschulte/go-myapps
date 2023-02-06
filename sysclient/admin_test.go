package sysclient_test

import (
	"testing"

	"github.com/ricoschulte/go-myapps/sysclient"
	"gotest.tools/assert"
)

func TestAdminMessage(t *testing.T) {
	tests := []struct {
		name string
		raw  []byte

		Type    byte
		Command []byte
		Data    []byte

		Error     bool
		ErrorText string
	}{
		{
			name:    "AdminSendIdentify",
			raw:     []byte{0x01, 0x00, 0x00, 0xff, 0x00, 0xff},
			Type:    0x01,
			Command: []byte{0x00, 0x00},
			Data:    []byte{0xff, 0x00, 0xff},

			Error:     false,
			ErrorText: "",
		},
		{
			name: "AdminReceiveSysclientPassword",
			raw:  []byte{0x01, 0x00, 0x01, 0xf1, 0xf2, 0xf3},

			Type:    0x01,
			Command: []byte{0x00, 0x01},
			Data:    []byte{0xf1, 0xf2, 0xf3},

			Error:     false,
			ErrorText: "",
		}, {
			name: "AdminReceiveChallenge",
			raw: []byte{0x01,
				0x00, 0x02,
				0xf1, 0xf2, 0xf3,
			},

			Type:    0x01,
			Command: []byte{0x00, 0x02},
			Data:    []byte{0xf1, 0xf2, 0xf3},

			Error:     false,
			ErrorText: "",
		}, {

			name: "AdminReceiveNewAdminPassword",
			raw: []byte{0x01,
				0x00, 0x03,
				0xf1, 0xf2, 0xf3,
			},

			Type:    0x01,
			Command: []byte{0x00, 0x03},
			Data:    []byte{0xf1, 0xf2, 0xf3},

			Error:     false,
			ErrorText: "",
		}, {
			name: "AdminProvision",
			raw: []byte{0x01,
				0x00, 0x04,
				0xf1, 0xf2, 0xf3,
			},

			Type:    0x01,
			Command: []byte{0x00, 0x04},
			Data:    []byte{0xf1, 0xf2, 0xf3},

			Error:     false,
			ErrorText: "",
		}, {
			name: "AdminProvisionResult",
			raw: []byte{0x01,
				0x00, 0x05,
				0xf1, 0xf2, 0xf3,
			},

			Type:    0x01,
			Command: []byte{0x00, 0x05},
			Data:    []byte{0xf1, 0xf2, 0xf3},

			Error:     false,
			ErrorText: "",
		}, {
			name: "AdminConfiguration",
			raw: []byte{0x01,
				0x00, 0x06,
				0xf1, 0xf2, 0xf3,
			},

			Type:      0x01,
			Command:   []byte{0x00, 0x06},
			Data:      []byte{0xf1, 0xf2, 0xf3},
			Error:     false,
			ErrorText: "",
		}, {
			name: "Not a AdminMessage",
			raw: []byte{0x02, // is a Tunnel message
				0x12, 0x34, 0x56, 0x78, 0x00, 0x00, 0x00, 0x03, 0xf1, 0xf2, 0xf3},

			Type:    0x01,
			Command: []byte{0x00, 0x00},
			Data:    []byte{0xff, 0x00, 0xff},

			Error:     true,
			ErrorText: "message is not a AdminMessage",
		}, {
			name: "Message with invalid length",
			raw:  []byte{0x01, 0x12 /* , one missing*/},

			Type:      0x01,
			Command:   []byte{0x00, 0x01},
			Data:      []byte{0xff, 0x00, 0xff},
			Error:     true,
			ErrorText: "invalid length of AdminMessage",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := sysclient.NewAdminMessage(test.raw)
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

				assert.Equal(t, 2, len(m.Command), "Command has invalid length")
				assert.DeepEqual(t, test.Command, m.Command)

				assert.Equal(t, len(test.Data), len(m.Data), "Data has invalid length")
				assert.DeepEqual(t, test.Data, m.Data)

				assert.DeepEqual(t, test.raw, m.AsBytes())
			}

		})
	}

}

func TestPasswordFromJson(t *testing.T) {
	tests := []struct {
		name string
		raw  string

		Password  string
		Error     bool
		ErrorText string
	}{
		{
			"a",
			`{"password":"AEF_06eOMcDANh5DKkklhVjKJ3lgVXVD"}`,
			"AEF_06eOMcDANh5DKkklhVjKJ3lgVXVD",
			false,
			"",
		},
		{
			"invalid json",
			`{"password":"AEF_06eOMcDANh5DKkklhVjKJ3lgVXVD"`,
			"AEF_06eOMcDANh5DKkklhVjKJ3lgVXVD",
			true,
			"unexpected end of JSON input",
		},
		{
			"empty json",
			``,
			"AEF_06eOMcDANh5DKkklhVjKJ3lgVXVD",
			true,
			"unexpected end of JSON input",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := sysclient.NewPassword([]byte(test.raw))
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

				assert.Equal(t, test.Password, m.Password)
			}

		})
	}

}

func TestChallengeFromJson(t *testing.T) {
	tests := []struct {
		name string
		raw  string

		Challenge string
		Error     bool
		ErrorText string
	}{
		{
			"ok",
			`{"challenge":"316668907"}`,
			"316668907",
			false,
			"",
		},
		{
			"invalid json",
			`{"challenge:"316668907"}`,
			"316668907",
			true,
			"invalid character '3' after object key",
		},
		{
			"empty bytes",
			``,
			"316668907",
			true,
			"unexpected end of JSON input",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := sysclient.NewChallenge([]byte(test.raw))
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

				assert.Equal(t, test.Challenge, m.Challenge)
			}

		})
	}

}

func TestAdministrativePasswordFromJson(t *testing.T) {
	tests := []struct {
		name string
		raw  string

		Seed      string
		Pwd       string
		Error     bool
		ErrorText string
	}{
		{
			"ok",
			`{"seed":"somethingsomething","pwd":"othersomethingsomething"}`,
			"somethingsomething",
			"othersomethingsomething",
			false,
			"",
		},
		{
			"invalid json",
			`{"seed":"somethingsomething""pwd":"othersomethingsomething"}`,
			"somethingsomething",
			"othersomethingsomething",
			true,
			`invalid character '"' after object key:value pair`,
		},
		{
			"empty bytes",
			``,
			"somethingsomething",
			"othersomethingsomething",
			true,
			"unexpected end of JSON input",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			m, err := sysclient.NewAdministrativePassword([]byte(test.raw))
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

				assert.Equal(t, test.Seed, m.Seed)
				assert.Equal(t, test.Pwd, m.Pwd)
			}

		})
	}

}

func TestDecryptAdminPassword(t *testing.T) {
	tests := []struct {
		name                  string
		Seed                  string
		Pwd                   string
		SysclientPassword     string
		expectedAdminPassword string
		Error                 bool
		ErrorText             string
	}{
		// #TODO it is unclear how the decrypted content should look like.
		// at the moment its everything, but not 'test'.
		// i dont know if the decrypted value is still encrypted in a way, a device can be configured with it.
		
		// {
		// 	name:                  "a",
		// 	Seed:                  "wJQog70GUK!GPTN2",
		// 	Pwd:                   "22f91ebb",
		// 	SysclientPassword:     "7iICN04agQq38j!hmbN7P5LLI1YePz1w",
		// 	expectedAdminPassword: "test",  <----- problem
		// 	Error:                 false,
		// 	ErrorText:             "",
		// },
		// {
		// 	name:                  "b",
		// 	Seed:                  "RlGIgQj2j!LUIgPL",
		// 	Pwd:                   "cfa78653",
		// 	SysclientPassword:     "N803pDZVpL_ccS_1tGwFjIN0cfw75MNo",
		// 	expectedAdminPassword: "test",  <----- problem
		// 	Error:                 false,
		// 	ErrorText:             "",
		// },
		{
			name: "empty seed",
			//Seed:                  "1234567812345678",
			Pwd:                   "sdfsdf",
			SysclientPassword:     "12345678123456781234567812345678",
			expectedAdminPassword: "sdfsdf",
			Error:                 true,
			ErrorText:             "seed cant be empty",
		},
		{
			name:                  "invalid seed length",
			Seed:                  "1234567812345678_",
			Pwd:                   "sdfsdf",
			SysclientPassword:     "12345678123456781234567812345678",
			expectedAdminPassword: "sdfsdf",
			Error:                 true,
			ErrorText:             "seed has a invalid length of '17'. should have a length of 16",
		},
		{
			name: "empty pwd",
			Seed: "1234567812345678",
			//Pwd:                   "sdfsdf",
			SysclientPassword:     "12345678123456781234567812345678",
			expectedAdminPassword: "sdfswrzw34",
			Error:                 true,
			ErrorText:             "pwd cant be empty",
		},
		{
			name: "empty sysclient",
			Seed: "1234567812345678",
			Pwd:  "sdfsdf",
			//SysclientPassword:     "12345678123456781234567812345678",
			expectedAdminPassword: "defdsshh",
			Error:                 true,
			ErrorText:             "sysclientPassword cant be empty",
		},
		{
			name:                  "invalid sysclientpassword length",
			Seed:                  "1234567812345678",
			Pwd:                   "sdfsdf",
			SysclientPassword:     "1234567812345678123__4567812345678",
			expectedAdminPassword: "defdsshh",
			Error:                 true,
			ErrorText:             "sysclientPassword has a invalid length of '34'. should have a length of 32",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bin := []byte{sysclient.MessageTypeAdmin}
			bin = append(bin, sysclient.AdminReceiveNewAdminPassword...)
			m, err_am := sysclient.NewAdminMessage(bin)
			if err_am != nil {
				t.Fatal("error creating test AM instance")
			}

			adminpwd_bytes, err := m.DecryptAdminPassword(test.SysclientPassword, test.Seed, test.Pwd)
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

				assert.Equal(t, test.expectedAdminPassword, string(adminpwd_bytes))
			}

		})
	}

}

func TestGetLoginDigest(t *testing.T) {

	tests := []struct {
		name           string
		id             string
		product        string
		version        string
		challenge      string
		password       string
		expectedDigest string
		Error          bool
		ErrorText      string
	}{
		{
			name:           "a",
			id:             "f19033480af9",
			product:        "IP232",
			version:        "13r2 dvl [13.4250/131286/1300]",
			challenge:      "1446931255",
			password:       "2jjH!u3ucXscEzHq8X!l83BX3!U8TPwA",
			expectedDigest: "13cc86b3f0691e2ce1a626a9963a5aefeefcb78f7e5b23a8ede10e8204247c97",
			Error:          false,
			ErrorText:      "",
		},
		{
			name:           "b",
			id:             "f19033480af9",
			product:        "IP232",
			version:        "13r2 dvl [13.4250/131286/1300]",
			challenge:      "616692991",
			password:       "2jjH!u3ucXscEzHq8X!l83BX3!U8TPwA",
			expectedDigest: "dcbf704513ba5f771aabc5ba553ab82928c51cd69ef1c010174086f5b82414a8",
			Error:          false,
			ErrorText:      "",
		},
		{
			name:           "c",
			id:             "f19033480af9",
			product:        "IP222",
			version:        "13r2 dvl [13.4250/131286/1300]",
			challenge:      "931977739",
			password:       "q19rykt45KceQbEyN8rXPLQlb4VZRvn7",
			expectedDigest: "75c2d186d072f9a22e797c4c1978d8f2f00185c2cce28d02515617368aec4ed4",
			Error:          false,
			ErrorText:      "",
		},
		{
			name:           "d",
			id:             "f19033480af9",
			product:        "IP222",
			version:        "13r2 dvl [13.4250/131286/1300]",
			challenge:      "201737720",
			password:       "q19rykt45KceQbEyN8rXPLQlb4VZRvn7",
			expectedDigest: "939f89c2fe801cf4a2217b3db94e1f9c1e24544b03a303f716e9c60878eff9d6",
			Error:          false,
			ErrorText:      "",
		},
		{
			name:           "e",
			id:             "f19033480af9",
			product:        "IP222",
			version:        "13r2 dvl [13.4250/131286/1300]",
			challenge:      "571341503",
			password:       "boO!fTk!YepvjmYPS6ejgTVzh6DrD6vo",
			expectedDigest: "d8c5c3d118890df9a9439d40941f2e6be136942374a4736cd687ab4ea96f8f90",
			Error:          false,
			ErrorText:      "",
		},
		{
			name:           "f",
			id:             "f19033480af9",
			product:        "IP222",
			version:        "13r2 dvl [13.4250/131286/1300]",
			challenge:      "1404833035",
			password:       "boO!fTk!YepvjmYPS6ejgTVzh6DrD6vo",
			expectedDigest: "e8ce86a0395ad01d09dc08ca44195f3613bfb3d8407d9e9a9f2d15907be9f351",
			Error:          false,
			ErrorText:      "",
		},

		{
			name: "empty id",
			//id:             "f09033480af9",
			product:        "IP222",
			version:        "13r2 dvl [13.4250/131286/1300]",
			challenge:      "sdsdfsd",
			password:       "sdfsdf",
			expectedDigest: "",
			Error:          true,
			ErrorText:      "id cant be empty",
		},
		{
			name: "empty product",
			id:   "f09033480af9",
			//product:        "IP222",
			version:        "13r2 dvl [13.4250/131286/1300]",
			challenge:      "sdsdfsd",
			password:       "sdfsdf",
			expectedDigest: "",
			Error:          true,
			ErrorText:      "product cant be empty",
		},
		{
			name:    "empty version",
			id:      "f09033480af9",
			product: "IP222",
			//version:        "13r2 dvl [13.4250/131286/1300]",
			challenge:      "sdsdfsd",
			password:       "sdfsdf",
			expectedDigest: "",
			Error:          true,
			ErrorText:      "version cant be empty",
		},
		{
			name:    "empty challenge",
			id:      "f09033480af9",
			product: "IP222",
			version: "13r2 dvl [13.4250/131286/1300]",
			//challenge:      "sdsdfsd",
			password:       "sdfsdf",
			expectedDigest: "",
			Error:          true,
			ErrorText:      "challenge cant be empty",
		},
		{
			name:      "empty password",
			id:        "f09033480af9",
			product:   "IP222",
			version:   "13r2 dvl [13.4250/131286/1300]",
			challenge: "sdsdfsd",
			//password:       "sdfsdf",
			expectedDigest: "",
			Error:          true,
			ErrorText:      "password cant be empty",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bin := []byte{sysclient.MessageTypeAdmin}
			bin = append(bin, sysclient.AdminReceiveNewAdminPassword...)
			m, err_am := sysclient.NewAdminMessage(bin)
			if err_am != nil {
				t.Fatal("error creating test AM instance")
			}

			digest_string, err := m.GetLoginDigest(
				test.id,
				test.product,
				test.version,
				test.challenge,
				test.password,
			)

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
				assert.Equal(t, test.expectedDigest, string(digest_string))
			}
		})
	}
}
