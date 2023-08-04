package storage

import (
	"message_control/internal/message"
)

type Storage interface {
	AddNewMessage(msg message.Message) bool
	GetMessagesChat(user1, user2 string) []message.Message
	GetUsersList(user string) []string
}
