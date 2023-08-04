package postgres

import "message_control/internal/message"

func (pg Postgres) AddNewMessage(msg message.Message) bool {
	return false
}
func (pg Postgres) GetMessagesChat(user1, user2 string) []message.Message {
	return nil
}
func (pg Postgres) GetUsersList(user string) []string {
	return nil
}
