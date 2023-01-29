package connection_test

import (
	"encoding/json"
	"testing"

	"github.com/ricoschulte/go-myapps/connection"
	"github.com/stretchr/testify/assert"
)

func TestParsingCheckBuildResult(t *testing.T) {

	tests := []struct {
		name                string
		mt                  string
		version             string
		build               string
		launcherUpdateBuild string
		appStoreUrl         string
		raw                 string
	}{
		{"Acloud", "CheckBuildResult", "13r3", "137766", "137766", "https://store.innovaphone.com/release/download/",
			`{"mt":"CheckBuildResult","version":"13r3","build":"137766","launcherUpdateBuild":"137766","appStoreUrl":"https://store.innovaphone.com/release/download/"}`},
		{"Bfritz",
			"CheckBuildResult", "", "137656", "137656", "https://store.innovaphone.com/beta/download/",
			`{"mt":"CheckBuildResult","build":"137656","launcherUpdateBuild":"137656","appStoreUrl":"https://store.innovaphone.com/beta/download/"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var parsed connection.CheckBuildResult
			json.Unmarshal([]byte(test.raw), &parsed)

			assert.Equal(t, test.mt, parsed.Mt)
			assert.Equal(t, test.version, parsed.Version)
			assert.Equal(t, test.build, parsed.Build)
			assert.Equal(t, test.launcherUpdateBuild, parsed.LauncherUpdateBuild)
			assert.Equal(t, test.appStoreUrl, parsed.AppStoreUrl)

		})
	}

}

func TestParsingUpdateRegister(t *testing.T) {

	tests := []struct {
		name     string
		Mt       string
		Profile  string
		Tutorial string
		Signup   string
		Reset    string
		raw      string
	}{
		{"A",
			"UpdateRegister", "profile", "tutorial", "", "",
			`{"mt":"UpdateRegister","profile":"profile","tutorial":"tutorial"}`,
		},
		{"B",
			"UpdateRegister", "profile", "tutorials", "https://apps.fritz.box/fritz.box/usersapp/register.htm", "https://apps.fritz.box/fritz.box/usersapp/password.htm",
			`{"mt":"UpdateRegister","signup":"https://apps.fritz.box/fritz.box/usersapp/register.htm","reset":"https://apps.fritz.box/fritz.box/usersapp/password.htm","profile":"profile","tutorial":"tutorials"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var parsed connection.UpdateRegister
			json.Unmarshal([]byte(test.raw), &parsed)

			assert.Equal(t, test.Mt, parsed.Mt)
			assert.Equal(t, test.Profile, parsed.Profile)
			assert.Equal(t, test.Tutorial, parsed.Tutorial)
			assert.Equal(t, test.Signup, parsed.Signup)
			assert.Equal(t, test.Reset, parsed.Reset)

		})
	}

}

func TestParsingLoginInfoResult(t *testing.T) {

	tests := []struct {
		name           string
		Mt             string
		UserDigest     bool
		UserNtlm       bool
		UserOAuth2     bool
		UserOAuth2Name string

		SessionDigest bool

		raw string
	}{
		{"A",
			"LoginInfoResult",

			true,
			false,
			false,
			"auth OAuth2",
			true,
			`{"mt":"LoginInfoResult","user":{
				"digest":true,"ntlm":false,"oauth2":false,"oauth2Name":"auth OAuth2"},
				"session":{"digest":true}
			}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var parsed connection.LoginInfoResult
			json.Unmarshal([]byte(test.raw), &parsed)

			assert.Equal(t, test.Mt, parsed.Mt)
			assert.Equal(t, test.UserDigest, parsed.User.Digest)
			assert.Equal(t, test.UserNtlm, parsed.User.Ntlm)
			assert.Equal(t, test.UserOAuth2, parsed.User.OAuth2)
			assert.Equal(t, test.UserOAuth2Name, parsed.User.OAuth2Name)
			assert.Equal(t, test.SessionDigest, parsed.Session.Digest)

		})
	}

}

func TestParsingLogin(t *testing.T) {
	tests := []struct {
		name      string
		Mt        string
		Type      string
		UserAgent string

		Method   string
		Username string
		Nonce    string
		Response string

		raw string
	}{
		{"session without digest", "Login", "session", "myApps (Firefox)", "", "", "", "", `{"mt":"Login","type":"session","userAgent":"myApps (Firefox)"}`},
		{"session with digest", "Login", "session", "myApps (Firefox)", "digest", "c997302cfd92630121fd00505600a2cd", "749334693e72d301", "3684cb22dfdfea4bf03dcdceabd864053c7ec66767a8d996b06f8b94afb9307b",
			`{"mt":"Login","type":"session","method":"digest",
		"username":"c997302cfd92630121fd00505600a2cd","nonce":"749334693e72d301",
		"response":"3684cb22dfdfea4bf03dcdceabd864053c7ec66767a8d996b06f8b94afb9307b","userAgent":"myApps (Firefox)"}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var parsed connection.Login
			json.Unmarshal([]byte(test.raw), &parsed)

			assert.Equal(t, test.Mt, parsed.Mt)
			assert.Equal(t, test.Type, parsed.Type)
			assert.Equal(t, test.UserAgent, parsed.UserAgent)

			assert.Equal(t, test.Method, parsed.Method)
			assert.Equal(t, test.Username, parsed.Username)
			assert.Equal(t, test.Nonce, parsed.Nonce)
			assert.Equal(t, test.Response, parsed.Response)

		})
	}
}

