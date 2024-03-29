package pbxtableusers

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/ricoschulte/go-myapps/service"
)

type Column struct {
	Udate bool `json:"update"`
}

type ReplicateStart struct {
	service.BaseMessage
	Add     bool              `json:"add,omitempty"`
	Del     bool              `json:"del,omitempty"`
	Columns map[string]Column `json:"columns"`
	Pseudo  []string          `json:"pseudo"`
}

func NewReplicateStart(add, del bool, columns map[string]Column, pseudo []string, src string) *ReplicateStart {
	return &ReplicateStart{
		BaseMessage: service.BaseMessage{
			Api: "PbxTableUsers",
			Mt:  "ReplicateStart",
			Src: src,
		},
		Add:     add,
		Del:     del,
		Columns: columns,
		Pseudo:  pseudo,
	}
}

type ReplicateNext struct {
	service.BaseMessage
}

func NewReplicateNext(src string) *ReplicateNext {
	return &ReplicateNext{
		BaseMessage: service.BaseMessage{
			Api: "PbxTableUsers",
			Mt:  "ReplicateNext",
			Src: src,
		},
	}
}

type ReplicateStartResult struct {
	service.BaseMessage
	Guid service.Guid `json:"guid"`
	//Columns ReplicatedObject `json:"columns"`
}

type ReplicateNextResult struct {
	service.BaseMessage
	ReplicatedObject ReplicatedObject `json:"columns"`
}

type ReplicateUpdate struct {
	service.BaseMessage
	ReplicatedObject ReplicatedObject `json:"columns"`
}

type ReplicateAdd struct {
	service.BaseMessage
	ReplicatedObject ReplicatedObject `json:"columns"`
}

type ReplicateDel struct {
	service.BaseMessage
	ReplicatedObject ReplicatedObject `json:"columns"`
}

