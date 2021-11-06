package handler

import "log"

type dummyMessageHandler struct {
}

func (handler *dummyMessageHandler) SendMessage(message string) {
	log.Println("Message received:", message)
}

func (handler *dummyMessageHandler) Close() {

}

//NewDummyMessageHandler Factory function for creating new DummyMessageHandler
func NewDummyMessageHandler() (*dummyMessageHandler, error) {
	handler := new(dummyMessageHandler)
	return handler, nil
}