func TestParsingAuthenticate(t *testing.T) {
	tests := []struct {
		name      string
		Mt        string
		Type      string
		Method    string
		Domain    string
		Challenge string
		raw       string
	}{
		{
			"A",
			"Authenticate",
			"session",
			"digest",
			"pbx.company.com",
			"84463cc884463c89",
			`{"mt":"Authenticate","type":"session","method":"digest","domain":"pbx.company.com","challenge":"84463cc884463c89"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var parsed connection.Authenticate
			json.Unmarshal([]byte(test.raw), &parsed)

			assert.Equal(t, test.Mt, parsed.Mt)
			assert.Equal(t, test.Type, parsed.Type)
			assert.Equal(t, test.Method, parsed.Method)
			assert.Equal(t, test.Domain, parsed.Domain)
			assert.Equal(t, test.Challenge, parsed.Challenge)
		})
	}
}

func TestParsingLoginResult(t *testing.T) {
	tests := []struct {
		Name string
		Mt   string

		Error     int
		ErrorText string

		InfoUserDomain        string
		InfoUserSip           string
		InfoUserGuid          string
		InfoUserDn            string
		InfoUserNum           string
		InfoUserEmail         string
		InfoUserPrefixIntl    string
		InfoUserPrefixNtl     string
		InfoUserPrefixSubs    string
		InfoUserPrefixArea    string
		InfoUserPrefixCountry string
		InfoAccount           map[string]interface{}
		InfoSettings          map[string]interface{}
		InfoSessionUsr        string
		InfoSessionPwd        string
		InfoAltHttp           string
		InfoAlt               string

		ExpectedDigest string
		ExpectedRaw    string
	}{
		{
			Name:                  "A",
			Mt:                    "LoginResult",
			Error:                 1,
			ErrorText:             "",
			InfoUserDomain:        "company.com",
			InfoUserSip:           "usersip",
			InfoUserGuid:          "2c4a96f920ad5e01241c00505600a2cd",
			InfoUserDn:            "User Name",
			InfoUserNum:           "202",
			InfoUserEmail:         "usere@email.de",
			InfoUserPrefixIntl:    "",
			InfoUserPrefixNtl:     "",
			InfoUserPrefixSubs:    "",
			InfoUserPrefixArea:    "",
			InfoUserPrefixCountry: "",
			InfoAccount:           map[string]interface{}{},
			InfoSettings:          map[string]interface{}{},
			InfoSessionUsr:        "",
			InfoSessionPwd:        "",
			InfoAltHttp:           "company.com/PBX0",
			InfoAlt:               "alt.company.com/PBX0",
			ExpectedDigest:        "eb9f658f9128ad4f78bd13dae5faaa59b114d092862943f7fe219096c0726e37",
			ExpectedRaw: `{"mt":"LoginResult","info":{
				"user":{
					"domain":"company.com","sip":"usersip","guid":"2c4a96f920ad5e01241c00505600a2cd","dn":"User Name","num":"202","email":"usere@email.de","prefix":{}
				},
				"account":{},
				"settings":{},
				"altHttp":"company.com/PBX0","alt":"alt.company.com/PBX0"
			},"digest":"eb9f658f9128ad4f78bd13dae5faaa59b114d092862943f7fe219096c0726e37"}`,
		},
		{
			Name:                  "B",
			Mt:                    "LoginResult",
			Error:                 1,
			ErrorText:             "",
			InfoUserDomain:        "example.com",
			InfoUserSip:           "usersip",
			InfoUserGuid:          "2c4a96f920ad5e01241c00505600a2cd",
			InfoUserDn:            "User Name äöüß",
			InfoUserNum:           "202",
			InfoUserEmail:         "usere@example.de",
			InfoUserPrefixIntl:    "000",
			InfoUserPrefixNtl:     "00",
			InfoUserPrefixSubs:    "0",
			InfoUserPrefixArea:    "40",
			InfoUserPrefixCountry: "49",
			InfoAccount:           map[string]interface{}{},
			InfoSettings:          map[string]interface{}{},
			InfoSessionUsr:        "c53a54119e819cc2cc4ce8e385ad462e9162b516f7292cc071aac075fe292c33",
			InfoSessionPwd:        "eee03a3752e6801c5ab14d2735edca9a7d5fdea0440c8a",
			InfoAltHttp:           "example.com/PBX0",
			InfoAlt:               "alt.example.com/PBX0",
			ExpectedDigest:        "eb9f658f9128ad4f78bd13dae5faaa59b114d092862943f7fe219096c0726e37",
			ExpectedRaw: `{"mt":"LoginResult","info":{
				"user":{
					"domain":"example.com","sip":"usersip","guid":"2c4a96f920ad5e01241c00505600a2cd","dn":"User Name äöüß","num":"202","email":"usere@example.de",
					"prefix":{"intl":"000","ntl":"00","subs":"0","area":"40","country":"49"}
				},
				"account":{},
				"settings":{},
				"session":{"usr":"c53a54119e819cc2cc4ce8e385ad462e9162b516f7292cc071aac075fe292c33","pwd":"eee03a3752e6801c5ab14d2735edca9a7d5fdea0440c8a"},
				"altHttp":"example.com/PBX0","alt":"alt.example.com/PBX0"
			},"digest":"eb9f658f9128ad4f78bd13dae5faaa59b114d092862943f7fe219096c0726e37"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			var parsed connection.LoginResult
			err := json.Unmarshal([]byte(test.ExpectedRaw), &parsed)
			if err != nil {
				t.Fatalf("Error while parsing json: %v", err)
			}
			assert.Equal(t, test.Mt, parsed.Mt)
			assert.Equal(t, test.InfoUserDomain, parsed.Info.User.Domain)
			assert.Equal(t, test.InfoUserSip, parsed.Info.User.Sip)
			assert.Equal(t, test.InfoUserGuid, parsed.Info.User.Guid)
			assert.Equal(t, test.InfoUserDn, parsed.Info.User.Dn)
			assert.Equal(t, test.InfoUserNum, parsed.Info.User.Num)
			assert.Equal(t, test.InfoUserEmail, parsed.Info.User.Email)
			assert.Equal(t, test.InfoAccount, parsed.Info.Account)
			assert.Equal(t, test.InfoUserPrefixIntl, parsed.Info.User.Prefix.Intl)
			assert.Equal(t, test.InfoUserPrefixNtl, parsed.Info.User.Prefix.Ntl)
			assert.Equal(t, test.InfoUserPrefixSubs, parsed.Info.User.Prefix.Subs)
			assert.Equal(t, test.InfoUserPrefixArea, parsed.Info.User.Prefix.Area)
			assert.Equal(t, test.InfoUserPrefixCountry, parsed.Info.User.Prefix.Country)
			assert.Equal(t, test.InfoSettings, parsed.Info.Settings)
			assert.Equal(t, test.InfoSessionUsr, parsed.Info.Session.Usr)
			assert.Equal(t, test.InfoSessionPwd, parsed.Info.Session.Pwd)
			assert.Equal(t, test.InfoAltHttp, parsed.Info.AltHttp)
			assert.Equal(t, test.InfoAlt, parsed.Info.Alt)
			assert.Equal(t, test.ExpectedDigest, parsed.Digest)
		})
	}
}

