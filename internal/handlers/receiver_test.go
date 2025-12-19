package handlers

import (
	"context"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
	"github.com/care-giver-app/care-giver-api/internal/response"
	"github.com/care-giver-app/care-giver-golang-common/pkg/receiver"
	"github.com/stretchr/testify/assert"
)

func TestHandleReceiver(t *testing.T) {
	tests := map[string]struct {
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
	}{
		"Happy Path - Got Receiver": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
				PathParameters: map[string]string{
					"receiverId": "Receiver#123",
				},
				QueryStringParameters: map[string]string{
					"userId": "User#123",
				},
			},
			expectedResponse: response.FormatResponse(receiver.Receiver{
				FirstName: "Success",
			}, http.StatusOK),
		},
		"Sad Path - Bad Path Parameters": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
				PathParameters: map[string]string{
					"receiverId": "BadValue",
				},
				QueryStringParameters: map[string]string{
					"userId": "User#123",
				},
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Bad Query Parameters": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
				PathParameters: map[string]string{
					"receiverId": "Receiver#123",
				},
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Error Getting Receiver From DB": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
				PathParameters: map[string]string{
					"receiverId": "Receiver#Error",
				},
				QueryStringParameters: map[string]string{
					"userId": "User#123",
				},
			},
			expectedResponse: response.CreateInternalServerErrorResponse(),
		},
		"Sad Path - Error Getting User From DB": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
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
				HTTPMethod: http.MethodGet,
				PathParameters: map[string]string{
					"receiverId": "Receiver#123",
				},
				QueryStringParameters: map[string]string{
					"userId": "User#NotACareGiver",
				},
			},
			expectedResponse: response.CreateAccessDeniedResponse(),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			params := HandlerParams{
				AppCfg:           appconfig.NewAppConfig(),
				Request:          tc.request,
				UserRepo:         testUserRepo,
				ReceiverRepo:     testReceiverRepo,
				RelationshipRepo: testRelationshipRepo,
			}
			resp, err := HandleReceiver(context.Background(), params)

			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResponse, resp)
		})
	}
}
