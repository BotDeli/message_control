package config

import "fmt"

type PostgresConfig struct {
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-default:""`
	Host     string `yaml:"host" env-default:"localhost"`
	Dbname   string `yaml:"dbname" env-required:"true"`
	Sslmode  string `yaml:"sslmode" env-default:"disable"`
}

func (config *PostgresConfig) GetDataSourceName() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		config.User,
		config.Password,
		config.Host,
		config.Dbname,
		config.Sslmode)
}
