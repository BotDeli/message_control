package handlers

import (
	"github.com/labstack/echo/v4"
	"message_control/internal/storage"
)

func handlerGetUserList(storage storage.MessageControl) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		//username, err := requestToUsername(ctx)
		//if err != nil {
		//	return err
		//}
		//return ctx.JSON(200, storage.GetUsersList(username))
		return nil
	}
}

func requestToUsername(ctx echo.Context) (string, error) {
	var user struct {
		Username string `json:"username"`
	}
	err := ctx.Bind(&user)
	return user.Username, err
}
