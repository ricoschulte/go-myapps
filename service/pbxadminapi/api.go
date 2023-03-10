package pbxadminapi

import (
	"github.com/ricoschulte/go-myapps/service"
	log "github.com/sirupsen/logrus"
)

type PbxAdminApi struct {
}

func (api *PbxAdminApi) GetApiName() string {
	return "PbxAdminApi"
}

func (api *PbxAdminApi) OnConnect(connection *service.AppServicePbxConnection) {
	log.WithField("api", api.GetApiName()).Warn("OnConnect not implemented ")
}

func (api *PbxAdminApi) OnDisconnect(connection *service.AppServicePbxConnection) {
	log.WithField("api", api.GetApiName()).Warn("OnDisconnect not implemented ")
}

func (api *PbxAdminApi) HandleMessage(connection *service.AppServicePbxConnection, msg *service.BaseMessage, message []byte) {
	log.WithField("api", api.GetApiName()).Warn("HandleMessage not implemented ")
}
