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

func formatResponse(resp interface{}, statusCode int) events.APIGatewayProxyResponse {
	respJson, err := json.Marshal(resp)
	if err != nil {
		panic("unable to create response")
	}
	respStr := string(respJson)

	return events.APIGatewayProxyResponse{
		Body:       respStr,
		StatusCode: statusCode,
	}
}

func CreateBadRequestResponse() events.APIGatewayProxyResponse {
	resp := &ErrorResponse{
		Status: "Bad Request",
	}

	return formatResponse(resp, http.StatusBadRequest)
}

func CreateResourceNotFoundResponse() events.APIGatewayProxyResponse {
	resp := &ErrorResponse{
		Status: "Resource Not Found",
	}
	return formatResponse(resp, http.StatusNotFound)
}

func CreateInternalServerErrorResponse() events.APIGatewayProxyResponse {
	resp := &ErrorResponse{
		Status: "Internal Server Error",
	}

	return formatResponse(resp, http.StatusInternalServerError)
}
