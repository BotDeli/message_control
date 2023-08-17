package message

import (
	"time"
)

type Message struct {
	From string    `json:"from"`
	To   string    `json:"to"`
	Text string    `json:"text"`
	Date time.Time `json:"date"`
}

func (m Message) getDateString() string {
	return m.Date.Format("YYYY:MM:DD")
}
