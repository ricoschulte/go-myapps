package connection

import (
	"fmt"
)

// the interface all message handlers must implement
type MessageHandler interface {
	GetMt() string // the MessageHandlerRegister matches this string against the incoming message MT
	HandleMessage(*MyAppsConnection, []byte) error
}

type MessageHandlerRegister struct {
	Handler []MessageHandler
}

func (hr *MessageHandlerRegister) AddHandler(handler MessageHandler) error {
	hr.Handler = append(hr.Handler, handler)
	return nil
}

func (hr *MessageHandlerRegister) HandleMessage(myAppsConnection *MyAppsConnection, key string, message []byte) error {
	for _, handler := range hr.Handler {
		if handler.GetMt() == key {
			return handler.HandleMessage(myAppsConnection, message)
		}
	}
	err := fmt.Errorf("no handler found for '%v'", key)
	return err
}