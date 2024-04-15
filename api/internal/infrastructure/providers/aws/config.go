package aws

import (
	"context"

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

		cfg, err := config.LoadDefaultConfig(CTX, optFns...)
		if err != nil {
			return nil, err
		}

		instance = &cfg
	}

	return instance, nil
}
