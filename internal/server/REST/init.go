package REST

import (
	"github.com/labstack/echo/v4"
	"log"
	"message_control/internal/config"
	"message_control/internal/server/REST/handlers"
	"message_control/internal/storage"
	"net/http"
)

func MustStartServer(cfg config.HTTPServerConfig, storage storage.Storage) {
	client := echo.New()
	handlers.InitHandlers(client)
	server := http.Server{
		Addr:              cfg.Address,
		Handler:           client,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		IdleTimeout:       cfg.IdleTimeout,
	}
	log.Fatal(server.ListenAndServe())
}
