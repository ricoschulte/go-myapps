package connection

import "fmt"

// the interface all src handlers must implement
type CallbackHandler interface {
	HandleCallbackMessage(*MyAppsConnection, []byte) error // the function that gets called when a messages received that matches the src
}

type CallbackHandlerRegisterItem struct {
	NumCallbacks      int
	ReceivedCallbacks int
	Handler           CallbackHandler
}

type CallbackHandlerRegister struct {
	Handler map[string]CallbackHandlerRegisterItem
}

func NewCallbackHandlerRegister() *CallbackHandlerRegister {
	return &CallbackHandlerRegister{
		Handler: map[string]CallbackHandlerRegisterItem{},
	}
}

func (cbr *CallbackHandlerRegister) Add(src string, numcallbacks int, handler CallbackHandler) error {
	cbr.Handler[src] = CallbackHandlerRegisterItem{
		NumCallbacks:      numcallbacks,
		ReceivedCallbacks: 0,
		Handler:           handler,
	}
	return nil
}

func (cbr *CallbackHandlerRegister) Remove(src string) error {
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

func (cbr *CallbackHandlerRegister) HandleMessage(myAppsConnection *MyAppsConnection, src string, message []byte) error {
	handler_item, has_handler := cbr.Handler[src]
	if has_handler {
		handler_item.Handler.HandleCallbackMessage(myAppsConnection, message)
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
