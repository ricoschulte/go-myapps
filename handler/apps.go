package handler

import (
	"github.com/ricoschulte/go-myapps/connection"

	"github.com/goccy/go-json"
)

type HandleUpdateAppsInfo struct {
}

func (m *HandleUpdateAppsInfo) GetMt() string {
	return "UpdateAppsInfo"
}

func (m *HandleUpdateAppsInfo) HandleMessage(myAppsConnection *connection.MyAppsConnection, message []byte) error {
	var msgin struct {
		Mt  string `json:"mt"`
		App struct {
			Name  string `json:"name"`
			Title string `json:"title"`
			Text  string `json:"text"`
			Url   string `json:"url"`
			Info  struct {
				Hidden bool `json:"hidden"`
			} `json:"info"`
		} `json:"app"`
	}
	err := json.Unmarshal(message, &msgin)
	if err != nil {
		return err
	}
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
