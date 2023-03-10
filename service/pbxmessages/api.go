package pbxmessages

import (
	"github.com/ricoschulte/go-myapps/service"
	log "github.com/sirupsen/logrus"
)

type PbxMessages struct {
}

func (api *PbxMessages) GetApiName() string {
	return "PbxMessages"
}

func (api *PbxMessages) OnConnect(connection *service.AppServicePbxConnection) {
	log.WithField("api", api.GetApiName()).Warn("OnConnect not implemented ")
}

func (api *PbxMessages) OnDisconnect(connection *service.AppServicePbxConnection) {
	log.WithField("api", api.GetApiName()).Warn("OnDisconnect not implemented ")
}

func (api *PbxMessages) HandleMessage(connection *service.AppServicePbxConnection, msg *service.BaseMessage, message []byte) {
	log.WithField("api", api.GetApiName()).Warn("HandleMessage not implemented ")
}
