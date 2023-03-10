package pbxapi

import (
	"fmt"

	"github.com/ricoschulte/go-myapps/service"
)

const (
	ActivityAvaiable   = ""
	ActivityAway       = "away"
	ActivityBusy       = "busy"
	ActivityDnd        = "dnd"
	ActivityOnThePhone = "on-the-phone" // can not be set
)

type SubscribePresence struct {
	service.BaseMessage
	Num string `json:"num,omitempty"`
	Sip string `json:"sip,omitempty"`
}

func NewSubscribePresenceWithNum(num, src string) *SubscribePresence {
	return &SubscribePresence{
		BaseMessage: service.BaseMessage{
			Api: "PbxApi",
			Mt:  "SubscribePresence",
			Src: src,
		},
		Num: num,
	}
}

func NewSubscribePresenceWithSip(sip, src string) *SubscribePresence {
	return &SubscribePresence{
		BaseMessage: service.BaseMessage{
			Api: "PbxApi",
			Mt:  "SubscribePresence",
			Src: src,
		},
		Sip: sip,
	}

}

type SubscribePresenceResult struct {
	service.BaseMessage
	Error     int    `json:"error,omitempty"`
	Errortext string `json:"errorText,omitempty"`
}

type PresenceState struct {
	service.BaseMessage
	Sip   string `json:"sip"`
	Dn    string `json:"dn"`
	Num   string `json:"num"`
	Email string `json:"email"`
}

type Presence struct {
	Contact string `json:"contact"` // tel:   im:
	Status  string `json:"status"`  // closed open
	Note    string `json:"note"`
}

type PresenceUpdate struct {
	service.BaseMessage
	Presence []Presence `json:"presence"`
}

func (update *PresenceUpdate) GetPresenceByContact(contact string) (*Presence, error) {
	for _, presence := range update.Presence {
		if presence.Contact == contact {
			return &presence, nil
		}
	}
	return nil, fmt.Errorf("no contact '%s'", contact)
}

// Sets the presence for a given contact of a user, defined by the SIP URI or GUID.
type SetPresence struct {
	service.BaseMessage
	Guid     string `json:"guid,omitempty"`
	Sip      string `json:"sip,omitempty"`
	Contact  string `json:"contact,omitempty"`  // tel:   im:
	Activity string `json:"activity,omitempty"` // busy, dnd ...
	Note     string `json:"note,omitempty"`
}

// Message sent back to confirm that setting the presence has been completed or failed.
type SetPresenceResult struct {
	service.BaseMessage
	Error     int    `json:"error,omitempty"`
	Errortext string `json:"errorText,omitempty"`
}

// Sets the presence for a given contact of a user, defined by the GUID.
func NewSetPresenceWithGuid(guid, contact, activity, note, src string) *SetPresence {
	return &SetPresence{
		BaseMessage: service.BaseMessage{
			Api: "PbxApi",
			Mt:  "SetPresence",
			Src: src,
		},
		Guid:     guid,
		Contact:  contact,
		Activity: activity,
		Note:     note,
	}
}

// Sets the presence for a given contact of a user, defined by the SIP URI.
func NewSetPresenceWithSip(sip, contact, activity, note, src string) *SetPresence {
	return &SetPresence{
		BaseMessage: service.BaseMessage{
			Api: "PbxApi",
			Mt:  "SetPresence",
			Src: src,
		},
		Sip:      sip,
		Contact:  contact,
		Activity: activity,
		Note:     note,
	}
}

// Requests information about the node of the user that is authenticated on the underlying AppWebsocket connection.
type GetNodeInfo struct {
	service.BaseMessage
}

func NewGetNodeInfo(src string) *GetNodeInfo {
	return &GetNodeInfo{
		BaseMessage: service.BaseMessage{
			Api: "PbxApi",
			Mt:  "GetNodeInfo",
			Src: src,
		},
	}
}

// Contains information about the node of the user that is authenticated on the underlying AppWebsocket connection.
type GetNodeInfoResult struct {
	service.BaseMessage
	Name                string `json:"name"`         // name Name of the node.
	PrefixInternational string `json:"prefix_intl"`  // prefix_intl Prefix for dialing international numbers.
	PrefixNational      string `json:"prefix_ntl"`   // prefix_ntl Prefix for dialing national numbers.
	PrefixSubscriber    string `json:"prefix_subs"`  // prefix_subs Subscriber number prefix.
	CountryCode         string `json:"country_code"` // country_code Country code.
}

// Adds a call to the PBX, which will result in a busy state of the user and also shows up as on the phone presence
type AddAlienCall struct {
	service.BaseMessage
}

func NewAddAlienCall(src string) *AddAlienCall {
	return &AddAlienCall{
		BaseMessage: service.BaseMessage{
			Api: "PbxApi",
			Mt:  "AddAlienCall",
			Src: src,
		},
	}
}

type AddAlienCallResult struct {
	service.BaseMessage
	Id int `json:"id"` // The id can be used to delete the call
}

// Deletes a call added with AddAlienCall. The id identifies the call to be deleted
type DelAlienCall struct {
	service.BaseMessage
	Id int `json:"id"` // The id can be used to delete the call
}

func NewDelAlienCall(id int, src string) *DelAlienCall {
	return &DelAlienCall{
		BaseMessage: service.BaseMessage{
			Api: "PbxApi",
			Mt:  "DelAlienCall",
			Src: src,
		},
		Id: id,
	}
}

type DelAlienCallResult struct {
	service.BaseMessage
}

type PbxApiEvent struct {
	Type                    int
	Connection              *service.AppServicePbxConnection
	PresenceState           *PresenceState
	PresenceUpdate          *PresenceUpdate
	SubscribePresenceResult *SubscribePresenceResult
	SetPresenceResult       *SetPresenceResult
	GetNodeInfoResult       *GetNodeInfoResult
	AddAlienCallResult      *AddAlienCallResult
	DelAlienCallResult      *DelAlienCallResult
}

const PbxApiEventDisconnect = -20
const PbxApiEventConnect = -10
const PbxApiEventSubscribePresenceResult = 10
const PbxApiEventPresenceState = 20
const PbxApiEventPresenceUpdate = 30
const PbxApiEventSetPresenceResult = 40
const PbxApiEventGetNodeInfoResult = 50
const PbxApiEventAddAlienCallResult = 60
const PbxApiEventDelAlienCallResult = 65
