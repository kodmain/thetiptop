package aws

import (
	"context"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var instance *aws.Config
var CTX = context.Background()

func Connect(profil ...string) (*aws.Config, error) {
	if instance == nil {
		optFns := []func(*config.LoadOptions) error{}
		for _, p := range profil {
			optFns = append(optFns, config.WithSharedConfigProfile(p))
		}

		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			client := &http.Client{Timeout: 500 * time.Millisecond}
			_, err := client.Get("http://169.254.169.254/latest/meta-data/")
			if err == nil {
				// L'application s'exécute sur une instance EC2, utiliser l'adresse IP locale
				return aws.Endpoint{
					URL: "http://169.254.169.254",
				}, nil
			}
			// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})

		optFns = append(optFns, config.WithEndpointResolverWithOptions(customResolver))

		cfg, err := config.LoadDefaultConfig(CTX, optFns...)
		if err != nil {
			return nil, err
		}

		instance = &cfg
	}

	return instance, nil
}
