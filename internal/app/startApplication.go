package app

import (
	"log"
	"message_control/internal/config"
	"message_control/internal/server/serverGRPC"
	"message_control/internal/storage/postgres"
)

func StartApplication() {
	cfg := config.MustReadConfig()

	controller := postgres.MustNewMessageControl(cfg.Postgres)

	log.Fatal(serverGRPC.StartServer(cfg.GRPCServer, controller))
}
