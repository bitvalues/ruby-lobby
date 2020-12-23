package config

import (
	"os"
	"strconv"
)

type Config struct {
	AuthServerPort int
	DataServerPort int
	ViewServerPort int
}

func GetConfig() Config {
	cfg := Config{
		AuthServerPort: 54231,
		DataServerPort: 54230,
		ViewServerPort: 54001,
	}

	authServerPort, err := strconv.Atoi(os.Getenv("AUTH_SERVER_PORT"))
	if err == nil {
		cfg.AuthServerPort = authServerPort
	}

	dataServerPort, err := strconv.Atoi(os.Getenv("DATA_SERVER_PORT"))
	if err == nil {
		cfg.DataServerPort = dataServerPort
	}

	viewServerPort, err := strconv.Atoi(os.Getenv("VIEW_SERVER_PORT"))
	if err == nil {
		cfg.ViewServerPort = viewServerPort
	}

	return cfg
}
