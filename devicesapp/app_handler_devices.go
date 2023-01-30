package devicesapp

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ricoschulte/go-myapps/appservice"
)

type HandleAppLoginResult struct {
	DevicesApp *DevicesApp
}

func (m *HandleAppLoginResult) GetMt() string {
	return "AppLoginResult"
}

func (m *HandleAppLoginResult) HandleMessage(appserviceclient *appservice.AppServiceClient, message []byte) error {
	appserviceclient.Printf("%s %v", m.GetMt(), string(message))
	var msgl appservice.AppLoginResult
	if err := json.Unmarshal(message, &msgl); err != nil {
		appserviceclient.Println("error unmarshalling:", m.GetMt(), err)
	}

	if msgl.Ok {
		userinfo, _ := json.Marshal(NewGetUserInfo())
		appserviceclient.Send(userinfo)

		setdevicefilter, _ := json.Marshal(NewSetDeviceFiler(true))
		appserviceclient.Send(setdevicefilter)

		getdomains, _ := json.Marshal(NewGetDomains(true))
		appserviceclient.Send(getdomains)

		getunassigneddevices, _ := json.Marshal(NewGetUnassignedDevicesCount(true))
		appserviceclient.Send(getunassigneddevices)

		// getdevices, _ := json.Marshal(NewGetDevices(true, "", "", "", true))
		// appserviceclient.Send(getdevices)
	}

	return nil
}

type HandleGetUserInfoResult struct {
	DevicesApp *DevicesApp
}

func (m *HandleGetUserInfoResult) GetMt() string {
	return "GetUserInfoResult"
}

func (m *HandleGetUserInfoResult) HandleMessage(appserviceclient *appservice.AppServiceClient, message []byte) error {
	appserviceclient.Printf("%s %v", m.GetMt(), string(message))
	var msgl GetUserInfoResult
	if err := json.Unmarshal(message, &msgl); err != nil {
		appserviceclient.Println("error unmarshalling:", m.GetMt(), err)
	}
	return nil
}

type HandleGetDomainsResult struct {
	DevicesApp *DevicesApp
}

func (m *HandleGetDomainsResult) GetMt() string {
	return "GetDomainsResult"
}

func (m *HandleGetDomainsResult) HandleMessage(appserviceclient *appservice.AppServiceClient, message []byte) error {
	appserviceclient.Printf("%s %v", m.GetMt(), string(message))
	var msgl GetDomainsResult
	if err := json.Unmarshal(message, &msgl); err != nil {
		appserviceclient.Println("error unmarshalling:", m.GetMt(), err)
	}

	listDomainIds := []string{}
	for _, domain := range msgl.Domains {
		listDomainIds = append(listDomainIds, fmt.Sprint(domain.Id))
		m.DevicesApp.Domains[domain.Id] = &domain
	}

	getdevices, _ := json.Marshal(NewGetDevices(true, strings.Join(listDomainIds, ","), "", "", false))
	appserviceclient.Send(getdevices)
	return nil
}

type HandleGetUnassignedDevicesCountResult struct {
	DevicesApp *DevicesApp
}

func (m *HandleGetUnassignedDevicesCountResult) GetMt() string {
	return "GetUnassignedDevicesCountResult"
}

func (m *HandleGetUnassignedDevicesCountResult) HandleMessage(appserviceclient *appservice.AppServiceClient, message []byte) error {
	appserviceclient.Printf("%s %v", m.GetMt(), string(message))
	var msgl GetUnassignedDevicesCountResult
	if err := json.Unmarshal(message, &msgl); err != nil {
		appserviceclient.Println("error unmarshalling:", m.GetMt(), err)
	}
	return nil
}

type HandleGetDevicesResult struct {
	DevicesApp *DevicesApp
}

func (m *HandleGetDevicesResult) GetMt() string {
	return "GetDevicesResult"
}

func (m *HandleGetDevicesResult) HandleMessage(appserviceclient *appservice.AppServiceClient, message []byte) error {
	appserviceclient.Printf("%s %v", m.GetMt(), string(message))
	var msgl GetDevicesResult
	if err := json.Unmarshal(message, &msgl); err != nil {
		appserviceclient.Println("error unmarshalling:", m.GetMt(), err)
	}
	appserviceclient.Printf("num devices %v", len(msgl.Devices))
	for _, device := range msgl.Devices {
		m.DevicesApp.Devices[device.Id] = &device
	}
	return nil
}

type HandleDeviceUpdate struct {
	DevicesApp *DevicesApp
}

func (m *HandleDeviceUpdate) GetMt() string {
	return "DeviceUpdate"
}

func (m *HandleDeviceUpdate) HandleMessage(appserviceclient *appservice.AppServiceClient, message []byte) error {
	appserviceclient.Printf("%s %v", m.GetMt(), string(message))
	var msgl DeviceUpdate
	if err := json.Unmarshal(message, &msgl); err != nil {
		appserviceclient.Println("error unmarshalling:", m.GetMt(), err)
	}
	m.DevicesApp.Devices[msgl.Device.Id] = &msgl.Device
	return nil
}
