package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
	"github.com/care-giver-app/care-giver-api/internal/response"
	"github.com/stretchr/testify/assert"
)

func makeTemplateParams(req events.APIGatewayProxyRequest) HandlerParams {
	return HandlerParams{
		AppCfg:  appconfig.NewAppConfig(),
		Request: req,
	}
}

func TestHandleListTrackerTemplates(t *testing.T) {
	tests := map[string]struct {
		request          events.APIGatewayProxyRequest
		expectedStatus   int
		checkBody        bool
		expectedMinCount int
	}{
		"Happy Path - Returns All 7 Templates": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedStatus:   http.StatusOK,
			checkBody:        true,
			expectedMinCount: 7,
		},
		"Sad Path - Missing JWT Claim": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: map[string]interface{}{},
				},
			},
			expectedStatus: http.StatusForbidden,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			params := makeTemplateParams(tc.request)
			resp, err := HandleListTrackerTemplates(context.Background(), params)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)

			if tc.checkBody {
				var body []interface{}
				assert.Nil(t, json.Unmarshal([]byte(resp.Body), &body))
				assert.GreaterOrEqual(t, len(body), tc.expectedMinCount)
			}
		})
	}
}

func TestHandleGetTrackerTemplate(t *testing.T) {
	tests := map[string]struct {
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
		checkName        string
	}{
		"Happy Path - Weight Template Found": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:     http.MethodGet,
				PathParameters: map[string]string{"templateName": "Weight"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: events.APIGatewayProxyResponse{StatusCode: http.StatusOK},
			checkName:        "Weight",
		},
		"Happy Path - URL Encoded Name": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:     http.MethodGet,
				PathParameters: map[string]string{"templateName": "Doctor%20Appointment"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: events.APIGatewayProxyResponse{StatusCode: http.StatusOK},
			checkName:        "Doctor Appointment",
		},
		"Sad Path - Missing JWT Claim": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:     http.MethodGet,
				PathParameters: map[string]string{"templateName": "Weight"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: map[string]interface{}{},
				},
			},
			expectedResponse: response.CreateAccessDeniedResponse(),
		},
		"Sad Path - Missing templateName": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:     http.MethodGet,
				PathParameters: map[string]string{},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Unknown Template Name": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:     http.MethodGet,
				PathParameters: map[string]string{"templateName": "NonExistentTemplate"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: response.CreateResourceNotFoundResponse(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			params := makeTemplateParams(tc.request)
			resp, err := HandleGetTrackerTemplate(context.Background(), params)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResponse.StatusCode, resp.StatusCode)

			if tc.checkName != "" && resp.StatusCode == http.StatusOK {
				var body map[string]interface{}
				assert.Nil(t, json.Unmarshal([]byte(resp.Body), &body))
				assert.Equal(t, tc.checkName, body["name"])
			}
		})
	}
}
