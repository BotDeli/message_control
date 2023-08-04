package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	HttpServer HTTPServerConfig `yaml:"http-server" env-required:"true"`
	Postgres   PostgresConfig   `yaml:"postgres" env-required:"true"`
}

type HTTPServerConfig struct {
	Address           string        `yaml:"address" env-default:"localhost:8080"`
	ReadHeaderTimeout time.Duration `yaml:"read-header-timeout" env-default:"10s"`
	IdleTimeout       time.Duration `yaml:"idle-timeout" env-default:"1m"`
}

const (
	pathIsNotSet = "Config path is not set"
	fileNotFound = "Config file not found"
	errorReading = "Error reading config file"
)

func MustReadConfig() Config {
	path := os.Getenv("ConfigPath")
	checkPathIsSet(path)
	checkFileExists(path)
	cfg := readConfig(path)
	return cfg
}

func checkPathIsSet(path string) {
	if path == "" {
		log.Fatal(pathIsNotSet)
	}
}

func checkFileExists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal(fileNotFound)
	}
}

func readConfig(path string) Config {
	var cfg Config
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		log.Fatal(errorReading)
	}
	return cfg
}
