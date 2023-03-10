package services

import (
	"github.com/ricoschulte/go-myapps/service"
	log "github.com/sirupsen/logrus"
)

type Services struct {
}

func (api *Services) GetApiName() string {
	return "Services"
}

func (api *Services) OnConnect(connection *service.AppServicePbxConnection) {
	log.WithField("api", api.GetApiName()).Warn("OnConnect not implemented ")
}

func (api *Services) OnDisconnect(connection *service.AppServicePbxConnection) {
	log.WithField("api", api.GetApiName()).Warn("OnDisconnect not implemented ")
}

func (api *Services) HandleMessage(connection *service.AppServicePbxConnection, msg *service.BaseMessage, message []byte) {
	log.WithField("api", api.GetApiName()).Warn("HandleMessage not implemented ")
}
