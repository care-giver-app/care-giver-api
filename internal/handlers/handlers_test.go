package handlers

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func TestGetRegisteredHandler(t *testing.T) {
	handlerOne := func(ctx context.Context, params HandlerParams) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{
			Body: "Handler One",
		}, nil
	}

	handlerTwo := func(ctx context.Context, params HandlerParams) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{
			Body: "Handler Two",
		}, nil
	}

	handlerThree := func(ctx context.Context, params HandlerParams) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{
			Body: "Handler Three",
		}, nil
	}

	handlersMap = map[Endpoint]HandlerFunc{
		{Path: "/testPathOne", Method: "POST"}:   handlerOne,
		{Path: "/test/path/two", Method: "GET"}:  handlerTwo,
		{Path: "/test/path/two", Method: "POST"}: handlerThree,
	}

	tests := map[string]struct {
		request         events.APIGatewayProxyRequest
		expectedHandler func(ctx context.Context, params HandlerParams) (events.APIGatewayProxyResponse, error)
	}{
		"Happy Path - Handler One": {
			request: events.APIGatewayProxyRequest{
				RequestContext: events.APIGatewayProxyRequestContext{
					ResourcePath: "/testPathOne",
				},
				HTTPMethod: "POST",
			},
			expectedHandler: handlerOne,
		},
		"Happy Path - Handler Two": {
			request: events.APIGatewayProxyRequest{
				RequestContext: events.APIGatewayProxyRequestContext{
					ResourcePath: "/test/path/two",
				},
				HTTPMethod: "GET",
			},
			expectedHandler: handlerTwo,
		},
		"Happy Path - Handler Two With Prefix": {
			request: events.APIGatewayProxyRequest{
				RequestContext: events.APIGatewayProxyRequestContext{
					ResourcePath: "/Stage/test/path/two",
				},
				HTTPMethod: "GET",
			},
			expectedHandler: handlerTwo,
		},
		"Happy Path - Handler Three With Same Path But Different Method": {
			request: events.APIGatewayProxyRequest{
				RequestContext: events.APIGatewayProxyRequestContext{
					ResourcePath: "/Stage/test/path/two",
				},
				HTTPMethod: "POST",
			},
			expectedHandler: handlerThree,
		},
		"Sad Path - Wrong Method": {
			request: events.APIGatewayProxyRequest{
				RequestContext: events.APIGatewayProxyRequestContext{
					ResourcePath: "/testPathOne",
				},
				HTTPMethod: "GET",
			},
		},
		"Sad Path - Invalid Path": {
			request: events.APIGatewayProxyRequest{
				RequestContext: events.APIGatewayProxyRequestContext{
					ResourcePath: "/bad/path",
				},
			},
		},
	}

	testHandlerRegistry := NewRegistry(nil, nil, nil, nil)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			handler, ok := testHandlerRegistry.GetHandler(tc.request)
			if tc.expectedHandler != nil {
				assert.True(t, ok)

				expectedResp, expectedErr := tc.expectedHandler(context.Background(), HandlerParams{})
				actualResp, actualErr := handler(context.Background(), HandlerParams{})
				assert.Equal(t, expectedResp, actualResp)
				assert.Equal(t, expectedErr, actualErr)
			} else {
				assert.False(t, ok)
			}
		})
	}
}

func TestRunHandler(t *testing.T) {
	testHandler := func(ctx context.Context, params HandlerParams) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{
			Body: "Test Response",
		}, nil
	}

	testRegistry := &Registry{
		AppCfg:       nil,
		UserRepo:     nil,
		ReceiverRepo: nil,
	}

	request := events.APIGatewayProxyRequest{
		RequestContext: events.APIGatewayProxyRequestContext{
			ResourcePath: "/test/path",
		},
		HTTPMethod: "POST",
	}

	response, err := testRegistry.RunHandler(context.Background(), testHandler, request)
	assert.Nil(t, err)
	assert.Equal(t, "Test Response", response.Body)
}
