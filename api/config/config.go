package config

import (
	"fmt"
	"os"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Mail      *mail.Service       `yaml:"mail"`
	Databases *database.Databases `yaml:"databases"`
}

func Load(path string) (*Config, error) {
	if path == "" {
		return nil, fmt.Errorf("path is required")
	}

	fileContents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	configuration := &Config{}
	err = yaml.Unmarshal(fileContents, configuration)
	if err != nil {
		return nil, err
	}

	return configuration, nil
}
