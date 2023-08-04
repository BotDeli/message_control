package app

import (
	"message_control/internal/config"
	"message_control/internal/server/REST"
	"message_control/internal/storage/postgres"
)

func StartApplication() {
	cfg := config.MustReadConfig()
	storage := postgres.MustNewStorage(cfg.Postgres)
	REST.MustStartServer(cfg.HttpServer, storage)
}
