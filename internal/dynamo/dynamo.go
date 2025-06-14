package dynamo

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
)

const (
	UserTablePrefix = "user-table"
	UserIDPrefix    = "User#"

	ReceiverTablePrefix = "receiver-table"
	ReceiverIDPrefix    = "Receiver#"

	localDockerEndpoint = "http://dynamodb-local:8000"
)

func CreateClient(cfg *appconfig.AppConfig) *dynamodb.Client {
	logger := cfg.Logger
	if cfg.Env == appconfig.LocalEnv {
		logger.Info("creating local dynamo db client")
		return createLocalClient(cfg.AWSConfig)
	}
	logger.Info("creating dynamo db client")
	return dynamodb.NewFromConfig(cfg.AWSConfig)
}

func createLocalClient(cfg aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(localDockerEndpoint)
	})
}
