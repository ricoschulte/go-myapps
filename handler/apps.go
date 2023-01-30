package handler

import (
	"encoding/json"

	"github.com/ricoschulte/go-myapps/connection"
)

/*
* adds the avaiable apps of a account to the connection
 */
type HandleUpdateAppsInfo struct {
}

func (m *HandleUpdateAppsInfo) GetMt() string {
	return "UpdateAppsInfo"
}

func (m *HandleUpdateAppsInfo) HandleMessage(myAppsConnection *connection.MyAppsConnection, message []byte) error {
	var msgin = connection.UpdateAppsInfo{}
	err := json.Unmarshal(message, &msgin)
	if err != nil {
		return err
	}

	myAppsConnection.Apps[msgin.App.Name] = &msgin.App
	return nil
}

type HandleUpdateAppsComplete struct{}

func (m *HandleUpdateAppsComplete) GetMt() string {
	return "UpdateAppsComplete"
}

func (m *HandleUpdateAppsComplete) HandleMessage(myAppsConnection *connection.MyAppsConnection, message []byte) error {
	var msgin struct {
		Mt         string `json:"mt"`
		DeviceApps []struct {
			Name      string `json:"name"`
			Title     string `json:"title"`
			Deviceapp string `json:"deviceapp"`
		} `json:"deviceApps"`
		Selected string `json:"selected"`
	}
	err := json.Unmarshal(message, &msgin)
	if err != nil {
		return err
	}
	myAppsConnection.Config.Printf("Number of DeviceApps: '%v'", len(msgin.DeviceApps))
	return nil
}
