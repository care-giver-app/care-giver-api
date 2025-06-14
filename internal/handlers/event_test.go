package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
	"github.com/care-giver-app/care-giver-api/internal/event"
	"github.com/care-giver-app/care-giver-api/internal/response"
	"github.com/stretchr/testify/assert"
)

func TestHandleReceiverEvent(t *testing.T) {
	tests := map[string]struct {
		requestMethod        string
		requestBody          map[string]interface{}
		expectedResponseBody map[string]interface{}
		expectedStatusCode   int
	}{
		"Happy Path - Event Added": {
			requestMethod: http.MethodPost,
			requestBody: map[string]interface{}{
				"receiverId": "Receiver#123",
				"userID":     "User#123",
				"type":       "Shower",
			},
			expectedResponseBody: map[string]interface{}{
				"status":     "Success",
				"receiverId": "Receiver#123",
			},
			expectedStatusCode: http.StatusOK,
		},
		"Happy Path - Event Added - With Timestamp": {
			requestMethod: http.MethodPost,
			requestBody: map[string]interface{}{
				"receiverId": "Receiver#123",
				"userID":     "User#123",
				"type":       "Shower",
				"timestamp":  "2023-10-01T12:00:00Z",
			},
			expectedResponseBody: map[string]interface{}{
				"status":     "Success",
				"receiverId": "Receiver#123",
			},
			expectedStatusCode: http.StatusOK,
		},
		"Happy Path - Event Added - With Data": {
			requestMethod: http.MethodPost,
			requestBody: map[string]interface{}{
				"receiverId": "Receiver#123",
				"userID":     "User#123",
				"type":       "Shower",
				"data": []event.DataPoint{
					{
						Name:  "some name",
						Value: "some value",
					},
				},
			},
			expectedResponseBody: map[string]interface{}{
				"status":     "Success",
				"receiverId": "Receiver#123",
			},
			expectedStatusCode: http.StatusOK,
		},
		"Sad Path - Bad Body": {
			requestMethod: http.MethodPost,
			requestBody: map[string]interface{}{
				"receiverId": false,
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		"Sad Path - Error Getting User": {
			requestMethod: http.MethodPost,
			requestBody: map[string]interface{}{
				"receiverId": "Receiver#123",
				"userId":     "User#Error",
				"type":       "Shower",
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		"Sad Path - User Is Not A Care Giver": {
			requestMethod: http.MethodPost,
			requestBody: map[string]interface{}{
				"receiverId": "Receiver#123",
				"userId":     "User#NotACareGiver",
				"type":       "Shower",
			},
			expectedStatusCode: http.StatusForbidden,
		},
		"Sad Path - Bad Event Name": {
			requestMethod: http.MethodPost,
			requestBody: map[string]interface{}{
				"receiverId": "Receiver#123",
				"userId":     "User#123",
				"type":       "badEventType",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		"Sad Path - Bad Event Data": {
			requestMethod: http.MethodPost,
			requestBody: map[string]interface{}{
				"receiverId": "Receiver#123",
				"userId":     "User#123",
				"type":       "Weight",
				"data":       "wrongDataType",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		"Sad Path - Error Adding Event": {
			requestMethod: http.MethodPost,
			requestBody: map[string]interface{}{
				"receiverId": "Receiver#Error",
				"userId":     "User#123",
				"type":       "Weight",
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			requestBody, err := json.Marshal(tc.requestBody)
			assert.Nil(t, err)
			req := events.APIGatewayProxyRequest{
				HTTPMethod: tc.requestMethod,
				Body:       string(requestBody),
			}

			params := HandlerParams{
				AppCfg:       appconfig.NewAppConfig(),
				Request:      req,
				UserRepo:     testUserRepo,
				ReceiverRepo: testReceiverRepo,
				EventRepo:    testEventRepo,
			}
			resp, err := HandleReceiverEvent(context.Background(), params)
			assert.Nil(t, err)
			assert.NotNil(t, resp)

			assert.Equal(t, tc.expectedStatusCode, resp.StatusCode)

			if tc.expectedStatusCode == http.StatusOK {
				var responseBody map[string]interface{}
				err = json.Unmarshal([]byte(resp.Body), &responseBody)
				assert.Nil(t, err)
				tc.expectedResponseBody["eventId"] = responseBody["eventId"]
				assert.Equal(t, tc.expectedResponseBody, responseBody)
			}
		})
	}
}

func TestHandleDeleteReceiverEvent(t *testing.T) {
	tests := map[string]struct {
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
	}{
		"Happy Path - Event Deleted Successfully": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodDelete,
				PathParameters: map[string]string{
					"eventId": "Event#123",
				},
				QueryStringParameters: map[string]string{
					"userId":     "User#123",
					"receiverId": "Receiver#123",
				},
			},
			expectedResponse: response.FormatResponse(
				map[string]string{
					"status": response.Success,
				}, http.StatusOK,
			),
		},
		"Sad Path - Bad Path Parameter - eventId": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodDelete,
				PathParameters: map[string]string{
					"NotEventId": "Event#123",
				},
				QueryStringParameters: map[string]string{
					"userId":     "User#123",
					"receiverId": "Receiver#123",
				},
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Bad Query Parameter - userId": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodDelete,
				PathParameters: map[string]string{
					"eventId": "Event#123",
				},
				QueryStringParameters: map[string]string{
					"receiverId": "Receiver#123",
				},
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Bad Query Parameter - receiverId": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodDelete,
				PathParameters: map[string]string{
					"eventId": "Event#123",
				},
				QueryStringParameters: map[string]string{
					"userId": "User#123",
				},
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Error Retrieving User": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodDelete,
				PathParameters: map[string]string{
					"eventId": "Event#123",
				},
				QueryStringParameters: map[string]string{
					"userId":     "User#Error",
					"receiverId": "Receiver#123",
				},
			},
			expectedResponse: response.CreateInternalServerErrorResponse(),
		},
		"Sad Path - User Is Not A Care Giver": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodDelete,
				PathParameters: map[string]string{
					"eventId": "Event#123",
				},
				QueryStringParameters: map[string]string{
					"userId":     "User#NotACareGiver",
					"receiverId": "Receiver#123",
				},
			},
			expectedResponse: response.CreateAccessDeniedResponse(),
		},
		"Sad Path - Error Adding Event": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodDelete,
				PathParameters: map[string]string{
					"eventId": "Event#123",
				},
				QueryStringParameters: map[string]string{
					"userId":     "User#123",
					"receiverId": "Receiver#Error",
				},
			},
			expectedResponse: response.CreateInternalServerErrorResponse(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			params := HandlerParams{
				AppCfg:       appconfig.NewAppConfig(),
				Request:      tc.request,
				UserRepo:     testUserRepo,
				ReceiverRepo: testReceiverRepo,
				EventRepo:    testEventRepo,
			}

			resp, err := HandleDeleteReceiverEvent(context.Background(), params)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResponse, resp)
		})
	}
}

