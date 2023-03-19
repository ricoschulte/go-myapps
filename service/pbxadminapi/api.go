package pbxadminapi

import (
	"encoding/json"
	"sync"

	"github.com/ricoschulte/go-myapps/service"
	log "github.com/sirupsen/logrus"
)

type PbxAdminApi struct {
	mu        sync.Mutex
	receivers []chan PbxAdminApiEvent
}

func NewPbxAdminApi() *PbxAdminApi {
	return &PbxAdminApi{}
}

func (api *PbxAdminApi) GetApiName() string {
	return "PbxAdminApi"
}

func (api *PbxAdminApi) OnConnect(connection *service.AppServicePbxConnection) {
	api.sendEvent(PbxAdminApiEvent{Type: PbxAdminApiEventConnect, Connection: connection})
}

func (api *PbxAdminApi) OnDisconnect(connection *service.AppServicePbxConnection) {
	api.sendEvent(PbxAdminApiEvent{Type: PbxAdminApiEventDisconnect, Connection: connection})
}

func (api *PbxAdminApi) HandleMessage(connection *service.AppServicePbxConnection, msg *service.BaseMessage, message []byte) {
	switch msg.Mt {
	case "GetAppLicsResult":
		msg := GetAppLicsResult{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Errorf("PbxApi: error unmarshalling message: %v", err)
		}
		api.sendEvent(PbxAdminApiEvent{Type: PbxAdminApiGetAppLicsResult, GetAppLicsResult: &msg, Connection: connection})
	default:
		log.Warn("unknown message received: %s", msg.Mt)
	}
}

func (api *PbxAdminApi) AddReceiver() chan PbxAdminApiEvent {
	api.mu.Lock()
	defer api.mu.Unlock()
	ch := make(chan PbxAdminApiEvent)
	api.receivers = append(api.receivers, ch)
	return ch
}

func (api *PbxAdminApi) RemoveReceiver(ch chan PbxAdminApiEvent) {
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

func (api *PbxAdminApi) sendEvent(event PbxAdminApiEvent) {
	api.mu.Lock()
	defer api.mu.Unlock()

	// Send the event to all registered receivers
	for _, ch := range api.receivers {
		ch <- event
	}
}
