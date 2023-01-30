package appservice

import (
	"fmt"
)

// the interface all message handlers must implement
type AppServiceMessageHandler interface {
	GetMt() string // the MessageHandlerRegister matches this string against the incoming message MT
	HandleMessage(*AppServiceClient, []byte) error
}

type AppServiceMessageHandlerRegister struct {
	Handler []AppServiceMessageHandler
}

func (hr *AppServiceMessageHandlerRegister) AddHandler(handler AppServiceMessageHandler) error {
	hr.Handler = append(hr.Handler, handler)
	return nil
}

func (hr *AppServiceMessageHandlerRegister) HandleMessage(appserviceclient *AppServiceClient, mt string, message []byte) error {
	handled := false

	for _, handler := range hr.Handler {
		if handler.GetMt() == mt {
			handler.HandleMessage(appserviceclient, message)
			handled = true
		}
	}

	if !handled {
		err := fmt.Errorf("no handler found for MT '%s': %v", mt, string(message))
		return err
	} else {
		return nil
	}
}
