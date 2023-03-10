package rcc

import (
	"github.com/ricoschulte/go-myapps/service"
	log "github.com/sirupsen/logrus"
)

type RCC struct {
}

func (api *RCC) GetApiName() string {
	return "RCC"
}

func (api *RCC) OnConnect(connection *service.AppServicePbxConnection) {
	log.WithField("api", api.GetApiName()).Warn("OnConnect not implemented ")
}

func (api *RCC) OnDisconnect(connection *service.AppServicePbxConnection) {
	log.WithField("api", api.GetApiName()).Warn("OnDisconnect not implemented ")
}

func (api *RCC) HandleMessage(connection *service.AppServicePbxConnection, msg *service.BaseMessage, message []byte) {
	log.WithField("api", api.GetApiName()).Warn("HandleMessage not implemented ")
}
