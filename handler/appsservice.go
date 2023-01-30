package handler

import (
	"encoding/json"

	"github.com/ricoschulte/go-myapps/appservice"
	"github.com/ricoschulte/go-myapps/connection"
)

/*
* start a connection to a appservice when the app is avaiable in the account
 */
type HandleAppService struct {
	Name                   string                                      // the name of the app
	MessageHandlerRegister appservice.AppServiceMessageHandlerRegister // list of message handler on the session
}

func (m *HandleAppService) GetMt() string {
	return "UpdateAppsInfo"
}

func (m *HandleAppService) HandleMessage(myAppsConnection *connection.MyAppsConnection, message []byte) error {
	var msgin = connection.UpdateAppsInfo{}
	err := json.Unmarshal(message, &msgin)
	if err != nil {
		return err
	}

	switch msgin.App.Name {
	case m.Name:
		go m.StartConnectionToAppService(myAppsConnection, msgin)
	default:
		//myAppsConnection.Config.Printf("uninteresting App %v: '%v'", msgin.App.Name, msgin.App.Url)
	}
	return nil
}

func (m *HandleAppService) StartConnectionToAppService(myAppsConnection *connection.MyAppsConnection, message connection.UpdateAppsInfo) error {
	myAppsConnection.Config.Printf("starting connection to: %v", message.App.Name)
	appconnect := appservice.AppServiceClient{
		MyAppsConnection:       myAppsConnection,
		AppInfo:                &message.App,
		MessageHandlerRegister: &m.MessageHandlerRegister,
	}

	err := appconnect.Connect()
	if err != nil {
		myAppsConnection.Config.Printf("error while starting connection to '%s': %v", message.App.Name, err)
	}
	return err
}
