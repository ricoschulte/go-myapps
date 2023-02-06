# Sysclient

An implementation of the Sysclient protocol described at [Sysclient protocol](https://sdk.innovaphone.com/13r3/doc/protocol/sysclient.htm)

- connects a Go client to the innovaphone Devices App as a sysclient like any other device.
- can be used to provide a web interface with text or other functions.
- or to record or check the configurations devices receive from the Devices app.
- useful during the development of myApps apps when they have to have access to devices. Allows to use mock/dummy devices instead of having to maintain real ones.

## Example usage of the libary

This example code deploys a sysclient that connects to the Devices app's websocket.
It then serves HTTP endpoints that are accessible through the interface in the Devices app.
It issues the requests with the received configurations on the console.

``` go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ricoschulte/go-myapps/sysclient"
)



func main() {
	// define the properties of the device that connects to the devices app as sysclient
	identity := sysclient.Identity{
		Id:      "f19033480af9",
		Product: "IP232",
		Version: "13r2 dvl [13.4250/131286/1300]",
		FwBuild: "134250",
		BcBuild: "131286",
		Major:   "13r2",
		Fw:      "ip222.bin",
		Bc:      "boot222.bin",
		Mini:    false,
		PbxActive: false,
		Other: false,
		Platform: sysclient.Platform{
			Type: "PHONE",
		},
		EthIfs: []sysclient.EthIf{
			{
				If:   "ETH0",
				Ipv4: "172.16.4.141",
				Ipv6: "2002:91fd:9d07:0:290:33ff:fe46:af2",
			},
		},
	}

	sc, err_creating_sysclient := sysclient.NewSysclient(
		// the identity from above
		identity,

		
		// the Devices App URL to connect to
        // ws[s]://<ip/host>[:<port>]/<domain></instance>/sysclients
		"wss://apps.company.com/company.com/devices/sysclients",

		// a timeout duration for the websocket
		time.Duration(2*time.Second),

		// InsecureSkipVerify: if true, the verification of the TLS certificate is skipped.
		// Do not use in production. You have been warned.
		false,

		// a SeveMux for requested HTTP request, see below
		getServeMux(),

		// filenames to store the password of the sysclient and a received and decoded admin password
		"sysclient_password.txt",
		"sysclient_administrativepassword.txt",
	)

	// when no error has happend while creating the client ...
	if err_creating_sysclient != nil {
		panic(err_creating_sysclient)
	}
	// ... start it.
	sc.Connect()

	select {} // keep the program running
}

// create a ServeMux for handling Tunnel Http Requests
func getServeMux() *http.ServeMux {
	mux := http.NewServeMux()

	// Serve a file for the /admin.xml path
	mux.HandleFunc("/admin.xml", func(w http.ResponseWriter, r *http.Request) {
		response_text := `<html><body>`
		response_text += fmt.Sprintf(`page on %s`, r.URL.Path)
		response_text += "</body></html>"

		w.Header().Add("Content-Length", fmt.Sprint(len(response_text)))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response_text))
	})

	// Serve a file for the /CMD0/mod_cmd.xml endpoint the will receive configurations
	mux.HandleFunc("/CMD0/mod_cmd.xml", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("---------------------------")
		fmt.Println("received on", r.URL.Path)
		fmt.Println("received headers", r.Header)
		for key, value := range r.URL.Query() {
			fmt.Println("received param", key, value)
		}

		response_text := `<html><body>`
		response_text += fmt.Sprintf(`mod_cmd.xml page on %s`, r.URL.Path)
		response_text += "</body></html>"

		w.Header().Add("Content-Length", fmt.Sprint(len(response_text)))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response_text))
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			// catch not existing path with error log
			response_text := fmt.Sprintf("Path %s not found", r.URL.Path)
			fmt.Println("REQUEST 404", r.URL.Path)
			w.Header().Add("Content-Length", fmt.Sprint(len(response_text)))
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(response_text))

			return
		}

		response_text := `<html><body>`
		response_text += fmt.Sprintf(`page on %s`, r.URL.Path)
		response_text += "</body></html>"

		w.Header().Add("Content-Length", fmt.Sprint(len(response_text)))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response_text))

	})

	return mux
}
```