type ReplicatedObject struct {
	Guid      string `json:"guid"`       // ReplicationString	guid	Globally unique identifier
	H323      string `json:"h323"`       // ReplicationString	h323	Username
	Pwd       string `json:"pwd"`        // ReplicationString	pwd	Password
	Cn        string `json:"cn"`         // ReplicationString	cn	Common name
	Dn        string `json:"dn"`         // ReplicationString	dn	Display name
	AppsMy    string `json:"apps-my"`    // ReplicationString	apps-my	List of the apps displayed on the home screen
	Config    string `json:"config"`     // ReplicationString	config	Config template
	Node      string `json:"node"`       // ReplicationString	node	Node
	Loc       string `json:"loc"`        // ReplicationString	loc	Location
	Hide      bool   `json:"hide"`       // ReplicationBool	hide	Hide from LDAP
	E164      string `json:"e164"`       // ReplicationString	e164	Phone number
	Cfpr      bool   `json:"cfpr"`       // ReplicationTristate	cfpr	Call forward based on Presence
	Tcfpr     string `json:"t-cfpr"`     // ReplicationTristate	t-cfpr	Call forward based on Presence inherited from the config template
	Pseudo    string `json:"pseudo"`     // ReplicationString	pseudo	XML Pseudo information of the object
	H323email bool   `json:"h323-email"` // ReplicationBool	h323-email	If true, the email is the username
	Apps      string `json:"apps"`       // ReplicationString	apps	List of the apps that the user has rights to access
	Fax       bool   `json:"fax"`        // ReplicationBool	fax	If true, the user has a fax license
	Emails    []struct {
		Email string `json:"email"` // ReplicationString	email	Email
	} `json:"emails"` // emails Table with the emails of the users
	Allows []struct {
		Name     string `json:"name"`     // ReplicationString	name	Filter name
		Grp      bool   `json:"grp"`      // ReplicationString	grp	If true, the name is a group name
		Visible  bool   `json:"visible"`  // ReplicationBool	visible	Visible
		Online   bool   `json:"online"`   // ReplicationBool	online	Online
		Presence bool   `json:"presence"` // ReplicationBool	presence	Presence
		Otf      bool   `json:"otf"`      // ReplicationBool	otf	On the phone
		Note     bool   `json:"note"`     // ReplicationBool	note	Presence note
		Dialog   bool   `json:"dialog"`   // ReplicationBool	dialog	Calls
		Ids      bool   `json:"ids"`      // ReplicationBool	ids	Calls with id
	} `json:"allows"` // allows Table with the visibility filters defined for the user
	Tallows []struct {
		Name     string `json:"name"`     // ReplicationString	name	Filter name
		Grp      bool   `json:"grp"`      // ReplicationString	grp	If true, the name is a group name
		Visible  bool   `json:"visible"`  // ReplicationBool	visible	Visible
		Online   bool   `json:"online"`   // ReplicationBool	online	Online
		Presence bool   `json:"presence"` // ReplicationBool	presence	Presence
		Otf      bool   `json:"otf"`      // ReplicationBool	otf	On the phone
		Note     bool   `json:"note"`     // ReplicationBool	note	Presence note
		Dialog   bool   `json:"dialog"`   // ReplicationBool	dialog	Calls
		Ids      bool   `json:"ids"`      // ReplicationBool	ids	Calls with id
	} `json:"t-allows"` // t-allows Table with the visibility filters defined on the config templates
	Grps []struct {
		Name string `json:"name"` // ReplicationString	name	Group name
		Mode string `json:"mode"` // ReplicationString	mode	Mode
		Dyn  string `json:"dyn"`  // ReplicationString	dyn	Dynamic
	} `json:"grps"` // grps Table with the users groups
	Devices []struct {
		Hw       string `json:"hw"`        // ReplicationString	hw	Hardware ID
		Text     string `json:"text"`      // ReplicationString	text	Name
		App      string `json:"app"`       // ReplicationString	app	App
		Admin    bool   `json:"admin"`     // ReplicationBool	admin	PBX Pwd
		Nofilter bool   `json:"no-filter"` // ReplicationBool	no-filter	No IP Filter
		Tls      bool   `json:"tls"`       // ReplicationBool	tls	TLS only
		Nomob    bool   `json:"no-mob"`    // ReplicationBool	no-mob	No Mobility
		Trusted  bool   `json:"trusted"`   // ReplicationBool	trusted	Reverse Proxy
		Sreg     bool   `json:"sreg"`      // ReplicationBool	sreg	Single Reg.
		Mr       bool   `json:"mr"`        // ReplicationBool	mr	Media Relay
		Voip     string `json:"voip"`      // ReplicationString	voip	Config VOIP
		Gkid     string `json:"gk-id"`     // ReplicationString	gk-id	Gatekeeper ID
		Prim     string `json:"prim"`      // ReplicationString	prim	Primary gatekeeper
	} `json:"devices"` // devices Table with the users devices
	Cds []struct {
		Type    string `json:"type"`     // ReplicationString	type	Diversion type (cfu` cfb or cfnr)
		Bool    string `json:"bool"`     // ReplicationString	bool	Boolean object
		Boolnot bool   `json:"bool-not"` // ReplicationBool		bool-not	Not flag (boolean object)
		E164    string `json:"e164"`     // ReplicationString	e164	Phone number
		H323    string `json:"h323"`     // ReplicationString	h323	Username
		Src     string `json:"src"`      // ReplicationString	src	Filters data on XML format
	} `json:"cds"` // cds Table with the users call diversions
	Forks []struct {
		E164     string `json:"e164"`     // ReplicationString	e164	Phone number
		H323     string `json:"h323"`     // ReplicationString	h323	Username
		Bool     string `json:"bool"`     // ReplicationString	bool	Boolean object
		Boolnot  bool   `json:"bool-not"` // ReplicationBool	bool-not	Not flag (boolean object)
		Mobility string `json:"mobility"` // ReplicationString	mobility	Mobility object
		App      string `json:"app"`      // ReplicationString	app	App
		Delay    int    `json:"delay"`    // ReplicationUnsigned	delay	Delay
		Hw       string `json:"hw"`       // ReplicationString	hw	Device
		Off      bool   `json:"off"`      // ReplicationBool	off	Disable
		Cw       bool   `json:"cw"`       // ReplicationBool	cw	Call-Waiting
		Min      int    `json:"min"`      // ReplicationUnsigned	min	Min-Alert
		Max      int    `json:"max"`      // ReplicationUnsigned	max	Max-Alert
	} `json:"forks"` // forks Table with the users forks
	Wakeups []struct {
		H        int    `json:"h"`        // ReplicationUnsigned	h	Hour
		M        int    `json:"m"`        // ReplicationUnsigned	m	Minute
		S        int    `json:"s"`        // ReplicationUnsigned	s	Second
		Name     string `json:"name"`     // ReplicationString	name
		Num      string `json:"num"`      // ReplicationString	num
		Retry    int    `json:"retry"`    // ReplicationUnsigned	retry
		Mult     bool   `json:"mult"`     // ReplicationBool	mult
		To       int    `json:"to"`       // ReplicationUnsigned	to
		Fallback string `json:"fallback"` // ReplicationString	fallback
		Bool     string `json:"bool"`     // ReplicationString	bool	Boolean object
		Boolnot  bool   `json:"bool-not"` // ReplicationBool	bool-not	Not flag (boolean object)
	} `json:"wakeups"` // wakeups Table with the users wakeups
}

