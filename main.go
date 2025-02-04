package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
	"github.com/care-giver-app/care-giver-api/internal/awsconfig"
	"github.com/care-giver-app/care-giver-api/internal/dynamo"
	"github.com/care-giver-app/care-giver-api/internal/handlers"
	"github.com/care-giver-app/care-giver-api/internal/log"
	"github.com/care-giver-app/care-giver-api/internal/repository"
	"github.com/care-giver-app/care-giver-api/internal/response"
)

const (
	functionName = "care-taker-api"

	userPath                   = "/user"
	userReceiversPath          = "/user/receivers"
	userPrimaryReceiverPath    = "/user/primary-receiver"
	userAdditionalReceiverPath = "/user/additional-receiver"

	receiverPath      = "/receiver"
	receiverEventPath = "/receiver/event"
)

var (
	dynamoClient *dynamodb.Client
	appCfg       *appconfig.AppConfig
	userRepo     *repository.UserRepository
	receiverRepo *repository.ReceiverRepository
)

func init() {
	appCfg := appconfig.NewAppConfig()

	logger := log.GetLoggerWithEnv(log.InfoLevel, appCfg.Env)
	logger.Sugar().Infof("initializing %s", functionName)

	cfg, err := awsconfig.GetAWSConfig(context.TODO(), appCfg.Env)
	if err != nil {
		logger.Sugar().Fatalf("Unable to load SDK config: %v", err)
	}

	appCfg.Logger = logger
	appCfg.AWSConfig = cfg

	dynamoClient = dynamo.CreateClient(appCfg)
	userRepo = repository.NewUserRespository(context.TODO(), appCfg, dynamoClient)
	receiverRepo = repository.NewReceiverRespository(context.TODO(), appCfg, dynamoClient)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	appCfg.Logger = log.GetLoggerWithEnv(log.InfoLevel, appCfg.Env)
	appCfg.Logger.Info("recieved event")

	ctx = repository.ContextWithUserRespository(ctx, userRepo)
	ctx = repository.ContextWithReceiverRespository(ctx, receiverRepo)

	switch req.Path {
	case userPath:
		return handlers.HandleUser(ctx, appCfg, req)
	case userReceiversPath:
		return handlers.HandleUserReceivers(ctx, appCfg, req)
	case userPrimaryReceiverPath:
		return handlers.HandleUserPrimaryReceiver(ctx, appCfg, req)
	case userAdditionalReceiverPath:
		return handlers.HandleUserAdditionalReceiver(ctx, appCfg, req)
	case receiverPath:
		return handlers.HandleReceiver(ctx, appCfg, req)
	case receiverEventPath:
		return handlers.HandleReceiverEvent(ctx, appCfg, req)
	default:
		return response.CreateBadRequestResponse(), nil
	}
}

func main() {
	lambda.Start(handler)
}
