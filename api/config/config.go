package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/aws/s3"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
)

const (
	DEFAULT = "default"
)

var (
	cfg *Config
)

type Config struct {
	Services map[string]struct {
		Database string `yaml:"database"`
		Mail     string `yaml:"mail"`
	} `yaml:"services"`
	Providers struct {
		Mails     map[string]*mail.Config     `yaml:"mails"`
		Databases map[string]*database.Config `yaml:"databases"`
	} `yaml:"providers"`
	Security struct {
		Validation struct {
			Expire string `yaml:"expire"`
		} `yaml:"validation"`
		JWT *jwt.JWT `yaml:"jwt"`
	} `yaml:"security"`
}

// Get Retrieve the value from cfg based on the provided key
// Retrieves a value from a config structure by key, following a path syntax (e.g. "parent.child").
// If the value is not found or is nil, it returns the provided defaultValue.
//
// Parameters:
// - key: string The key to access the value, formatted as "parent.child".
// - defaultValue: interface{} The value to return if the key is not found or the value is nil.
//
// Returns:
// - interface{} The retrieved value from cfg, or defaultValue if the key is not found or value is nil.
func Get(key string, defaultValue any) any {
	if cfg == nil {
		return defaultValue
	}

	path := strings.Split(key, ".")
	val := reflect.ValueOf(cfg)

	for _, elem := range path {
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}

		switch val.Kind() {
		case reflect.Struct:
			val = val.FieldByNameFunc(func(name string) bool {
				return strings.EqualFold(elem, name)
			})
		case reflect.Map:
			val = val.MapIndex(reflect.ValueOf(elem))
		default:
			return defaultValue
		}

		if !val.IsValid() {
			return defaultValue
		}
	}

	finalValue := convertValue(val.Interface(), defaultValue)
	return finalValue
}

func convertValue(val any, defaultValue any) any {
	switch defaultValue.(type) {
	case int:
		switch v := val.(type) {
		case string:
			if intValue, err := strconv.Atoi(v); err == nil {
				return intValue
			}
		case int:
			return v
		}
	case string:
		return reflect.ValueOf(val).String()
	case bool:
		switch v := val.(type) {
		case string:
			if boolValue, err := strconv.ParseBool(v); err == nil {
				return boolValue
			}
		case bool:
			return v
		}
	default:
		return val
	}

	return defaultValue
}

func GetInt(key string, defaultValue int) int {
	value := Get(key, defaultValue)

	if intValue, ok := value.(int); ok {
		return intValue
	}

	return defaultValue
}

func GetString(key string, defaultValue string) string {
	value := Get(key, defaultValue)

	if strValue, ok := value.(string); ok {
		return strValue
	}

	return defaultValue
}

func Reset() {
	cfg = nil
}

func Load(path *string) error {
	if path == nil || *path == "" {
		return fmt.Errorf("path is required")
	}

	var fileContents []byte
	var err error
	var workingDir string

	switch {
	case strings.HasPrefix(*path, "s3://"):
		workingDir, err = os.Getwd()
		if err != nil {
			return err
		}

		fileContents, err = loadFromS3(*path)
		if err != nil {
			return err
		}
	default:
		fileContents, err = os.ReadFile(*path)
		if err != nil {
			return err
		}
		abs, err := filepath.Abs(*path)
		if err != nil {
			return err
		}

		workingDir = filepath.Dir(abs)
	}

	fileContents = []byte(strings.ReplaceAll(string(fileContents), "${PWD}", workingDir))

	cfg = &Config{}
	err = yaml.Unmarshal(fileContents, cfg)

	if err != nil {
		return err
	}

	return cfg.Initialize()
}

func (cfg *Config) Initialize() error {
	if err := database.New(cfg.Providers.Databases); err != nil {
		return err
	}

	if err := mail.New(cfg.Providers.Mails); err != nil {
		return err
	}

	if err := jwt.New(cfg.Security.JWT); err != nil {
		return err
	}

	return nil
}

func loadFromS3(s3Path string) ([]byte, error) {
	s3URL := strings.SplitN(s3Path[len("s3://"):], "/", 2)
	if len(s3URL) < 2 {
		return nil, fmt.Errorf("invalid S3 path")
	}

	bucket := s3URL[0]
	item := s3URL[1]

	service, err := s3.New()
	if err != nil {
		return nil, err
	}

	output, err := service.GetObject(&bucket, &item)

	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}
