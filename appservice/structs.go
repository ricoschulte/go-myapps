package appservice

import "github.com/ricoschulte/go-myapps/connection"

type AppChallengeResult struct {
	Mt        string `json:"mt"`
	Challenge string `json:"challenge"`
}

type AppGetLogin struct {
	Mt        string `json:"mt"`
	Src       string `json:"src"`
	App       string `json:"app"`
	Challenge string `json:"challenge"`
}

func NewAppGetLogin(app, challenge string) AppGetLogin {
	return AppGetLogin{
		Mt:        "AppGetLogin",
		Src:       connection.GetRandomHexString(10),
		App:       app,
		Challenge: challenge,
	}
}

type Info struct {
	Appobj     string   `json:"appobj"`
	Appdn      string   `json:"appdn"`
	Appurl     string   `json:"appurl"`
	Pbx        string   `json:"pbx"`
	Cn         string   `json:"cn"`
	Unlicensed bool     `json:"unlicensed,omitempty"`
	Testmode   bool     `json:"testmode,omitempty"`
	Groups     []string `json:"groups"`
	Apps       []struct {
		Name string `json:"name"`
		Dn   string `json:"dn"`
	} `json:"apps"`
}

type AppGetLoginResult struct {
	Mt     string `json:"mt"`
	Src    string `json:"src"`
	Sip    string `json:"sip"`
	Guid   string `json:"guid"`
	Dn     string `json:"dn"`
	PbxObj string `json:"pbxObj"`
	Domain string `json:"domain"`
	App    string `json:"app"`
	Info   Info   `json:"info"`
	Digest string `json:"digest"`
	Salt   string `json:"salt"` // is this used again?
	Key    string `json:"key"`  // is this used again?
}

type AppLogin struct {
	Mt     string `json:"mt"`
	Sip    string `json:"sip"`
	Guid   string `json:"guid"`
	Dn     string `json:"dn"`
	PbxObj string `json:"pbxObj"`
	Domain string `json:"domain"`
	App    string `json:"app"`
	Info   Info   `json:"info"`
	Digest string `json:"digest"`
}

func NewAppLogin(apploginfrompbx AppGetLoginResult) AppLogin {

	return AppLogin{
		Mt:     "AppLogin",
		Sip:    apploginfrompbx.Sip,
		Guid:   apploginfrompbx.Guid,
		Dn:     apploginfrompbx.Dn,
		PbxObj: apploginfrompbx.PbxObj,
		Domain: apploginfrompbx.Domain,
		App:    apploginfrompbx.App,
		Info:   apploginfrompbx.Info,
		Digest: apploginfrompbx.Digest,
	}
}

type AppLoginResult struct {
	Mt  string `json:"mt"`
	App string `json:"app"`
	Ok  bool   `json:"ok"`
}
