package main

import (
	"sync"

	"github.com/ricoschulte/go-myapps/connection"
	"github.com/ricoschulte/go-myapps/handler"
)

func main() {
	var wg sync.WaitGroup

	var myAppConfigs []*connection.Config

	// add a account of a user of the same pbx
	myAppConfigs = append(myAppConfigs, &connection.Config{
		Host:               "192.168.178.200:443",
		Username:           "exampleUser",
		Password:           "examplePassword",
		InsecureSkipVerify: false,
		UserAgent:          "myApps Go client",
		SessionFilePath:    "myapps_session.json",
		Debug:              true,
	})

	// add a account of a second user of the same pbx
	myAppConfigs = append(myAppConfigs, &connection.Config{
		Host:            "192.168.178.200:443",
		Username:        "exampleUser2",
		Password:        "examplePassword2",
		UserAgent:       "myBot (Go)",
		SessionFilePath: "myapps_session_2.json",
		Debug:           true,
	})

	// add a account of a user of another pbx
	myAppConfigs = append(myAppConfigs, &connection.Config{
		Host:            "pbx.company.com",
		Username:        "exampleUser3",
		Password:        "examplePassword3",
		UserAgent:       "myApps (Go)",
		SessionFilePath: "myapps_session_3.json",
		Debug:           true,
	})

	for _, config := range myAppConfigs {
		// configure Handlers that should be used to handle messages for this connection
		config.Handler.AddHandler(&handler.HandleUpdateAppsInfo{})
		config.Handler.AddHandler(&handler.HandleUpdateAppsComplete{})
		config.Handler.AddHandler(&handler.HandleUpdateOwnPresence{})

		// start the session with this configuration
		go config.StartSession(&wg)
	}

	wg.Wait() // wait for all sessions to complete, to keep the sessions open and our program to run

	select {}
}
