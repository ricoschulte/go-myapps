package main

import (
	"sync"

	"github.com/ricoschulte/go-myapps/connection"
	"github.com/ricoschulte/go-myapps/handler"
)

func main() {
	var wg sync.WaitGroup

	accountConfig := connection.Config{
		Host:               "192.168.178.200:443",
		Username:           "exampleUser",
		Password:           "examplePassword",
		InsecureSkipVerify: true,
		UserAgent:          "myApps Go client",
		SessionFilePath:    "myapps_session.json",
		Debug:              true,
	}

	accountConfig.Handler.AddHandler(&handler.HandleUpdateAppsInfo{})
	accountConfig.Handler.AddHandler(&handler.HandleUpdateAppsComplete{})
	accountConfig.Handler.AddHandler(&handler.HandleUpdateOwnPresence{})

	go accountConfig.StartSession(&wg)

	wg.Wait()
	select {}
}
