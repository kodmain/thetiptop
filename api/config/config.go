package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Mail      *mail.Service       `yaml:"mail"`
	Databases *database.Databases `yaml:"databases"`
	JWT       *jwt.JWT            `yaml:"jwt"`
}

var (
	cfg *Config
)

func Get(key string) interface{} {
	if cfg == nil {
		return nil
	}

	path := strings.Split(key, ".")
	val := reflect.ValueOf(cfg)

	for _, elem := range path {
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		if val.Kind() == reflect.Struct {
			val = val.FieldByNameFunc(func(name string) bool {
				return strings.EqualFold(elem, name)
			})
		}

		if !val.IsValid() || (val.Kind() == reflect.Ptr && val.IsNil()) {
			return nil
		}
	}

	// Check if the value is a zero value.
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.IsValid() && !val.IsZero() {
		return val.Interface()
	}

	return nil
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

	cfg = configuration

	return configuration, nil
}
