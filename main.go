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
	dynamoClient *dynamodb.Client
	appCfg       *appconfig.AppConfig
	userRepo     *repository.UserRepository
	receiverRepo *repository.ReceiverRepository

	PathToHandlerMap = map[string]func(ctx context.Context, appCfg *appconfig.AppConfig, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
		"/user":                     handlers.HandleCreateUser,
		"/user/{userId}":            handlers.HandleGetUser,
		"/user/primary-receiver":    handlers.HandleUserPrimaryReceiver,
		"/user/additional-receiver": handlers.HandleUserAdditionalReceiver,
		"/receiver/{receiverId}":    handlers.HandleReceiver,
		"/receiver/event":           handlers.HandleReceiverEvent,
	}
)

func init() {
	appCfg = appconfig.NewAppConfig()

	logger := log.GetLoggerWithEnv(log.InfoLevel, appCfg.Env)
	logger.Sugar().Infof("initializing %s", functionName)

	cfg, err := awsconfig.GetAWSConfig(context.TODO(), appCfg.Env)
	if err != nil {
		logger.Sugar().Fatalf("Unable to load SDK config: %v", err)
	}

	appCfg.Logger = logger
	appCfg.AWSConfig = cfg

	dynamoClient = dynamo.CreateClient(appCfg)

	logger.Sugar().Info("initializing user respository")
	userRepo = repository.NewUserRespository(context.TODO(), appCfg, dynamoClient)

	logger.Sugar().Info("initializing receiver respository")
	receiverRepo = repository.NewReceiverRespository(context.TODO(), appCfg, dynamoClient)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	appCfg.Logger.Info("recieved event")

	ctx = repository.ContextWithUserRespository(ctx, userRepo)
	ctx = repository.ContextWithReceiverRespository(ctx, receiverRepo)

	appCfg.Logger.Sugar().Infof("event path: %s", req.RequestContext.Path)

	//TODO: Create handler param struct to pass in repositories as well
	if handler, ok := PathToHandlerMap[req.RequestContext.Path]; ok {
		return handler(ctx, appCfg, req)
	}

	appCfg.Logger.Error("unsuported request path", zap.Any(log.PathLogKey, req.RequestContext.Path))
	return response.CreateBadRequestResponse(), nil
}

func main() {
	lambda.Start(handler)
}
