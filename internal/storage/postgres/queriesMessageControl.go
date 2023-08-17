package postgres

import (
	"database/sql"
	"fmt"
	"message_control/errorHandle"
	"message_control/internal/message"
	"message_control/internal/server/serverGRPC/pb"
	"time"
)

func (pg Postgres) AddNewMessage(msg message.Message) (bool, error) {
	query := `
	INSERT INTO chat("from_user", "to_user", "text", "date") 
	VALUES ($1, $2, $3, $4)
			`

	_, err := pg.Database.Exec(
		query,
		msg.From,
		msg.To,
		msg.Text,
		msg.Date,
	)

	if err != nil {
		errorHandle.Commit(path, "messageControl", "AddNewMessage", err)
	}
	return err == nil, err
}
func (pg Postgres) GetMessagesChat(username, buddy string) ([]*pb.ChatMessage, error) {
	query := `
	SELECT * 
	FROM chat 
	WHERE (from_user = $1 OR to_user = $1) 
	  AND (from_user = $2 OR to_user = $2) 
	ORDER BY date DESC
			`

	var (
		rows *sql.Rows
		err  error
	)

	messages := make([]*pb.ChatMessage, 0, 10)

	if rows, err = sendQuery(pg, "GetMessagesChat", query, username, buddy); err != nil {
		return messages, err
	}

	var (
		from, to, text string
		date           time.Time
	)

	for rows.Next() {
		err = rows.Scan(&from, &to, &text, &date)
		if err != nil {
			fmt.Println(err)
			continue
		}

		messages = append(messages, &pb.ChatMessage{
			Msg: &pb.BodyMessage{
				From: from,
				To:   to,
				Text: text,
			},
			Date: dateFormat(date),
		})
	}
	return messages, nil
}

func (pg Postgres) GetFriendsList(username string) ([]*pb.Friend, error) {
	query := `SELECT DISTINCT 
    		CASE 
    		    WHEN from_user = $1 THEN to_user 
    		    WHEN to_user = $1 THEN from_user 
    		    ELSE $1 
    		END as username,
    		MAX(date)
			FROM chat
			GROUP BY username
			ORDER BY MAX(date) DESC
		`

	rows, err := sendQuery(pg, "GetFriendsList", query, username)

	friends := make([]*pb.Friend, 0, 10)

	if err != nil {
		return friends, err
	}

	var (
		friendName string
		date       time.Time
	)

	for rows.Next() {
		err = rows.Scan(&friendName, &date)
		if err == nil && friendName != "" && friendName != username {
			friends = append(friends,
				&pb.Friend{
					Username: friendName,
					Date:     dateFormat(date),
				},
			)
		}
	}

	return friends, nil
}

func sendQuery(pg Postgres, nameFunction, query string, args ...any) (*sql.Rows, error) {
	rows, err := pg.Database.Query(query, args...)
	if err != nil {
		errorHandle.Commit(path, "queriesMessageControl", nameFunction, err)
	}
	return rows, err
}

func dateFormat(date time.Time) string {
	return date.Format("02.01.2006")
}
