package devicesapp

type SetDeviceFiler struct {
	Mt          string `json:"mt"`
	RecvUpdates bool   `json:"recvUpdates"`
}

func NewSetDeviceFiler(recvUpdates bool) SetDeviceFiler {
	return SetDeviceFiler{
		Mt:          "SetDeviceFiler",
		RecvUpdates: recvUpdates,
	}
}

// --------------------------------------
type GetUserInfo struct {
	Mt string `json:"mt"`
}

func NewGetUserInfo() GetUserInfo {
	return GetUserInfo{
		Mt: "GetUserInfo",
	}
}

type GetUserInfoResult struct {
	Mt    string `json:"mt"`
	Key   string `json:"major"`
	Major string `json:""`
	Admin bool   `json:"admin"`
}

// --------------------------------------
type GetDomains struct {
	Mt          string `json:"mt"`
	RecvUpdates bool   `json:"recvUpdates"`
}

func NewGetDomains(recvUpdates bool) GetDomains {
	return GetDomains{
		Mt:          "GetDomains",
		RecvUpdates: recvUpdates,
	}
}

type Domain struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name"`
	DeployAdminPasswords bool   `json:"deployAdminPasswords"`
	IsInstanceDomain     bool   `json:"isInstanceDomain"`
	Categories           []struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Config bool   `json:"config"`
	} `json:"categories"`
}

type GetDomainsResult struct {
	Mt      string   `json:"mt"`
	Last    bool     `json:"last"`
	Domains []Domain `json:"domains"`
}

// --------------------------------------
type GetUnassignedDevicesCount struct {
	Mt string `json:"mt"`
}

func NewGetUnassignedDevicesCount(recvUpdates bool) GetUnassignedDevicesCount {
	return GetUnassignedDevicesCount{
		Mt: "GetUnassignedDevicesCount",
	}
}

type GetUnassignedDevicesCountResult struct {
	Mt    string `json:"mt"`
	Count int    `json:"count"`
}

// --------------------------------------
type GetDevices struct {
	Mt            string `json:"mt"`
	RecvUpdates   bool   `json:"recvUpdates"`
	DomainIds     string `json:"domainIds"`
	Categories    string `json:"categories"`
	Subcategories string `json:"subcategories"`
	Unassigned    bool   `json:"unassigned"`
}

func NewGetDevices(recvUpdates bool, DomainIds string, categories string, subcategories string, unassigned bool) GetDevices {
	return GetDevices{
		Mt:            "GetDevices",
		RecvUpdates:   recvUpdates,
		DomainIds:     DomainIds,
		Categories:    categories,
		Subcategories: subcategories,
		Unassigned:    unassigned,
	}
}

type Device struct {
	Id     int `json:"id"`
	EthIfs []struct {
		If    string `json:"if"`
		Ipv4  string `json:"ipv4"`
		Ipva6 string `json:"ipv6"`
	} `json:"ethIfs"`
	HwId      string `json:"hwId"`
	Name      string `json:"name"`
	DomainId  int    `json:"domainId"`
	Product   string `json:"product"`
	Version   string `json:"version"`
	Type      string `json:"type"`
	PbxActive bool   `json:"pbxActive"`
	Online    bool   `json:"online"`
}

type GetDevicesResult struct {
	Mt      string   `json:"mt"`
	Last    bool     `json:"last"`
	Devices []Device `json:"devices"`
}

type DeviceUpdate struct {
	Mt     string `json:"mt"`
	Device Device `json:"device"`
}
