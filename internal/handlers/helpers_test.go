package handlers

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

func ptr(v string) *string {
	return &v
}

func TestValidatePathParameters(t *testing.T) {
	tests := map[string]struct {
		request     events.APIGatewayProxyRequest
		param       string
		idPrefix    string
		expectedId  string
		expectedErr *string
	}{
		"Happy Path - Validated Parameter with #": {
			request: events.APIGatewayProxyRequest{PathParameters: map[string]string{
				"testParam": "test#123",
			}},
			param:      "testParam",
			idPrefix:   "test",
			expectedId: "test#123",
		},
		"Happy Path - Validated Parameter with %23": {
			request: events.APIGatewayProxyRequest{PathParameters: map[string]string{
				"testParam": "test%23123",
			}},
			param:      "testParam",
			idPrefix:   "test",
			expectedId: "test#123",
		},
		"Sad Path - No Path Parameters": {
			request:     events.APIGatewayProxyRequest{PathParameters: map[string]string{}},
			expectedErr: ptr("no path parameters provided"),
		},
		"Sad Path - Too Many Parameters": {
			request: events.APIGatewayProxyRequest{PathParameters: map[string]string{
				"testParam":    "test#123",
				"anotherParam": "test#456",
			}},
			param:       "testParam",
			idPrefix:    "test",
			expectedErr: ptr("too many path parameters provided"),
		},
		"Sad Path - Wrong Format": {
			request: events.APIGatewayProxyRequest{PathParameters: map[string]string{
				"testParam": "test!123",
			}},
			param:       "testParam",
			idPrefix:    "test",
			expectedErr: ptr("id is not formatted correctly"),
		},
		"Sad Path - Wrong Parameter Provided": {
			request: events.APIGatewayProxyRequest{PathParameters: map[string]string{
				"badParam": "test#123",
			}},
			param:       "testParam",
			idPrefix:    "test",
			expectedErr: ptr("invalid path parameters"),
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			id, err := validatePathParameters(tc.request, tc.param, tc.idPrefix)

			if tc.expectedErr != nil {
				assert.Equal(t, *tc.expectedErr, err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedId, id)
			}
		})
	}
}

func TestValidateQueryParameters(t *testing.T) {
	tests := map[string]struct {
		request       events.APIGatewayProxyRequest
		param         string
		expectedValue string
		expectedErr   string
	}{
		"Happy Path": {
			request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"queryParam": "queryValue",
				},
			},
			param:         "queryParam",
			expectedValue: "queryValue",
		},
		"Sad Path - No Query Parameters": {
			request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{},
			},
			param:       "queryParam",
			expectedErr: "no query parameters provided",
		},
		"Sad Path - Empty Query Param": {
			request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"queryParam": "",
				},
			},
			param:       "queryParam",
			expectedErr: "query parameter value is empty",
		},
		"Sad Path - Not Found": {
			request: events.APIGatewayProxyRequest{
				QueryStringParameters: map[string]string{
					"queryParam": "queryValue",
				},
			},
			param:       "diffParam",
			expectedErr: "query parameter 'diffParam' not found",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			v, err := validateQueryParameters(tc.request, tc.param)
			if tc.expectedErr != "" {
				assert.Equal(t, tc.expectedErr, err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedValue, v)
			}
		})
	}
}

type TestRequestStruct struct {
	FieldOne   string            `json:"fieldOne" validate:"required"`
	FieldTwo   float64           `json:"fieldTwo" validate:"required"`
	FieldThree map[string]string `json:"fieldThree"`
}

func TestReadRequestBody(t *testing.T) {
	tests := map[string]struct {
		requestBody    string
		expectedStruct TestRequestStruct
		expectErr      bool
	}{
		"Happy Path - Only Required Fields": {
			requestBody: "{\"fieldOne\": \"test\", \"fieldTwo\": 20.4}",
			expectedStruct: TestRequestStruct{
				FieldOne: "test",
				FieldTwo: 20.4,
			},
		},
		"Happy Path - All Fields": {
			requestBody: "{\"fieldOne\": \"test\", \"fieldTwo\": 20.4, \"fieldThree\": {\"Hello\": \"World\"}}",
			expectedStruct: TestRequestStruct{
				FieldOne: "test",
				FieldTwo: 20.4,
				FieldThree: map[string]string{
					"Hello": "World",
				},
			},
		},
		"Sad Path - Missing Required Field": {
			requestBody: "{\"fieldTwo\": 20.4, \"fieldThree\": {\"Hello\": \"World\"}}",
			expectErr:   true,
		},
		"Sad Path - Wrong Field Type": {
			requestBody: "{\"fieldOne\": \"test\", \"fieldTwo\": \"NAN\", \"fieldThree\": {\"Hello\": \"World\"}}",
			expectErr:   true,
		},
		"Sad Path - Unknown Field": {
			requestBody: "{\"fieldOne\": \"test\", \"fieldTwo\": 20.4, \"fieldFour\": {\"Hello\": \"World\"}}",
			expectErr:   true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var requestStruct TestRequestStruct
			err := readRequestBody(tc.requestBody, &requestStruct)

			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedStruct, requestStruct)
			}
		})
	}
}
