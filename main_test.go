package main

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/handlers"
	"github.com/care-giver-app/care-giver-api/internal/response"
	"github.com/stretchr/testify/assert"
)

type MockRegistry struct{}

func (m *MockRegistry) GetHandler(request events.APIGatewayProxyRequest) (func(ctx context.Context, params handlers.HandlerParams) (events.APIGatewayProxyResponse, error), bool) {
	fmt.Printf("Received request: %v\n", request)
	if request.RequestContext.ResourcePath == "/good/path" {
		return func(ctx context.Context, params handlers.HandlerParams) (events.APIGatewayProxyResponse, error) {
			return events.APIGatewayProxyResponse{
				Body: "Handler One",
			}, nil
		}, true
	}
	return nil, false
}

func (m *MockRegistry) RunHandler(ctx context.Context, handler func(ctx context.Context, params handlers.HandlerParams) (events.APIGatewayProxyResponse, error), request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if handler != nil {
		params := handlers.HandlerParams{
			Request: request,
		}
		return handler(ctx, params)
	}
	return events.APIGatewayProxyResponse{}, errors.New("handler not supported")
}

func TestHandler(t *testing.T) {
	tests := map[string]struct {
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
	}{
		"Happy Path - Handler One": {
			request: events.APIGatewayProxyRequest{
				RequestContext: events.APIGatewayProxyRequestContext{
					ResourcePath: "/good/path",
				},
				HTTPMethod: "POST",
			},
			expectedResponse: events.APIGatewayProxyResponse{
				Body: "Handler One",
			},
		},
		"Sad Path - Invalid Path": {
			request: events.APIGatewayProxyRequest{
				RequestContext: events.APIGatewayProxyRequestContext{
					ResourcePath: "/bad/path",
				},
				HTTPMethod: "GET",
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
	}

	handlerRegistry = &MockRegistry{}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp, err := handler(context.Background(), tc.request)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResponse, resp)
		})
	}
}
