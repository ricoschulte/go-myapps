package appservice

import "fmt"

// the interface all src handlers must implement
type AppServiceCallbackHandler interface {
	HandleCallbackMessage(*AppServiceClient, []byte) error // the function that gets called when a messages received that matches the src
}

type AppServiceCallbackHandlerRegisterItem struct {
	NumCallbacks      int
	ReceivedCallbacks int
	Handler           AppServiceCallbackHandler
}

type AppServiceCallbackHandlerRegister struct {
	Handler map[string]AppServiceCallbackHandlerRegisterItem
}

func NewAppServiceCallbackHandlerRegister() *AppServiceCallbackHandlerRegister {
	return &AppServiceCallbackHandlerRegister{
		Handler: map[string]AppServiceCallbackHandlerRegisterItem{},
	}
}

func (cbr *AppServiceCallbackHandlerRegister) Add(src string, numcallbacks int, handler AppServiceCallbackHandler) error {
	cbr.Handler[src] = AppServiceCallbackHandlerRegisterItem{
		NumCallbacks:      numcallbacks,
		ReceivedCallbacks: 0,
		Handler:           handler,
	}
	return nil
}

func (cbr *AppServiceCallbackHandlerRegister) Remove(src string) error {
	_, ok := cbr.Handler[src]
	if ok {
		// there is a handler for that src
		delete(cbr.Handler, src)
		return nil
	} else {
		// no handler registered for src
		return fmt.Errorf("key not found in handlers: %s", src)
	}

}

func (cbr *AppServiceCallbackHandlerRegister) HandleMessage(appserviceclient *AppServiceClient, src string, message []byte) error {
	handler_item, has_handler := cbr.Handler[src]
	if has_handler {
		handler_item.Handler.HandleCallbackMessage(appserviceclient, message)
		handler_item.ReceivedCallbacks += 1
		if handler_item.ReceivedCallbacks >= handler_item.NumCallbacks {
			cbr.Remove(src)
		}
		return nil
	} else {
		err := fmt.Errorf("no handler found for SRC '%v'", src)
		return err
	}
}
