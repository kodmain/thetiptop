package config

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gopkg.in/yaml.v3"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
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

	var fileContents []byte
	var err error

	if strings.HasPrefix(path, "s3://") {
		fileContents, err = loadFromS3(context.Background(), path)
		if err != nil {
			return nil, err
		}
	} else {
		fileContents, err = os.ReadFile(path)
		if err != nil {
			return nil, err
		}
	}

	configuration := &Config{}
	err = yaml.Unmarshal(fileContents, configuration)
	if err != nil {
		return nil, err
	}

	cfg = configuration

	return configuration, nil
}

func loadFromS3(ctx context.Context, s3Path string) ([]byte, error) {
	s3URL := strings.SplitN(s3Path[len("s3://"):], "/", 2)
	if len(s3URL) < 2 {
		return nil, fmt.Errorf("invalid S3 path")
	}

	bucket := s3URL[0]
	item := s3URL[1]

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	output, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &item,
	})
	if err != nil {
		return nil, err
	}
	defer output.Body.Close()

	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, output.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
