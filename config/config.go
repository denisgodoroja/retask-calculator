package config

import (
	"os"
)

type Config struct {
	Server   Server
	Database Database
}

type Server struct {
	Port string
}

type Database struct {
	Host     string
	User     string
	Password string
	Name     string
}

func LoadConfigs() (*Config, error) {
	cfg := &Config{
		Server: Server{
			Port: os.Getenv("PORT"),
		},
		Database: Database{
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Host:     os.Getenv("DB_HOST"),
			Name:     os.Getenv("DB_DATABASE"),
		},
	}

	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}

	return cfg, nil
}
