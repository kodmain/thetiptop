package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/aws/s3"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
)

var (
	cfg *Config
)

type Config struct {
	Providers struct {
		Mails     map[string]*mail.Config     `yaml:"mails"`
		Databases map[string]*database.Config `yaml:"databases"`
	} `yaml:"providers"`
	JWT *jwt.JWT `yaml:"jwt"`
}

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

func Reset() {
	cfg = nil
}

func Load(path *string) error {
	if path == nil || *path == "" {
		return fmt.Errorf("path is required")
	}

	var fileContents []byte
	var err error

	switch {
	case strings.HasPrefix(*path, "s3://"):
		fileContents, err = loadFromS3(*path)
	default:
		fileContents, err = os.ReadFile(*path)
	}

	if err != nil {
		return err
	}

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

	if err := jwt.New(cfg.JWT); err != nil {
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
