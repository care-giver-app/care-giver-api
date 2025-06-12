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
	"go.uber.org/zap"
)

const (
	functionName = "care-taker-api"
)

var (
	dynamoClient    *dynamodb.Client
	appCfg          *appconfig.AppConfig
	userRepo        *repository.UserRepository
	receiverRepo    *repository.ReceiverRepository
	eventRepo       *repository.EventRepository
	handlerRegistry handlers.RegistryProvider
)

func init() {
	appCfg = appconfig.NewAppConfig()
	appCfg.Logger.Sugar().Infof("initializing %s", functionName)

	cfg, err := awsconfig.GetAWSConfig(context.TODO(), appCfg.Env)
	if err != nil {
		appCfg.Logger.Sugar().Fatalf("Unable to load SDK config: %v", err)
	}

	appCfg.AWSConfig = cfg

	dynamoClient = dynamo.CreateClient(appCfg)

	appCfg.Logger.Info("initializing user respository")
	userRepo = repository.NewUserRespository(context.TODO(), appCfg, dynamoClient)

	appCfg.Logger.Info("initializing receiver respository")
	receiverRepo = repository.NewReceiverRespository(context.TODO(), appCfg, dynamoClient)

	appCfg.Logger.Info("initializing event repository")
	eventRepo = repository.NewEventRespository(context.TODO(), appCfg, dynamoClient)

	appCfg.Logger.Info("initializing handler registry")
	handlerRegistry = handlers.NewRegistry(appCfg, userRepo, receiverRepo, eventRepo)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	appCfg.Logger.Info("recieved event")

	if handler, ok := handlerRegistry.GetHandler(req); ok {
		return handlerRegistry.RunHandler(ctx, handler, req)
	}

	appCfg.Logger.Error("unsupported request path", zap.Any(log.PathLogKey, req.RequestContext.ResourcePath), zap.Any(log.MethodLogKey, req.HTTPMethod))
	return response.CreateBadRequestResponse(), nil
}

func main() {
	lambda.Start(handler)
}
