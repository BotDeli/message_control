package REST

import (
	"github.com/labstack/echo/v4"
	"log"
	"message_control/internal/config"
	"message_control/internal/server/REST/handlers"
	"message_control/internal/server/REST/middleware"
	"message_control/internal/storage"
	"net/http"
)

func MustStartServer(cfg config.HTTPServerConfig, storage storage.MessageControl) {
	client := echo.New()
	client.Use(middleware.CheckUUIDRequest(cfg.UUID))
	handlers.InitHandlers(client, storage)
	startServer(cfg, client)
}

func startServer(cfg config.HTTPServerConfig, client *echo.Echo) {
	server := http.Server{
		Addr:              cfg.Address,
		Handler:           client,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
	log.Fatal(server.ListenAndServe())
}
