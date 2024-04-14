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
	BUILD_COMMIT  string               // The commit hash of the build, useful for tracking specific builds in version control.
	BUILD_VERSION string = "local"     // The version of the build, defaults to the value in DEFAULT_VERSION.
	APP_NAME      string = "TheTipTop" // The name of the application, defaults to the value in DEFAULT_APP_NAME.
	HOSTNAME      string = "localhost" // The hostname of the server, used for generating TLS certificates.
	cfg           *Config

	DEFAULT_DB_NAME string = "default"
	DEFAULT_CONFIG  string = "s3://config.kodmain/config.yml"

	PORT_HTTP  string = ":80"
	PORT_HTTPS string = ":443"
)

type Config struct {
	Mail      *mail.Service       `yaml:"mail"`
	Databases *database.Databases `yaml:"databases"`
	JWT       *jwt.JWT            `yaml:"jwt"`
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

func Load(path string) error {
	if path == "" {
		return fmt.Errorf("path is required")
	}

	var fileContents []byte
	var err error

	if strings.HasPrefix(path, "s3://") {
		fileContents, err = loadFromS3(path)
		if err != nil {
			return err
		}
	} else {
		fileContents, err = os.ReadFile(path)
		if err != nil {
			return err
		}
	}

	configuration := &Config{}
	err = yaml.Unmarshal(fileContents, configuration)
	if err != nil {
		return err
	}

	cfg = configuration

	return cfg.Initialize()
}

func (cfg *Config) Initialize() error {
	if err := database.New(cfg.Databases); err != nil {
		return err
	}

	if err := mail.New(cfg.Mail); err != nil {
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
