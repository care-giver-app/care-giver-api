package response

import (
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

type TestResponseStruct struct {
	Body   string `json:"body"`
	Status string `json:"status"`
}

func TestFormatResponse(t *testing.T) {
	expectedResponse := events.APIGatewayProxyResponse{
		Body:       "{\"body\":\"MyTestBody\",\"status\":\"Success\"}",
		StatusCode: http.StatusOK,
	}

	respStruct := &TestResponseStruct{
		Body:   "MyTestBody",
		Status: "Success",
	}

	resp := FormatResponse(respStruct, http.StatusOK)
	assert.Equal(t, expectedResponse, resp)
}

func TestResponses(t *testing.T) {
	tests := map[string]struct {
		function         func() events.APIGatewayProxyResponse
		expectedResponse events.APIGatewayProxyResponse
	}{
		"Bad Request": {
			function: CreateBadRequestResponse,
			expectedResponse: events.APIGatewayProxyResponse{
				Body:       "{\"status\":\"Bad Request\"}",
				StatusCode: http.StatusBadRequest,
			},
		},
		"Internal Server Error": {
			function: CreateInternalServerErrorResponse,
			expectedResponse: events.APIGatewayProxyResponse{
				Body:       "{\"status\":\"Internal Server Error\"}",
				StatusCode: http.StatusInternalServerError,
			},
		},
		"Resource Not Found": {
			function: CreateResourceNotFoundResponse,
			expectedResponse: events.APIGatewayProxyResponse{
				Body:       "{\"status\":\"Resource Not Found\"}",
				StatusCode: http.StatusNotFound,
			},
		},
		"Access Denied": {
			function: CreateAccessDeniedResponse,
			expectedResponse: events.APIGatewayProxyResponse{
				Body:       "{\"status\":\"Access Denied\"}",
				StatusCode: http.StatusForbidden,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			resp := tc.function()
			assert.Equal(t, tc.expectedResponse, resp)
		})
	}
}
