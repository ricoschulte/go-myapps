# go-myapps

A golang client to connect to one or more myApps accounts of users at a innovaphone pbx via websocket.

## Introduction

The go-myapps library is a client library for the myApps API, written in Go. This library provides a simple and convenient way to interact with the myApps API, allowing you to easily authenticate, retrieve data, and handle events from myApps. The library is located at github.com/ricoschulte/go-myapps. 

## Getting started

To get started with the go-myapps library, you'll need to import the connection and handler packages.

``` GO
import (
	"github.com/ricoschulte/go-myapps/connection"
	"github.com/ricoschulte/go-myapps/handler"
)
```
This will allow you to use the "connection" and "handler" packages from "go-myapps" in your code.

In order to connect to a myApps server, you need to create a connection.Config struct. This struct holds all of the information needed to connect to the myApps server and to configure your session.

Here is an example of how to create a connection.Config struct:


``` GO
accountConfig := &connection.Config{
	Host:            "pbx.company.com",
	Username:        "exampleUser",
	Password:        "examplePassword",
	UserAgent:       "myApps (Go)",
	SessionFilePath: "myapps_session.json",
	Debug:           true,
    InsecureSkipVerify: false,
}
```

In this example, the following information is provided:

- **Host**: The hostname or IP address of the myApps server, including the port number.
- **Username**: The username of the myApps account you want to use.
- **Password**: The password for the myApps account you want to use.
- **UserAgent**: The user agent that will be sent to the myApps server. This is used to identify sessions of the client at the Account Security list within the myApps Clients.
- **SessionFilePath**: The file path where the session keys state will be stored. This allows you to resume a session after a disconnect. Please note that they are (for now) unencrypted stored.
- **Debug**: A boolean value indicating whether or not to enable debug logging. Default is false, meaning no debug messages.
- **InsecureSkipVerify**: A boolean value indicating whether or not to verify the SSL/TLS certificate. Default is false, so connections are aborted, if the Host does not provide a valid certificate.

You can use as many Accounts as you like, even accounts on different hosts/pbx.

## Configuring Handlers

Handlers are functions that are called when a certain type of message is received from the myApps server. You can use these to handle messages in your own way.

To configure handlers for your myApps connection, you need to create instances of handler structs and add them to the Handler field of your connection.Config struct.

Here is an example of how to add three handlers to your connection.Config struct:

``` GO
// configure Handlers that should be used to handle messages for this connection
accountConfig.Handler.AddHandler(&handler.HandleUpdateAppsInfo{})
accountConfig.Handler.AddHandler(&handler.HandleUpdateAppsComplete{})
accountConfig.Handler.AddHandler(&handler.HandleUpdateOwnPresence{})
```

In this example, the HandleUpdateAppsInfo, HandleUpdateAppsComplete, and HandleUpdateOwnPresence handlers are added to the config.Handler field.

## Starting the session

To start the myApps session, you need to call the StartSession method on your connection.Config struct. This will start the session and connect to the myApps server.

Here is an example of how to start the session for a single connection.Config struct:

``` GO
var wg sync.WaitGroup
accountConfig.StartSession(&wg)
wg.Wait() 
```

In this example, the StartSession method is called with a sync.WaitGroup value as its argument.
The wait group is used to wait for all sessions in that WaitGroup to complete, to keep the sessions open.

## Examples

A complete example of a program using this client can be found in `example.go`.

### Single user

``` GO
package main

import (
	"sync"

	"github.com/ricoschulte/go-myapps/connection"
	"github.com/ricoschulte/go-myapps/handler"
)

func main() {
	var wg sync.WaitGroup

	config := connection.Config{
		Host:               "192.168.178.200:443",
		Username:           "exampleUser",
		Password:           "examplePassword",
		InsecureSkipVerify: true,
		UserAgent:          "myApps Go client",
		SessionFilePath:    "myapps_session.json",
		Debug:              true,
	}

	config.Handler.AddHandler(&handler.HandleUpdateAppsInfo{})
	config.Handler.AddHandler(&handler.HandleUpdateAppsComplete{})
	config.Handler.AddHandler(&handler.HandleUpdateOwnPresence{})

	go config.StartSession(&wg)

	wg.Wait()
	select {}
}
```

### Multible users of different pbx

Connects to the accounts on two different Hosts

``` GO
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
```

## About ©

[myApps](https://www.innovaphone.com/en/myapps/what-is-myapps.html) is a product of [innovaphone AG](https://www.innovaphone.com).

Documentation of the API used in this client can be found at [ innovaphone App SDK](https://sdk.innovaphone.com/).