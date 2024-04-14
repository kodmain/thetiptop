package s3

import (
	"bytes"
	"context"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/aws"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/buffer"
)

var instance *Service

type ServiceInterface interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

type Service struct {
	API ServiceInterface
}

func New() (*Service, error) {
	if instance != nil {
		return instance, nil
	}

	cfg, err := aws.Connect()
	if err != nil {
		return nil, err
	}

	instance = &Service{
		API: s3.NewFromConfig(*cfg),
	}

	return instance, nil
}

func (s *Service) GetObject(bucket *string, item *string) (*bytes.Buffer, error) {
	output, err := s.API.GetObject(aws.CTX, &s3.GetObjectInput{
		Bucket: bucket,
		Key:    item,
	})

	if err != nil {
		return nil, err
	}

	return buffer.Read(output.Body)
}
