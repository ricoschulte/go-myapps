package pbximpersonation

import (
	"github.com/ricoschulte/go-myapps/service"
	log "github.com/sirupsen/logrus"
)

type PbxImpersonation struct {
}

func (api *PbxImpersonation) GetApiName() string {
	return "PbxImpersonation"
}

func (api *PbxImpersonation) OnConnect(connection *service.AppServicePbxConnection) {
	log.WithField("api", api.GetApiName()).Warn("OnConnect not implemented ")
}

func (api *PbxImpersonation) OnDisconnect(connection *service.AppServicePbxConnection) {
	log.WithField("api", api.GetApiName()).Warn("OnDisconnect not implemented ")
}

func (api *PbxImpersonation) HandleMessage(connection *service.AppServicePbxConnection, msg *service.BaseMessage, message []byte) {
	log.WithField("api", api.GetApiName()).Warn("HandleMessage not implemented ")
}