func TestParsingSessionAdded(t *testing.T) {
	tests := []struct {
		name      string
		Mt        string
		ID        string
		UserAgent string
		Timestamp int64
		raw       string
	}{
		{
			"A",
			"SessionAdded",
			"4b9c710caaca63015649000c2914757a",
			"myApps (Chrome)",
			1674226510000,
			`{"mt":"SessionAdded","id":"4b9c710caaca63015649000c2914757a","info":{"userAgent":"myApps (Chrome)","timestamp":1674226510000}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var parsed connection.SessionAdded
			json.Unmarshal([]byte(test.raw), &parsed)

			assert.Equal(t, test.Mt, parsed.Mt)
			assert.Equal(t, test.ID, parsed.Id)
			assert.Equal(t, test.UserAgent, parsed.Info.UserAgent)
			assert.Equal(t, test.Timestamp, parsed.Info.Timestamp)
		})
	}
}

func TestParsingUpdateOwnPresence(t *testing.T) {
	tests := []struct {
		name     string
		Mt       string
		Presence []struct {
			Contact  string
			Activity string
			Status   string
		}
		raw string
	}{
		{"A", "UpdateOwnPresence", []struct {
			Contact  string
			Activity string
			Status   string
		}{
			{"tel:", "busy", "closed"},
			{"im:", "", "open"},
		}, `{"mt":"UpdateOwnPresence","presence":[{"contact":"tel:","activity":"busy","status":"closed"},{"contact":"im:","activity":"","status":"open"}]}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var parsed connection.UpdateOwnPresence
			json.Unmarshal([]byte(test.raw), &parsed)

			assert.Equal(t, test.Mt, parsed.Mt)
			assert.Equal(t, 2, len(parsed.Presence))
			assert.Equal(t, test.Presence[0].Contact, parsed.Presence[0].Contact)
			assert.Equal(t, test.Presence[0].Activity, parsed.Presence[0].Activity)
			assert.Equal(t, test.Presence[0].Status, parsed.Presence[0].Status)
			assert.Equal(t, test.Presence[1].Contact, parsed.Presence[1].Contact)
			assert.Equal(t, test.Presence[1].Activity, parsed.Presence[1].Activity)
			assert.Equal(t, test.Presence[1].Status, parsed.Presence[1].Status)
		})
	}
}

