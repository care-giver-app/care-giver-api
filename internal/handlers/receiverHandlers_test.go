package handlers

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
	"github.com/care-giver-app/care-giver-api/internal/receiver"
	receiverEvents "github.com/care-giver-app/care-giver-api/internal/receiver/events"
	"github.com/care-giver-app/care-giver-api/internal/response"
	"github.com/care-giver-app/care-giver-api/internal/user"
	"github.com/stretchr/testify/assert"
)

type MockUserRepo struct{}

func (mu *MockUserRepo) CreateUser(u user.User) error {
	switch u.Email {
	case "good@test.com":
		return nil
	case "error@test.com":
		return errors.New("error creating user")
	}
	return errors.New("unsupported mock")
}

func (mu *MockUserRepo) GetUser(uid string) (user.User, error) {
	switch uid {
	case "User#123":
		return user.User{
			PrimaryCareReceivers: []string{"Receiver#123", "Receiver#123Error"},
		}, nil
	case "User#NotACareGiver":
		return user.User{}, nil
	case "User#Error":
		return user.User{}, errors.New("error getting user from db")
	}
	return user.User{}, errors.New("unsupported mock")
}

func (mu *MockUserRepo) UpdateReceiverList(uid string, rid string, listName string) error {
	switch uid {
	case "User#123":
		return nil
	case "User#ListError":
		return errors.New("error updating receiver list")
	}
	return errors.New("unsupported mock")
}

type MockReceiverRepo struct{}

func (mr *MockReceiverRepo) CreateReceiver(r receiver.Receiver) error {
	switch r.FirstName {
	case "Good":
		return nil
	case "Error":
		return errors.New("error creating receiver")
	}
	return errors.New("unsupported mock")
}

func (mr *MockReceiverRepo) GetReceiver(rid string) (receiver.Receiver, error) {
	switch rid {
	case "Receiver#123":
		return receiver.Receiver{
			FirstName: "Success",
		}, nil
	case "Receiver#Error":
		return receiver.Receiver{}, errors.New("error retrieving from db")
	}
	return receiver.Receiver{}, errors.New("unsupported mock")
}

func (mr *MockReceiverRepo) UpdateReceiver(rid string, evt receiverEvents.Event, eventName string) error {
	switch rid {
	case "Receiver#123":
		return nil
	case "Receiver#123Error":
		return errors.New("error updating receiver in db")
	}
	return errors.New("unsupported mock")
}

var (
	testUserRepo     = &MockUserRepo{}
	testReceiverRepo = &MockReceiverRepo{}
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
			},
			expectedResponse: response.FormatResponse(receiver.Receiver{
				FirstName: "Success",
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
					"receiverId": "BadValue",
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
			resp, err := HandleReceiver(context.Background(), params)

			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResponse, resp)
		})
	}
}

func TestHandleReceiverEvent(t *testing.T) {
	tests := map[string]struct {
		request          events.APIGatewayProxyRequest
		expectedResponse events.APIGatewayProxyResponse
	}{
		"Happy Path - Event Added": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "{\"receiverId\": \"Receiver#123\", \"userID\": \"User#123\", \"eventName\": \"showers\"}",
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
		"Sad Path - Bad Body": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "{\"receiverId\": false}",
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Error Getting User": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "{\"receiverId\": \"Receiver#123\", \"userID\": \"User#Error\", \"eventName\": \"showers\"}",
			},
			expectedResponse: response.CreateInternalServerErrorResponse(),
		},
		"Sad Path - User Is Not A Care Giver": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "{\"receiverId\": \"Receiver#123\", \"userID\": \"User#NotACareGiver\", \"eventName\": \"showers\"}",
			},
			expectedResponse: response.CreateAccessDeniedResponse(),
		},
		"Sad Path - Bad Event Name": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "{\"receiverId\": \"Receiver#123\", \"userID\": \"User#123\", \"eventName\": \"badEventName\"}",
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Bad Event Data": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "{\"receiverId\": \"Receiver#123\", \"userID\": \"User#123\", \"eventName\": \"weights\", \"event\": {\"badField\": \"false\"}}",
			},
			expectedResponse: response.CreateBadRequestResponse(),
		},
		"Sad Path - Error Updating Receiver": {
			request: events.APIGatewayProxyRequest{
				HTTPMethod: http.MethodPost,
				Body:       "{\"receiverId\": \"Receiver#123Error\", \"userID\": \"User#123\", \"eventName\": \"showers\"}",
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
			resp, err := HandleReceiverEvent(context.Background(), params)

			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResponse, resp)
		})
	}
}
