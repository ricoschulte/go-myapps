package connection

type CheckBuildResult struct {
	Mt                  string `json:"mt"`
	Version             string `json:"version"`
	Build               string `json:"build"`
	LauncherUpdateBuild string `json:"launcherUpdateBuild"`
	AppStoreUrl         string `json:"appStoreUrl"`
}

type SubscribeRegister struct {
	MT string `json:"mt"`
}

type UpdateRegister struct {
	Mt       string `json:"mt"`
	Profile  string `json:"profile"`
	Tutorial string `json:"tutorial"`
	Signup   string `json:"signup"`
	Reset    string `json:"reset"`
}

type LoginInfo struct {
	MT string `json:"mt"`
}

type LoginInfoResult struct {
	Mt        string `json:"mt"`
	Error     int    `json:"error"`
	ErrorText string `json:"errorText"`
	User      struct {
		Digest     bool   `json:"digest"`
		Ntlm       bool   `json:"ntlm"`
		OAuth2     bool   `json:"oauth2"`
		OAuth2Name string `json:"oauth2Name"`
	} `json:"user"`
	Session struct {
		Digest bool `json:"digest"`
	} `json:"session"`
}

type Login struct {
	Mt        string `json:"mt"`
	Type      string `json:"type"` // user/session
	UserAgent string `json:"userAgent"`
	Method    string `json:"method,omitempty"`
	Username  string `json:"username,omitempty"`
	Nonce     string `json:"nonce,omitempty"`
	Response  string `json:"response,omitempty"`
}

type Authenticate struct {
	Mt        string `json:"mt"`
	Type      string `json:"type"`
	Method    string `json:"method"`
	Domain    string `json:"domain"`
	Challenge string `json:"challenge"`
}
type MyAppUserInfo struct {
	Domain string `json:"domain"`
	Sip    string `json:"sip"`
	Guid   string `json:"guid"`
	Dn     string `json:"dn"`
	Num    string `json:"num"`
	Email  string `json:"email"`
	Prefix struct {
		Intl    string `json:"intl"`
		Ntl     string `json:"ntl"`
		Subs    string `json:"subs"`
		Area    string `json:"area"`
		Country string `json:"country"`
	} `json:"prefix"`
}
type LoginResult struct {
	Mt   string `json:"mt"`
	Info struct {
		User     MyAppUserInfo          `json:"user"`
		Account  map[string]interface{} `json:"account"`
		Settings map[string]interface{} `json:"settings"`
		Session  struct {
			Usr string
			Pwd string
		} `json:"session"`
		AltHttp string `json:"altHttp"`
		Alt     string `json:"alt"`
	} `json:"info"`
	Digest    string `json:"digest"`
	Error     int    `json:"error"`
	ErrorText string `json:"errorText"`
}

const LoginResultInvalidParameters = 1    // {"mt":"LoginResult","error":1,"errorText":"Invalid parameters"}
const LoginResultAuthenticationFailed = 2 // {"mt":"LoginResult","error":2,"errorText":"Authentication failed"}
const LoginResultSessionExpired = 5       // {"mt":"LoginResult","error":5,"errorText":"Session expired"}

type SessionAdded struct {
	Mt   string `json:"mt"`
	Id   string `json:"id"`
	Info struct {
		UserAgent string `json:"userAgent"`
		Timestamp int64  `json:"timestamp"`
	} `json:"info"`
}

type UpdateOwnPresence struct {
	Mt       string `json:"mt"`
	Presence []struct {
		Contact  string `json:"contact"`
		Activity string `json:"activity"`
		Status   string `json:"status"`
	} `json:"presence"`
}

type Authorize struct {
	Mt   string `json:"mt"`
	Code int    `json:"code"`
}

type Redirect struct {
	Mt   string `json:"mt"`
	Info struct {
		Host    string `json:"host"`
		Mod     string `json:"mod"`
		Session struct {
			Usr string `json:"usr"`
			Pwd string `json:"pwd"`
		} `json:"session"`
	} `json:"info"`
	Digest string `json:"digest"`
}

type App struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
	Info  struct {
		Hidden bool `json:"hidden"`
	} `json:"info"`
}

type UpdateAppsInfo struct {
	Mt  string `json:"mt"`
	App App    `json:"app"`
}
