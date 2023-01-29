package handler

import (
	"github.com/ricoschulte/go-myapps/connection"
)

type HandleUpdateOwnPresence struct{}

func (m *HandleUpdateOwnPresence) GetMt() string {
	return "UpdateOwnPresence"
}

func (m *HandleUpdateOwnPresence) HandleMessage(myAppsConnection *connection.MyAppsConnection, message []byte) error {
	myAppsConnection.Config.Printf("Presence of %v %v has changed", myAppsConnection.User.Dn, myAppsConnection.Nonce)
	return nil
}
