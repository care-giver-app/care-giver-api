package main

import (
	"context"
	"strings"

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

	PathToHandlerMap = map[string]func(ctx context.Context, params handlers.HandlerParams) (events.APIGatewayProxyResponse, error){
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
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	appCfg.Logger.Info("recieved event", zap.Any("Request", req))

	params := handlers.HandlerParams{
		AppCfg:       appCfg,
		Request:      req,
		UserRepo:     userRepo,
		ReceiverRepo: receiverRepo,
	}

	if handler, ok := PathToHandlerMap[removePathPrefix(req.RequestContext.Path)]; ok {
		return handler(ctx, params)
	}

	appCfg.Logger.Error("unsupported request path", zap.Any(log.PathLogKey, req.RequestContext.Path))
	return response.CreateBadRequestResponse(), nil
}

func removePathPrefix(path string) string {
	pathPrefixes := []string{"/Stage", "/Prod"}
	for _, prefix := range pathPrefixes {
		path = strings.TrimPrefix(path, prefix)
	}
	return path
}

func main() {
	lambda.Start(handler)
}