func TestParsingAuthorize(t *testing.T) {

	tests := []struct {
		name string
		mt   string
		code int
		raw  string
	}{
		{"A", "Authorize", 2001, `{"mt":"Authorize","code":2001}`},
		{"B", "Authorize", 2031, `{"mt":"Authorize","code":2031}`},
		{"C", "Authorize", 2011, `{"mt":"Authorize","code":2011}`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var parsed connection.Authorize
			json.Unmarshal([]byte(test.raw), &parsed)

			assert.Equal(t, test.mt, parsed.Mt)
			assert.Equal(t, test.code, parsed.Code)

		})
	}

}

func TestParsingRedirect(t *testing.T) {
	//

	tests := []struct {
		name           string
		mt             string
		infoHost       string
		infoMod        string
		InfoSessionUsr string
		InfoSessionPwd string
		digest         string
		raw            string
	}{
		{"A",
			"Redirect",
			"secondary1.example.de",
			"PBX0",
			"6020b754d37061acb8b123487f2fe69ab6f2232f1240683d65e2e193532196c0",
			"09e4274e110647721671234d2f6c1c546b52ff4f6513c2",
			"ac1c25f5e0e1cafdb61c31412345811b52d290e51e5f2d1021022ecbd49b8e7d",
			`{
				"mt":"Redirect",
				"info":{
					"host":"secondary1.example.de",
					"mod":"PBX0",
					"session":{
						"usr":"6020b754d37061acb8b123487f2fe69ab6f2232f1240683d65e2e193532196c0",
						"pwd":"09e4274e110647721671234d2f6c1c546b52ff4f6513c2"
					}
				},
				"digest":"ac1c25f5e0e1cafdb61c31412345811b52d290e51e5f2d1021022ecbd49b8e7d"
			}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var parsed connection.Redirect
			json.Unmarshal([]byte(test.raw), &parsed)

			assert.Equal(t, test.mt, parsed.Mt)
			assert.Equal(t, test.infoHost, parsed.Info.Host)
			assert.Equal(t, test.infoMod, parsed.Info.Mod)
			assert.Equal(t, test.InfoSessionUsr, parsed.Info.Session.Usr)
			assert.Equal(t, test.InfoSessionPwd, parsed.Info.Session.Pwd)
			assert.Equal(t, test.digest, parsed.Digest)

		})
	}

}
