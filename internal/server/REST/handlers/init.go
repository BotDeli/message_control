package handlers

import (
	"github.com/labstack/echo/v4"
	"message_control/internal/storage"
)

func InitHandlers(client *echo.Echo, storage storage.MessageControl) {
	client.POST("/newMessage", handlerNewMessage(storage))
	client.GET("/userList", handlerGetUserList(storage))
}