func TestHandleGetReceiverEvents(t *testing.T) {
	tests := map[string]struct {
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
	}{
		"Happy Path - Events Retrieved Successfully": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodDelete,
				PathParameters: map[string]string{
					"receiverId": "Receiver#123",
				},
				QueryStringParameters: map[string]string{
					"userId": "User#123",
				},
			},
			expectedResponse: response.FormatResponse(
				[]event.Entry{
					{
						EventID:    "Event#123",
						ReceiverID: "Receiver#123",
					},
				}, http.StatusOK,
			),
		},
		"Sad Path - Bad Path Parameter - receiverId": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodDelete,
				PathParameters: map[string]string{
					"notReceiverId": "Receiver#123",
				},
				QueryStringParameters: map[string]string{
					"userId": "User#123",
				},
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Bad Query Parameter - userId": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodDelete,
				PathParameters: map[string]string{
					"receiverId": "Receiver#123",
				},
				QueryStringParameters: map[string]string{},
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Error Retrieving User": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodDelete,
				PathParameters: map[string]string{
					"receiverId": "Receiver#123",
				},
				QueryStringParameters: map[string]string{
					"userId": "User#Error",
				},
			},
			expectedResponse: response.CreateInternalServerErrorResponse(),
		},
		"Sad Path - User Is Not A Care Giver": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodDelete,
				PathParameters: map[string]string{
					"receiverId": "Receiver#123",
				},
				QueryStringParameters: map[string]string{
					"userId": "User#NotACareGiver",
				},
			},
			expectedResponse: response.CreateAccessDeniedResponse(),
		},
		"Sad Path - Error Getting Event": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodDelete,
				PathParameters: map[string]string{
					"receiverId": "Receiver#Error",
				},
				QueryStringParameters: map[string]string{
					"userId": "User#123",
				},
			},
			expectedResponse: response.CreateInternalServerErrorResponse(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			params := HandlerParams{
				AppCfg:       appconfig.NewAppConfig(),
				Request:      tc.request,
				UserRepo:     testUserRepo,
				ReceiverRepo: testReceiverRepo,
				EventRepo:    testEventRepo,
			}

			resp, err := HandleGetReceiverEvents(context.Background(), params)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResponse, resp)
		})
	}
}