// returns the Pseudo Type of the Object based on the Xml based Pseudo Data
// for user objects that do not have a explicit pseudo type a by default a pseudo type "user" is returned
// when the pseudo type is empty it is assumed that the object is a "user object"
// if the xml based infos could not be parsed, a error and a empty string is returned
func (obj *ReplicatedObject) GetPseudoType() (string, error) {
	if obj.Pseudo == "" {
		return PseudoTypeUser, nil
	}
	var pseudo PseudoTypeXml
	err := xml.Unmarshal([]byte(obj.Pseudo), &pseudo)
	if err != nil {
		return "", fmt.Errorf("Failed to parse PseudoType XML: %s", err)
	}

	// make sure that every type we do not know results in a error
	for _, s := range PseudoTypes {
		if s == pseudo.Type {
			return s, nil
		}
	}
	return "", errors.New("could not get a valid pseudo type from the object. none of the known types are found")

}

// root element of XML pased Pseudo Type
type PseudoTypeXml struct {
	XMLName xml.Name         `xml:"pseudo"`
	Type    string           `xml:"type,attr"`
	App     PseudoTypeXmlApp `xml:"app"`
}

type PseudoTypeXmlApp struct {
	URL           string `xml:"url,attr"`
	WebSocket     bool   `xml:"websocket,attr"`
	PBX           bool   `xml:"pbx,attr"`
	PBXSignal     bool   `xml:"pbxsignal,attr"`
	EPSignal      bool   `xml:"epsignal,attr"`
	Messages      bool   `xml:"messages,attr"`
	TableUsers    bool   `xml:"tableusers,attr"`
	Admin         bool   `xml:"admin,attr"`
	Services      bool   `xml:"services,attr"`
	RCC           bool   `xml:"rcc,attr"`
	Impersonation bool   `xml:"impersonation,attr"`
}

var (
	PseudoTypeAp             = "ap"
	PseudoTypeApp            = "app"
	PseudoTypeBcConference   = "bc_conf"
	PseudoTypeBroadcast      = "broadcast"
	PseudoTypeConference     = "conference"
	PseudoTypeConfig         = "config"
	PseudoTypeDect           = "dect"
	PseudoTypeDirsearch      = "dirsearch"
	PseudoTypeDtmfCtrl       = "dtmf-ctrl"
	PseudoTypeUc             = "uc"
	PseudoTypeFax            = "fax"
	PseudoTypeIcp            = "icp"
	PseudoTypeLdap           = "ldap"
	PseudoTypeMCastAnnounce  = "multicast"
	PseudoTypeMessageWaiting = "mwi"
	PseudoTypeMessages       = "messages"
	PseudoTypeMobility       = "mobility"
	PseudoTypeNode           = "node"
	PseudoTypeNumberMap      = "map"
	PseudoTypePush           = "push"
	PseudoTypeQuickdial      = "qdial"
	PseudoTypeSessionBorder  = "sbc"
	PseudoTypeSettings       = "settings"
	PseudoTypeVoicemail      = "vm"
	PseudoTypeGw             = "gw"
	PseudoTypePbx            = "loc"
	PseudoTypeWaitingQueue   = "waiting"
	PseudoTypeExecutive      = "executive"
	PseudoTypeBoolean        = "bool"
	PseudoTypeTrunk          = "trunk"
	PseudoTypeUser           = "user" // by default, the user object have not pseudo type aka a empty value ""
)

