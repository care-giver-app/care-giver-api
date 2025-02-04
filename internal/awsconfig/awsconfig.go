package awsconfig

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
)

var (
	USEastRegion = "us-east-1"
)

func GetAWSConfig(ctx context.Context, env string) (aws.Config, error) {
	if env == appconfig.LocalEnv {
		return getLocalAWSConfig(ctx)
	}
	return getAWSConfig(ctx)
}

func getAWSConfig(ctx context.Context) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(USEastRegion),
	)

	return cfg, err
}

func getLocalAWSConfig(ctx context.Context) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(USEastRegion),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)

	return cfg, err
}
