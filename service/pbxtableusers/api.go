package pbxtableusers

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/ricoschulte/go-myapps/service"
	log "github.com/sirupsen/logrus"
)

type PbxTableUsers struct {
	ReplicatedObjects      map[string]ReplicatedObject // synced objects
	ReplicatedObjectsMutex sync.Mutex
	mu                     sync.Mutex
	receivers              []chan PbxTableUsersEvent
}

func NewPbxTableUsers() *PbxTableUsers {
	return &PbxTableUsers{
		ReplicatedObjects: map[string]ReplicatedObject{},
	}
}

func (api *PbxTableUsers) GetApiName() string {
	return "PbxTableUsers"
}

func (api *PbxTableUsers) OnConnect(connection *service.AppServicePbxConnection) {
	api.sendEvent(PbxTableUsersEvent{Type: PbxTableUsersEventConnect, Connection: connection})
}

func (api *PbxTableUsers) OnDisconnect(connection *service.AppServicePbxConnection) {
	api.sendEvent(PbxTableUsersEvent{Type: PbxTableUsersEventDisconnect, Connection: connection})
}

func (api *PbxTableUsers) HandleMessage(connection *service.AppServicePbxConnection, msg *service.BaseMessage, message []byte) {
	api.ReplicatedObjectsMutex.Lock()
	defer api.ReplicatedObjectsMutex.Unlock()
	switch msg.Mt {
	case "ReplicateStartResult":
		msg := ReplicateStartResult{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Errorf("PbxApi: error unmarshalling message: %v", err)
		}
		mbytes, _ := json.Marshal(NewReplicateNext("src_" + strconv.FormatInt(time.Now().UnixNano(), 10)))
		connection.WriteMessage(mbytes)
	case "ReplicateNextResult":
		msg := ReplicateNextResult{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Errorf("PbxApi: error unmarshalling message: %v", err)
		}
		if len(msg.ReplicatedObject.Guid) > 0 {
			api.ReplicatedObjects[msg.ReplicatedObject.Guid] = msg.ReplicatedObject
			api.sendEvent(PbxTableUsersEvent{Type: PbxTableUsersEventInitial, Object: &msg.ReplicatedObject, Connection: connection})

			mbytes, _ := json.Marshal(NewReplicateNext("src_" + strconv.FormatInt(time.Now().UnixNano(), 10)))
			connection.WriteMessage(mbytes)
		} else {
			api.sendEvent(PbxTableUsersEvent{Type: PbxTableUsersEventInitialDone, Connection: connection})
		}
	case "ReplicateAdd":
		msg := ReplicateAdd{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Errorf("PbxApi: error unmarshalling message: %v", err)
		}
		api.ReplicatedObjects[msg.ReplicatedObject.Guid] = msg.ReplicatedObject
		api.sendEvent(PbxTableUsersEvent{Type: PbxTableUsersEventAdd, Object: &msg.ReplicatedObject, Connection: connection})
	case "ReplicateUpdate":
		msg := ReplicateUpdate{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Errorf("PbxApi: error unmarshalling message: %v", err)
		}
		api.ReplicatedObjects[msg.ReplicatedObject.Guid] = msg.ReplicatedObject
		api.sendEvent(PbxTableUsersEvent{Type: PbxTableUsersEventUpdate, Object: &msg.ReplicatedObject, Connection: connection})
	case "ReplicateDel":
		msg := ReplicateDel{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Errorf("PbxApi: error unmarshalling message: %v", err)
		}
		delete(api.ReplicatedObjects, msg.ReplicatedObject.Guid)
		api.sendEvent(PbxTableUsersEvent{Type: PbxTableUsersEventDelete, Object: &msg.ReplicatedObject, Connection: connection})

	default:
		log.Warn("unknown message received: %s", string(message))
	}
}

func (api *PbxTableUsers) AddReceiver() chan PbxTableUsersEvent {
	api.mu.Lock()
	defer api.mu.Unlock()
	ch := make(chan PbxTableUsersEvent)
	api.receivers = append(api.receivers, ch)
	return ch
}

func (api *PbxTableUsers) RemoveReceiver(ch chan PbxTableUsersEvent) {
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

func (api *PbxTableUsers) sendEvent(event PbxTableUsersEvent) {
	api.mu.Lock()
	defer api.mu.Unlock()

	// Send the event to all registered receivers
	for _, ch := range api.receivers {
		ch <- event
	}
}

// returns the current replicated objects in a go routine save way
func (api *PbxTableUsers) GetReplicatedObjects() map[string]ReplicatedObject {
	// lock the map to prevent concurrent writes/reads
	api.ReplicatedObjectsMutex.Lock()
	defer api.ReplicatedObjectsMutex.Unlock()
	listcopy := make(map[string]ReplicatedObject)
	for key, value := range api.ReplicatedObjects {
		listcopy[key] = value
	}
	return listcopy

}
