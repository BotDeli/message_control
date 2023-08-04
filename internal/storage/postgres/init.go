package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"message_control/errorHandle"
	"message_control/internal/config"
	"message_control/internal/storage"
)

const (
	errConnectionError = "connection error"
	path               = "message_control/internal/storage/postgres"
)

type DB interface {
	Close() error
}

type Postgres struct {
	Database DB
}

func MustNewStorage(cfg config.PostgresConfig) storage.Storage {
	dataSourceName := cfg.GetDataSourceName()
	database, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		log.Fatal(errorHandle.ErrorFormatString(path, "init.go", "MustNewStorage", err))
	}

	if err = database.Ping(); err != nil {
		log.Fatal(errConnectionError)
	}

	pg := Postgres{Database: database}
	return pg
}

func (pg Postgres) Disconnect() {
	if err := pg.Database.Close(); err != nil {
		log.Println(errorHandle.ErrorFormatString(path, "storage.go", "Disconnect", err))
	}
}
