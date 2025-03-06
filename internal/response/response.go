package response

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ErrorResponse struct {
	DeveloperText string `json:"developerText,omitempty"`
	Status        string `json:"status"`
}

func FormatResponse(resp interface{}, statusCode int) events.APIGatewayProxyResponse {
	respJson, err := json.Marshal(resp)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "Failed to create response body",
			StatusCode: statusCode,
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       string(respJson),
		StatusCode: statusCode,
	}
}

func CreateBadRequestResponse() events.APIGatewayProxyResponse {
	resp := &ErrorResponse{
		Status: "Bad Request",
	}

	return FormatResponse(resp, http.StatusBadRequest)
}

func CreateResourceNotFoundResponse() events.APIGatewayProxyResponse {
	resp := &ErrorResponse{
		Status: "Resource Not Found",
	}
	return FormatResponse(resp, http.StatusNotFound)
}

func CreateInternalServerErrorResponse() events.APIGatewayProxyResponse {
	resp := &ErrorResponse{
		Status: "Internal Server Error",
	}

	return FormatResponse(resp, http.StatusInternalServerError)
}

func CreateAccessDeniedResponse() events.APIGatewayProxyResponse {
	resp := &ErrorResponse{
		Status: "Access Denied",
	}

	return FormatResponse(resp, http.StatusForbidden)
}
