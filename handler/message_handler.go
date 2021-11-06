package handler

type MessageHandler interface {
	SendMessage(string)
	Close()
}
