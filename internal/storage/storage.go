package storage

import (
	"message_control/internal/message"
)

type MessageControl interface {
	AddNewMessage(msg message.Message) bool
	GetMessagesChat(username, buddy string) ([]message.Message, error)
	GetUsersList(username string) ([]ChatUser, error)
}

type ChatUser struct {
	Username string `json:"username"`
	Read     bool   `json:"read"`
}
