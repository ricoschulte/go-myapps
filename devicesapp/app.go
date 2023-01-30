package devicesapp

import "github.com/ricoschulte/go-myapps/handler"

type DevicesApp struct {
	Handler *handler.HandleAppService
	Devices map[int]*Device // holds a map of devices by their Id
	Domains map[int]*Domain // holds a map of domains by their Id
}

func NewDevicesApp() *DevicesApp {
	devicesapp := &DevicesApp{}
	devicesapp.Handler = &handler.HandleAppService{
		Name: "devices-api",
	}
	devicesapp.Devices = map[int]*Device{}
	devicesapp.Domains = map[int]*Domain{}

	// register the handler on the appservice
	devicesapp.Handler.MessageHandlerRegister.AddHandler(&HandleAppLoginResult{DevicesApp: devicesapp})
	devicesapp.Handler.MessageHandlerRegister.AddHandler(&HandleGetUserInfoResult{DevicesApp: devicesapp})
	devicesapp.Handler.MessageHandlerRegister.AddHandler(&HandleGetDomainsResult{DevicesApp: devicesapp})
	devicesapp.Handler.MessageHandlerRegister.AddHandler(&HandleGetUnassignedDevicesCountResult{DevicesApp: devicesapp})
	devicesapp.Handler.MessageHandlerRegister.AddHandler(&HandleGetDevicesResult{DevicesApp: devicesapp})
	devicesapp.Handler.MessageHandlerRegister.AddHandler(&HandleDeviceUpdate{DevicesApp: devicesapp})

	// #TODO
	// {"mt":"DeviceRemoved","device":{"id":28,"hwId":"docker0","domainId":1,"product":"Docker Host","version":"v99999sss","type":"VA","pbxActive":false,"online":false}}
	// und DeviceAdded unter garantie auch noch
	// dito für domains
	return devicesapp
}
