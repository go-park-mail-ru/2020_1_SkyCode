package tools

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Database struct {
		Name     string `json:"name"`
		Host     string `json:"host"`
		Port     uint16 `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"database"`

	Server struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"server"`
}

func LoadConf(filename string) (*Config, error) {
	c := &Config{}

	file, err := os.Open(filename)

	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err != nil {
		return nil, err
	}

	parser := json.NewDecoder(file)

	if err := parser.Decode(c); err != nil {
		return nil, err
	}

	return c, nil
}
