package sysclient

import "encoding/json"

type Platform struct {
	Type string `json:"type"`
	Fxs  bool   `json:"fxs"`
}

type EthIf struct {
	If   string `json:"if"`
	Ipv4 string `json:"ipv4"`
	Ipv6 string `json:"ipv6"`
}

type Identity struct {
	Id        string   `json:"id"`
	Product   string   `json:"product"`
	Version   string   `json:"version"`
	FwBuild   string   `json:"fwBuild"`
	BcBuild   string   `json:"bcBuild"`
	Major     string   `json:"major"`
	Fw        string   `json:"fw"`
	Bc        string   `json:"bc"`
	Mini      bool     `json:"mini"`
	PbxActive bool     `json:"pbxActive"`
	Other     bool     `json:"other"`
	Platform  Platform `json:"platform"`
	Digest    string   `json:"digest,omitempty"`
	EthIfs    []EthIf  `json:"ethIfs,omitempty"` // up to 3 eth interfaces are allowed. 
}

func NewIdentity(json_bytes []byte) (*Identity, error) {
	var identity Identity
	err := json.Unmarshal(json_bytes, &identity)
	if err != nil {
		return nil, err
	}
	return &identity, nil
}

func (i Identity) ToBytes() ([]byte, error) {
	msgbytes, err_to_json := json.Marshal(i)
	if err_to_json != nil {
		return nil, err_to_json
	}
	return msgbytes, nil
}
