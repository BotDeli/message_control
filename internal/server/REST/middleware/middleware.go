package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func CheckUUIDRequest(expectedUUID string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			requestUUID := ctx.Request().FormValue("uuid")
			if requestUUID == expectedUUID {
				return next(ctx)
			}
			return ctx.JSON(http.StatusUnauthorized, nil)
		}
	}
}
