package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Addr string
}

type DatabaseConfig struct {
	ConnString   string
	Username     string
	Password     string
	Host         string
	Port         string
	DatabaseName string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, err
	}

	return Config{
		Server: ServerConfig{
			Addr: os.Getenv("SERVER_ADDRESS"),
		},
		Database: DatabaseConfig{
			ConnString:   os.Getenv("POSTGRES_CONN"),
			Username:     os.Getenv("POSTGRES_USERNAME"),
			Password:     os.Getenv("POSTGRES_PASSWORD"),
			Host:         os.Getenv("POSTGRES_HOST"),
			Port:         os.Getenv("POSTGRES_POST"),
			DatabaseName: os.Getenv("POSTGRES_DATABASE"),
		},
	}, nil

}
