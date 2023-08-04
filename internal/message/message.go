package message

import "time"

type Message struct {
	Text string    `json:"text"`
	Time time.Time `json:"time"`
	From string    `json:"from"`
	To   string    `json:"to"`
}
