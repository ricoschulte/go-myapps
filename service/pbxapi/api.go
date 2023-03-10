package pbxapi

import (
	"encoding/json"
	"sync"

	"github.com/ricoschulte/go-myapps/service"
	log "github.com/sirupsen/logrus"
)

type PbxApi struct {
	mu        sync.Mutex
	receivers []chan PbxApiEvent
}

func NewPbxApi() *PbxApi {
	return &PbxApi{}
}

func (api *PbxApi) GetApiName() string {
	return "PbxApi"
}

func (api *PbxApi) OnConnect(connection *service.AppServicePbxConnection) {
	api.sendEvent(PbxApiEvent{Type: PbxApiEventConnect, Connection: connection})
}

func (api *PbxApi) OnDisconnect(connection *service.AppServicePbxConnection) {
	api.sendEvent(PbxApiEvent{Type: PbxApiEventDisconnect, Connection: connection})
}

func (api *PbxApi) HandleMessage(connection *service.AppServicePbxConnection, msg *service.BaseMessage, message []byte) {

	switch msg.Mt {
	case "SubscribePresenceResult":
		msg := SubscribePresenceResult{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Errorf("PbxApi: error unmarshalling message: %v", err)
		}
		if msg.Error != 0 {
			log.Errorf("SubscribePresenceResult: error %d: %s", msg.Error, msg.Errortext)
		}
		api.sendEvent(PbxApiEvent{Type: PbxApiEventSubscribePresenceResult, SubscribePresenceResult: &msg, Connection: connection})
	case "PresenceState":
		msg := PresenceState{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Errorf("PbxApi: error unmarshalling message: %v", err)
		}
		api.sendEvent(PbxApiEvent{Type: PbxApiEventPresenceState, PresenceState: &msg, Connection: connection})
	case "PresenceUpdate":
		msg := PresenceUpdate{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Errorf("PbxApi: error unmarshalling message: %v", err)
		}
		api.sendEvent(PbxApiEvent{Type: PbxApiEventPresenceUpdate, PresenceUpdate: &msg, Connection: connection})
	case "SetPresenceResult":
		msg := SetPresenceResult{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Errorf("PbxApi: error unmarshalling message: %v", err)
		}
		if msg.Error != 0 {
			log.Errorf("SetPresenceResult: error %d: %s", msg.Error, msg.Errortext)
		}
		api.sendEvent(PbxApiEvent{Type: PbxApiEventSetPresenceResult, SetPresenceResult: &msg, Connection: connection})
	case "GetNodeInfoResult":
		msg := GetNodeInfoResult{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Errorf("PbxApi: error unmarshalling message: %v", err)
		}
		api.sendEvent(PbxApiEvent{Type: PbxApiEventGetNodeInfoResult, GetNodeInfoResult: &msg, Connection: connection})
	case "AddAlienCallResult":
		msg := AddAlienCallResult{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Errorf("PbxApi: error unmarshalling message: %v", err)
		}
		api.sendEvent(PbxApiEvent{Type: PbxApiEventAddAlienCallResult, AddAlienCallResult: &msg, Connection: connection})
	case "DelAlienCallResult":
		msg := DelAlienCallResult{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Errorf("PbxApi: error unmarshalling message: %v", err)
		}
		api.sendEvent(PbxApiEvent{Type: PbxApiEventDelAlienCallResult, DelAlienCallResult: &msg, Connection: connection})
	default:
		log.Warn("unknown message received: %s", msg.Mt)
	}
}

func (api *PbxApi) AddReceiver() chan PbxApiEvent {
	api.mu.Lock()
	defer api.mu.Unlock()
	ch := make(chan PbxApiEvent)
	api.receivers = append(api.receivers, ch)
	return ch
}

func (api *PbxApi) RemoveReceiver(ch chan PbxApiEvent) {
	api.mu.Lock()
	defer api.mu.Unlock()

	for i, c := range api.receivers {
		if c == ch {
			// Remove the channel from the slice
			api.receivers = append(api.receivers[:i], api.receivers[i+1:]...)
			// Close the channel to signal the receiver that it should stop listening
			close(ch)
			return
		}
	}
}

func (api *PbxApi) sendEvent(event PbxApiEvent) {
	api.mu.Lock()
	defer api.mu.Unlock()

	// Send the event to all registered receivers
	for _, ch := range api.receivers {
		ch <- event
	}
}
