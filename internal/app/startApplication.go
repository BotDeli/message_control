package app

import (
	"message_control/internal/config"
	"message_control/internal/server/REST"
	"message_control/internal/storage/postgres"
)

func StartApplication() {
	cfg := config.MustReadConfig()

	controller := postgres.MustNewMessageControl(cfg.Postgres)

	REST.MustStartServer(cfg.HttpServer, controller)
}