// List of known Pseudo Types of pbx objects
// Use this list with care. If i use this list with >20 items with 'ReplicateStart' the pbx bootloops with 'Unexpected Restart WATCHDOG (2) of 137818'
var PseudoTypes = []string{
	PseudoTypeAp,
	PseudoTypeApp,
	PseudoTypeBcConference,
	PseudoTypeBroadcast,
	PseudoTypeConference,
	PseudoTypeConfig,
	PseudoTypeDect,
	PseudoTypeDirsearch,
	PseudoTypeDtmfCtrl,
	PseudoTypeUc,
	PseudoTypeFax,
	PseudoTypeIcp,
	PseudoTypeLdap,

	PseudoTypeMCastAnnounce,
	PseudoTypeMessageWaiting,
	PseudoTypeMessages,
	PseudoTypeMobility,
	PseudoTypeNode,
	PseudoTypeNumberMap,
	PseudoTypePush,
	PseudoTypeQuickdial,
	PseudoTypeSessionBorder,
	PseudoTypeSettings,
	PseudoTypeVoicemail,
	PseudoTypeGw,
	PseudoTypePbx,
	PseudoTypeWaitingQueue,
	PseudoTypeExecutive,
	PseudoTypeBoolean,
	PseudoTypeTrunk,
	"", // users / PseudoTypeUser
}

var AllFields = []string{}
var AllColumns = map[string]Column{} // all known Columns to subscribe to

func init() {

	AllFields = append(AllFields, TableUsers...)
	AllFields = append(AllFields, "emails")
	AllFields = append(AllFields, "allows")
	AllFields = append(AllFields, "t-allows")
	AllFields = append(AllFields, "grps")
	AllFields = append(AllFields, "devices")
	AllFields = append(AllFields, "cds")
	AllFields = append(AllFields, "forks")
	AllFields = append(AllFields, "wakeups")

	for _, col := range AllFields {
		AllColumns[col] = Column{Udate: true}
	}
}

/*
Tables
*/

// Table with the user data
var TableUsers = []string{
	"guid",       // ReplicationString	guid	Globally unique identifier
	"h323",       // ReplicationString	h323	Username
	"pwd",        // ReplicationString	pwd	Password
	"cn",         // ReplicationString	cn	Common name
	"dn",         // ReplicationString	dn	Display name
	"apps-my",    // ReplicationString	apps-my	List of the apps displayed on the home screen
	"config",     // ReplicationString	config	Config template
	"node",       // ReplicationString	node	Node
	"loc",        // ReplicationString	loc	Location
	"hide",       // ReplicationBool	hide	Hide from LDAP
	"e164",       // ReplicationString	e164	Phone number
	"cfpr",       // ReplicationTristate	cfpr	Call forward based on Presence
	"t-cfpr",     // ReplicationTristate	t-cfpr	Call forward based on Presence inherited from the config template
	"pseudo",     // ReplicationString	pseudo	Pseudo information of the object
	"h323-email", // ReplicationBool	h323-email	If true, the email is the username
	"apps",       // ReplicationString	apps	List of the apps that the user has rights to access
	"fax",        // ReplicationBool	fax	If true, the user has a fax license
}

type PbxTableUsersEvent struct {
	Type       int
	Connection *service.AppServicePbxConnection
	Object     *ReplicatedObject
}

const PbxTableUsersEventDisconnect = -20
const PbxTableUsersEventConnect = -10
const PbxTableUsersEventInitial = 0
const PbxTableUsersEventInitialDone = 5
const PbxTableUsersEventAdd = 10
const PbxTableUsersEventUpdate = 20
const PbxTableUsersEventDelete = 30
