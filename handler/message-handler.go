package handler

type MessageHandler interface {
	Message(string)
	Close()
}
