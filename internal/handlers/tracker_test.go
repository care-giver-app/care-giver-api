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

func makeTrackerParams(req events.APIGatewayProxyRequest) HandlerParams {
	return HandlerParams{
		AppCfg:           appconfig.NewAppConfig(),
		Request:          req,
		UserRepo:         testUserRepo,
		ReceiverRepo:     testReceiverRepo,
		RelationshipRepo: testRelationshipRepo,
		TrackerRepo:      testTrackerRepo,
	}
}

func authedRequest(uid string) map[string]interface{} {
	return map[string]interface{}{
		"custom:db_user_id": uid,
	}
}

func TestHandleCreateTracker(t *testing.T) {
	validBody, _ := json.Marshal(map[string]interface{}{
		"receiverId": "Receiver#123",
		"name":       "Walk",
		"kind":       "event",
		"fields":     []interface{}{},
		"icon":       "assets/icon.svg",
		"color":      map[string]string{"primary": "#000", "secondary": "#fff"},
	})

	tests := map[string]struct {
		request        events.APIGatewayProxyRequest
		expectedStatus int
	}{
		"Happy Path - Tracker Created": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       string(validBody),
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedStatus: http.StatusOK,
		},
		"Sad Path - Missing JWT Claim": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       string(validBody),
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: map[string]interface{}{},
				},
			},
			expectedStatus: http.StatusForbidden,
		},
		"Sad Path - Bad Body": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       `{"receiverId": false}`,
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		"Sad Path - Not A CareGiver": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       string(validBody),
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#NotACareGiver"),
				},
			},
			expectedStatus: http.StatusForbidden,
		},
		"Sad Path - Name Conflict": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body: func() string {
					b, _ := json.Marshal(map[string]interface{}{
						"receiverId": "Receiver#123",
						"name":       "Duplicate Name",
						"kind":       "event",
						"fields":     []interface{}{},
						"icon":       "assets/icon.svg",
						"color":      map[string]string{"primary": "#000", "secondary": "#fff"},
					})
					return string(b)
				}(),
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedStatus: http.StatusConflict,
		},
		"Sad Path - Repo Error": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body: func() string {
					b, _ := json.Marshal(map[string]interface{}{
						"receiverId": "Receiver#Error",
						"name":       "Some Tracker",
						"kind":       "event",
						"fields":     []interface{}{},
						"icon":       "assets/icon.svg",
						"color":      map[string]string{"primary": "#000", "secondary": "#fff"},
					})
					return string(b)
				}(),
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedStatus: http.StatusInternalServerError,
		},
		"Sad Path - Invalid Kind": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body: func() string {
					b, _ := json.Marshal(map[string]interface{}{
						"receiverId": "Receiver#123",
						"name":       "Bad Tracker",
						"kind":       "invalid_kind",
						"fields":     []interface{}{},
						"icon":       "assets/icon.svg",
						"color":      map[string]string{"primary": "#000", "secondary": "#fff"},
					})
					return string(b)
				}(),
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			params := makeTrackerParams(tc.request)
			resp, err := HandleCreateTracker(context.Background(), params)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
		})
	}
}

func TestHandleListTrackers(t *testing.T) {
	tests := map[string]struct {
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
	}{
		"Happy Path - Trackers Listed": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
				PathParameters: map[string]string{
					"receiverId": "Receiver#123",
				},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: events.APIGatewayProxyResponse{StatusCode: http.StatusOK},
		},
		"Happy Path - Empty List Returns Array": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
				PathParameters: map[string]string{
					"receiverId": "Receiver#Empty",
				},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: events.APIGatewayProxyResponse{StatusCode: http.StatusOK},
		},
		"Sad Path - Missing JWT Claim": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:     http.MethodGet,
				PathParameters: map[string]string{"receiverId": "Receiver#123"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: map[string]interface{}{},
				},
			},
			expectedResponse: response.CreateAccessDeniedResponse(),
		},
		"Sad Path - Bad Path Parameter": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:     http.MethodGet,
				PathParameters: map[string]string{"wrongParam": "Receiver#123"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Not A CareGiver": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:     http.MethodGet,
				PathParameters: map[string]string{"receiverId": "Receiver#123"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#NotACareGiver"),
				},
			},
			expectedResponse: response.CreateAccessDeniedResponse(),
		},
		"Sad Path - Repo Error": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:     http.MethodGet,
				PathParameters: map[string]string{"receiverId": "Receiver#Error"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: response.CreateInternalServerErrorResponse(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			params := makeTrackerParams(tc.request)
			resp, err := HandleListTrackers(context.Background(), params)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResponse.StatusCode, resp.StatusCode)

			if tc.expectedResponse.StatusCode == http.StatusOK {
				var body interface{}
				assert.Nil(t, json.Unmarshal([]byte(resp.Body), &body))
				assert.IsType(t, []interface{}{}, body, "list response must be a JSON array")
			}
		})
	}
}

