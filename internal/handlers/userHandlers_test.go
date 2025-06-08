package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
	"github.com/care-giver-app/care-giver-api/internal/response"
	"github.com/care-giver-app/care-giver-api/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestHandleCreateUser(t *testing.T) {
	tests := map[string]struct {
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
		isSuccess        bool
	}{
		"Happy Path - User Added": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "{\"email\": \"good@test.com\", \"firstName\":\"Demo\", \"lastName\":\"Daniel\", \"password\":\"myPass\"}",
			},
			isSuccess: true,
		},
		"Sad Path - Wrong Method": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "BadMethod",
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Bad Request Body": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "{\"email\": false}",
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Error Add User To DB": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "{\"email\": \"error@test.com\", \"firstName\":\"Demo\", \"lastName\":\"Daniel\", \"password\":\"myPass\"}",
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
			}
			resp, err := HandleCreateUser(context.Background(), params)

			expectedResponse := tc.expectedResponse
			if tc.isSuccess {
				var respStruct CreateUserResponse
				testErr := json.Unmarshal([]byte(resp.Body), &respStruct)
				assert.Nil(t, testErr)

				expectedResp := CreateUserResponse{
					UserID: respStruct.UserID,
					Status: response.Success,
				}

				expectedResponse = response.FormatResponse(expectedResp, http.StatusOK)
			}

			assert.Nil(t, err)
			assert.Equal(t, expectedResponse, resp)
		})
	}
}

func TestHandleGetUser(t *testing.T) {
	tests := map[string]struct {
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
	}{
		"Happy Path - User Retrieved": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
				PathParameters: map[string]string{
					"userId": "User#123",
				},
			},
			expectedResponse: response.FormatResponse(user.User{
				PrimaryCareReceivers: []string{"Receiver#123", "Receiver#123Error"},
			}, http.StatusOK),
		},
		"Sad Path - Wrong Method": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "BadMethod",
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Bad Path Parameters": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
				PathParameters: map[string]string{
					"userId": "BadValue",
				},
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Error Getting User From DB": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodGet,
				PathParameters: map[string]string{
					"userId": "User#Error",
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
			}
			resp, err := HandleGetUser(context.Background(), params)

			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResponse, resp)
		})
	}
}

func TestHandleUserPrimaryReceiver(t *testing.T) {
	tests := map[string]struct {
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
		isSuccess        bool
	}{
		"Happy Path - User Added": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "{\"userId\": \"User#123\", \"firstName\":\"Good\", \"lastName\":\"Daniel\"}",
			},
			isSuccess: true,
		},
		"Sad Path - Wrong Method": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "BadMethod",
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Bad Request Body": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "{\"userId\": false}",
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Error Creating Receiver": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "{\"userId\": \"User#123\", \"firstName\":\"Error\", \"lastName\":\"Daniel\"}",
			},
			expectedResponse: response.CreateInternalServerErrorResponse(),
		},
		"Sad Path - Error Updating Receiver List": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "{\"userId\": \"User#ListError\", \"firstName\":\"Good\", \"lastName\":\"Daniel\"}",
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
			}
			resp, err := HandleUserPrimaryReceiver(context.Background(), params)

			expectedResponse := tc.expectedResponse
			if tc.isSuccess {
				var respStruct PrimaryReceiverResponse
				testErr := json.Unmarshal([]byte(resp.Body), &respStruct)
				assert.Nil(t, testErr)

				expectedResp := PrimaryReceiverResponse{
					ReceiverID: respStruct.ReceiverID,
					Status:     response.Success,
				}

				expectedResponse = response.FormatResponse(expectedResp, http.StatusOK)
			}

			assert.Nil(t, err)
			assert.Equal(t, expectedResponse, resp)
		})
	}
}

func TestHandleUserAdditionalReceiver(t *testing.T) {
	tests := map[string]struct {
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
	}{
		"Happy Path - User Added": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "{\"userId\": \"User#123\", \"receiverId\":\"Receiver#789\"}",
			},
			expectedResponse: response.FormatResponse(map[string]string{
				"status": "success",
			}, http.StatusOK),
		},
		"Sad Path - Wrong Method": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: "BadMethod",
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Bad Request Body": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "{\"userId\": false}",
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Error Updating Receiver List": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "{\"userId\": \"User#ListError\", \"receiverId\":\"Receiver#789\"}",
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
			}
			resp, err := HandleUserAdditionalReceiver(context.Background(), params)

			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResponse, resp)
		})
	}
}
