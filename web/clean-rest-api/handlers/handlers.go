package handlers

import (
	"clean-rest-api/model"
)

type MessageHandler struct {
	messages []model.Message
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		messages: []model.Message{
			{Id: 1, Title: "demo", Body: "demo message"},
			{Id: 2, Title: "random", Body: "random message"},
		},
	}
}

func (m *MessageHandler) GetAllMessages() []model.Message {
	return m.messages
}
