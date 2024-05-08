package aws

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/kodmain/thetiptop/api/env"
)

var instance *aws.Config
var CTX = context.Background()

func Connect(profil ...string) (*aws.Config, error) {
	if instance == nil {

		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true, // Ignorer la vérification des certificats TLS
				},
			},
		}

		optFns := []func(*config.LoadOptions) error{
			config.WithDefaultRegion("eu-west-3"),
			config.WithSharedConfigProfile(env.DEFAULT_AWS_PROFILE),
			config.WithHTTPClient(httpClient),
		}

		for _, p := range profil {
			optFns = append(optFns, config.WithSharedConfigProfile(p))
		}

		cfg, err := config.LoadDefaultConfig(CTX, optFns...)
		if err != nil {
			return nil, err
		}

		instance = &cfg
	}

	return instance, nil
}
