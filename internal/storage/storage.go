package storage

import (
	"message_control/internal/message"
	"message_control/internal/server/serverGRPC/pb"
)

type MessageControl interface {
	AddNewMessage(msg message.Message) (bool, error)
	GetFriendsList(username string) ([]*pb.Friend, error)
	GetMessagesChat(username, friend string) ([]*pb.ChatMessage, error)
}
