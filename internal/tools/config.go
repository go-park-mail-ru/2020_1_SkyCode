package tools

import (
	"github.com/joho/godotenv"
	"os"
)

type Database struct {
	Name     string `json:"name"`
	Host     string `json:"host"`
	Port     string `json:"port"`
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
	if err := godotenv.Load(); err != nil {
		return nil, errors.New("No such env file")
	}

	c := &Config{
		Database: Database{
			Name:     os.Getenv("PSQL_DB"),
			Host:     os.Getenv("PSQL_HOST"),
			Port:     os.Getenv("PSQL_PORT"),
			User:     os.Getenv("PSQL_USER"),
			Password: os.Getenv("PSQL_PASSWORD")},

		ApiKeys: ApiKeys{
			YandexGeoCoder: os.Getenv("YANDEX_GEOCODER"),
		},
	}

	return c, nil
}
