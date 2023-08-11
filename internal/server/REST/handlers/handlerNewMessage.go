package handlers

import (
	"github.com/labstack/echo/v4"
	"message_control/internal/message"
	"message_control/internal/storage"
	"net/http"
	"time"
)

type handleMessage struct {
	From string `json:"from"`
	To   string `json:"to"`
	Text string `json:"text"`
}

func handlerNewMessage(storage storage.MessageControl) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if msg, err := ctxToMessage(ctx); err != nil || !storage.AddNewMessage(msg) {
			return ctx.String(http.StatusBadRequest, "Invalid message")
		}
		return ctx.String(http.StatusCreated, "created new message")
	}
}

func ctxToMessage(ctx echo.Context) (message.Message, error) {
	var msg handleMessage
	err := ctx.Bind(&msg)
	return message.Message{
		From: msg.From,
		To:   msg.To,
		Text: msg.Text,
		Date: time.Now(),
		Read: false,
	}, err
}
