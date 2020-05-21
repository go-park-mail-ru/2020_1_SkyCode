package tools

import (
	"os"
	"strconv"
)

type Database struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     uint64 `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type ApiKeys struct {
	YandexGeoCoder string `json:"yandex_geocoder"`
}
type Config struct {
	Database Database
	Server   Server
	ApiKeys  ApiKeys
}

func LoadConf() (*Config, error) {
	port, err := strconv.ParseUint(os.Getenv("PSQL_PORT"), 10, 64)

	if err != nil {
		return nil, err
	}
	c := &Config{
		Database: Database{
			Name:     os.Getenv("PSQL_DB"),
			Host:     os.Getenv("PSQL_HOST"),
			Port:     port,
			User:     os.Getenv("PSQL_USER"),
			Password: os.Getenv("PSQL_PASSWORD")},

		ApiKeys: ApiKeys{
			YandexGeoCoder: os.Getenv("YANDEX_GEOCODER"),
		},
	}

	return c, nil
}