func TestHandleGetTracker(t *testing.T) {
	tests := map[string]struct {
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
	}{
		"Happy Path - Tracker Found": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:     http.MethodGet,
				PathParameters: map[string]string{"trackerId": "Tracker#123"},
				QueryStringParameters: map[string]string{
					"receiverId": "Receiver#123",
				},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: events.APIGatewayProxyResponse{StatusCode: http.StatusOK},
		},
		"Sad Path - Missing JWT Claim Returns 403 Not 404": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodGet,
				PathParameters:        map[string]string{"trackerId": "Tracker#123"},
				QueryStringParameters: map[string]string{"receiverId": "Receiver#123"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: map[string]interface{}{},
				},
			},
			expectedResponse: response.CreateAccessDeniedResponse(),
		},
		"Sad Path - Not A CareGiver Returns 403 Not 404": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodGet,
				PathParameters:        map[string]string{"trackerId": "Tracker#123"},
				QueryStringParameters: map[string]string{"receiverId": "Receiver#123"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#NotACareGiver"),
				},
			},
			expectedResponse: response.CreateAccessDeniedResponse(),
		},
		"Sad Path - Tracker Not Found": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodGet,
				PathParameters:        map[string]string{"trackerId": "Tracker#NotFound"},
				QueryStringParameters: map[string]string{"receiverId": "Receiver#123"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: response.CreateResourceNotFoundResponse(),
		},
		"Sad Path - Bad Path Parameter": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodGet,
				PathParameters:        map[string]string{"wrong": "Tracker#123"},
				QueryStringParameters: map[string]string{"receiverId": "Receiver#123"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Missing receiverId Query Param": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodGet,
				PathParameters:        map[string]string{"trackerId": "Tracker#123"},
				QueryStringParameters: map[string]string{},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Repo Error": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodGet,
				PathParameters:        map[string]string{"trackerId": "Tracker#Error"},
				QueryStringParameters: map[string]string{"receiverId": "Receiver#123"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: response.CreateInternalServerErrorResponse(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			params := makeTrackerParams(tc.request)
			resp, err := HandleGetTracker(context.Background(), params)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResponse.StatusCode, resp.StatusCode)
		})
	}
}

func TestHandleUpdateTracker(t *testing.T) {
	validUpdateBody, _ := json.Marshal(map[string]interface{}{
		"name": "Updated Name",
	})

	tests := map[string]struct {
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
	}{
		"Happy Path - Tracker Updated": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodPut,
				PathParameters:        map[string]string{"trackerId": "Tracker#123"},
				QueryStringParameters: map[string]string{"receiverId": "Receiver#123"},
				Body:                  string(validUpdateBody),
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: events.APIGatewayProxyResponse{StatusCode: http.StatusOK},
		},
		"Sad Path - Missing JWT Claim": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodPut,
				PathParameters:        map[string]string{"trackerId": "Tracker#123"},
				QueryStringParameters: map[string]string{"receiverId": "Receiver#123"},
				Body:                  string(validUpdateBody),
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: map[string]interface{}{},
				},
			},
			expectedResponse: response.CreateAccessDeniedResponse(),
		},
		"Sad Path - Tracker Not Found": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodPut,
				PathParameters:        map[string]string{"trackerId": "Tracker#NotFound"},
				QueryStringParameters: map[string]string{"receiverId": "Receiver#123"},
				Body:                  string(validUpdateBody),
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: response.CreateResourceNotFoundResponse(),
		},
		"Sad Path - Immutable Field kind In Body": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodPut,
				PathParameters:        map[string]string{"trackerId": "Tracker#123"},
				QueryStringParameters: map[string]string{"receiverId": "Receiver#123"},
				Body:                  `{"kind": "event"}`,
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Not A CareGiver": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodPut,
				PathParameters:        map[string]string{"trackerId": "Tracker#123"},
				QueryStringParameters: map[string]string{"receiverId": "Receiver#123"},
				Body:                  string(validUpdateBody),
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#NotACareGiver"),
				},
			},
			expectedResponse: response.CreateAccessDeniedResponse(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			params := makeTrackerParams(tc.request)
			resp, err := HandleUpdateTracker(context.Background(), params)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResponse.StatusCode, resp.StatusCode)
		})
	}
}

func TestHandleDeleteTracker(t *testing.T) {
	tests := map[string]struct {
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
	}{
		"Happy Path - Tracker Deleted": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodDelete,
				PathParameters:        map[string]string{"trackerId": "Tracker#123"},
				QueryStringParameters: map[string]string{"receiverId": "Receiver#123"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: response.FormatResponse(map[string]string{"status": response.Success}, http.StatusOK),
		},
		"Sad Path - Missing JWT Claim": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodDelete,
				PathParameters:        map[string]string{"trackerId": "Tracker#123"},
				QueryStringParameters: map[string]string{"receiverId": "Receiver#123"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: map[string]interface{}{},
				},
			},
			expectedResponse: response.CreateAccessDeniedResponse(),
		},
		"Sad Path - Not A CareGiver": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodDelete,
				PathParameters:        map[string]string{"trackerId": "Tracker#123"},
				QueryStringParameters: map[string]string{"receiverId": "Receiver#123"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#NotACareGiver"),
				},
			},
			expectedResponse: response.CreateAccessDeniedResponse(),
		},
		"Sad Path - Tracker Not Found": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodDelete,
				PathParameters:        map[string]string{"trackerId": "Tracker#NotFound"},
				QueryStringParameters: map[string]string{"receiverId": "Receiver#123"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: response.CreateResourceNotFoundResponse(),
		},
		"Sad Path - Missing receiverId Query Param": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodDelete,
				PathParameters:        map[string]string{"trackerId": "Tracker#123"},
				QueryStringParameters: map[string]string{},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Repo Delete Error": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod:            http.MethodDelete,
				PathParameters:        map[string]string{"trackerId": "Tracker#Error"},
				QueryStringParameters: map[string]string{"receiverId": "Receiver#123"},
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: authedRequest("User#123"),
				},
			},
			expectedResponse: response.CreateInternalServerErrorResponse(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			params := makeTrackerParams(tc.request)
			resp, err := HandleDeleteTracker(context.Background(), params)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResponse.StatusCode, resp.StatusCode)
		})
	}
}
