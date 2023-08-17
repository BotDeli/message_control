package app

import (
	"log"
	"message_control/internal/config"
	"message_control/internal/server/grpc"
	"message_control/internal/storage/postgres"
)

func StartApplication() {
	cfg := config.MustReadConfig()

	controller := postgres.MustNewMessageControl(cfg.Postgres)

	log.Fatal(grpc.StartServer(cfg.GRPCServer, controller))
}
