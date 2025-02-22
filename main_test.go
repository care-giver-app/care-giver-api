package main

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
	"github.com/care-giver-app/care-giver-api/internal/response"
	"github.com/stretchr/testify/assert"
)

func testHandlerOne(ctx context.Context, appCfg *appconfig.AppConfig, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body: "Handler One",
	}, nil
}

func testHandlerTwo(ctx context.Context, appCfg *appconfig.AppConfig, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body: "Handler Two",
	}, nil
}

func TestHandler(t *testing.T) {
	PathToHandlerMap = map[string]func(ctx context.Context, appCfg *appconfig.AppConfig, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error){
		"/testPathOne":   testHandlerOne,
		"/test/path/two": testHandlerTwo,
	}

	tests := map[string]struct {
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
	}{
		"Happy Path - Handler One": {
			request: events.APIGatewayProxyRequest{
				RequestContext: events.APIGatewayProxyRequestContext{
					Path: "/testPathOne",
				},
			},
			expectedResponse: events.APIGatewayProxyResponse{
				Body: "Handler One",
			},
		},
		"Happy Path - Handler Two": {
			request: events.APIGatewayProxyRequest{
				RequestContext: events.APIGatewayProxyRequestContext{
					Path: "/test/path/two",
				},
			},
			expectedResponse: events.APIGatewayProxyResponse{
				Body: "Handler Two",
			},
		},
		"Sad Path - Invalid Path": {
			request: events.APIGatewayProxyRequest{
				RequestContext: events.APIGatewayProxyRequestContext{
					Path: "/bad/path",
				},
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := handler(context.Background(), tc.request)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResponse, resp)
		})
	}
}
