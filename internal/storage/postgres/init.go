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

type Postgres struct {
	Database DB
}

type DB interface {
	Close() error
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
}

func MustNewMessageControl(cfg config.PostgresConfig) storage.MessageControl {
	dataSourceName := cfg.GetDataSourceName()
	database, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		errorHandle.Commit(path, "init.go", "MustNewStorage", err)
	}

	if err = database.Ping(); err != nil {
		log.Fatal(errConnectionError)
	}

	pg := Postgres{Database: database}
	return pg
}

func (pg Postgres) Disconnect() {
	if err := pg.Database.Close(); err != nil {
		errorHandle.Commit(path, "queriesMessageControl.go", "Disconnect", err)
	}
}
