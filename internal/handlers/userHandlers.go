package handlers

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
)

func HandleReceiver(ctx context.Context, appCfg *appconfig.AppConfig, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{}, nil
}

func HandleReceiverEvent(ctx context.Context, appCfg *appconfig.AppConfig, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{}, nil
}
