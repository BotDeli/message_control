package postgres

import (
	"database/sql"
	"message_control/errorHandle"
	"message_control/internal/message"
	"message_control/internal/storage"
	"time"
)

func (pg Postgres) AddNewMessage(msg message.Message) bool {
	query := `
	INSERT INTO chat("from_user", "to_user", "text", "date", "read") 
	VALUES ($1, $2, $3, $4, $5)
			`

	_, err := pg.Database.Exec(
		query,
		msg.From,
		msg.To,
		msg.Text,
		msg.Date,
		msg.Read,
	)

	if err != nil {
		errorHandle.Commit(path, "messageControl", "AddNewMessage", err)
	}
	return err == nil
}
func (pg Postgres) GetMessagesChat(username, buddy string) ([]message.Message, error) {
	query := `
	SELECT * 
	FROM chat 
	WHERE (from_user = $1 OR to_user = $1) AND (from_user = $2 OR to_user = $2) ORDER BY date
			`

	var (
		rows *sql.Rows
		err  error
	)

	if rows, err = sendQuery(pg, "GetMessagesChat", query, username, buddy); err != nil {
		return nil, err
	}

	var (
		from, to, text string
		date           time.Time
		read           bool
		messages       = make([]message.Message, 0, 10)
	)

	for rows.Next() {
		err = rows.Scan(&from, &to, &text, &date, &read)
		if err != nil {
			continue
		}

		messages = append(messages, message.Message{
			From: from,
			To:   to,
			Text: text,
			Date: date,
			Read: read,
		})
	}
	return messages, nil
}

func (pg Postgres) GetFriendsList(username string) ([]storage.ChatUser, error) {
	var (
		lastInputUsers, lastOutputUsers []storage.ChatUser
		err                             error
	)

	if lastInputUsers, err = getUsersMessagesToMe(pg, username); err != nil {
		errorHandle.Commit(path, "queriesMessageControl", "GetFriendsList", err)
		return nil, err
	}

	if lastOutputUsers, err = getUsersMessagesSendI(pg, username); err != nil {
		errorHandle.Commit(path, "queriesMessageControl", "GetFriendsList", err)
		return nil, err
	}
	return append(lastInputUsers, lastOutputUsers...), nil
}

func getUsersMessagesToMe(pg Postgres, username string) ([]storage.ChatUser, error) {
	query := `
	SELECT from_user, SUM(CASE WHEN read = false THEN 1 ELSE 0 END) as count_false 
	FROM chat 
	WHERE to_user = $1 
	GROUP BY from_user
			`
	var (
		rows *sql.Rows
		err  error
	)

	if rows, err = sendQuery(pg, "getUsersMessagesToMe", query, username); err != nil {
		return nil, err
	}

	var (
		from  string
		count int
		chats = make([]storage.ChatUser, 0, 10)
	)

	for rows.Next() {
		err = rows.Scan(&from, &count)
		if err != nil {
			continue
		}

		chats = append(chats, storage.ChatUser{
			Username: from,
			Read:     count == 0,
		})
	}
	return chats, nil
}

func getUsersMessagesSendI(pg Postgres, username string) ([]storage.ChatUser, error) {
	query := `
	 SELECT DISTINCT to_user  
	 FROM chat
	 WHERE from_user = $1 AND to_user NOT IN (
	 		SELECT from_user 
    		FROM chat
    		WHERE to_user = $2 
    		GROUP BY from_user
    	)
    		`

	var (
		rows *sql.Rows
		err  error
	)

	if rows, err = sendQuery(pg, "getUsersMessagesSendI", query, username, username); err != nil {
		return nil, err
	}

	var (
		to    string
		chats = make([]storage.ChatUser, 0, 10)
	)

	for rows.Next() {
		err = rows.Scan(&to)
		if err != nil {
			continue
		}

		chats = append(chats, storage.ChatUser{
			Username: to,
			Read:     false,
		})
	}
	return chats, nil
}

func sendQuery(pg Postgres, nameFunction, query string, args ...any) (*sql.Rows, error) {
	rows, err := pg.Database.Query(query, args...)
	if err != nil {
		errorHandle.Commit(path, "queriesMessageControl", nameFunction, err)
	}
	return rows, err
}
