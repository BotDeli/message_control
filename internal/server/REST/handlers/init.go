package handlers

import "github.com/labstack/echo/v4"

func InitHandlers(client *echo.Echo) {
	client.GET("/", func(ctx echo.Context) error {
		return ctx.String(200, "Hi")
	})
}
