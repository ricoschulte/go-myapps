package pbxadminapi

import (
	"strconv"

	"github.com/ricoschulte/go-myapps/service"
)

// Request the available App Licenses
type GetPbxLicenses struct {
	service.BaseMessage
}

func NewGetPbxLicenses(src string) *GetPbxLicenses {
	return &GetPbxLicenses{
		BaseMessage: service.BaseMessage{
			Api: "PbxAdminApi",
			Mt:  "GetPbxLicenses",
			Src: src,
		},
	}
}

type GetPbxLicensesResult struct {
	service.BaseMessage
	Licenses []AppLic `json:"lic"`
}

// Request the available App Licenses
type GetAppLics struct {
	service.BaseMessage
}

func NewGetAppLics(src string) *GetAppLics {
	return &GetAppLics{
		BaseMessage: service.BaseMessage{
			Api: "PbxAdminApi",
			Mt:  "GetAppLics",
			Src: src,
		},
	}
}

type GetAppLicensesResult struct {
	service.BaseMessage
	Licenses []AppLic `json:"lic"`
}

type AppLic struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Count   string `json:"count"`
	Used    string `json:"used"`
	Local   string `json:"local"`
	Slaves  string `json:"slaves"`
	Key     string `json:"key"`
}

// returns the value as float64 or errValue when the string from the JSON couldnt be parsed as float
func (l *AppLic) GetCount(errValue float64) float64 {
	if s, err := strconv.ParseFloat(l.Count, 32); err == nil {
		return s
	} else {
		return errValue
	}
}

// returns the value as float64 or errValue when the string from the JSON couldnt be parsed as float
func (l *AppLic) GetUsed(errValue float64) float64 {
	if s, err := strconv.ParseFloat(l.Used, 32); err == nil {
		return s
	} else {
		return errValue
	}
}

// returns the value as float64 or errValue when the string from the JSON couldnt be parsed as float
func (l *AppLic) GetLocal(errValue float64) float64 {
	if s, err := strconv.ParseFloat(l.Local, 32); err == nil {
		return s
	} else {
		return errValue
	}
}

// returns the value as float64 or errValue when the string from the JSON couldnt be parsed as float
func (l *AppLic) GetSlaves(errValue float64) float64 {
	if s, err := strconv.ParseFloat(l.Slaves, 32); err == nil {
		return s
	} else {
		return errValue
	}
}

type PbxAdminApiEvent struct {
	Type                 int
	Connection           *service.AppServicePbxConnection
	GetAppLicensesResult *GetAppLicensesResult
	GetPbxLicensesResult *GetPbxLicensesResult
}

const PbxAdminApiEventDisconnect = -20
const PbxAdminApiEventConnect = -10
const PbxAdminApiGetAppLicsResult = 10
const PbxAdminApiGetPbxLicensesResult = 20
