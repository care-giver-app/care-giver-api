package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator/v10"
)

const (
	idSeparator           = "#"
	idSeparatorUrlEscaped = "%23"
)

func validatePathParameters(request events.APIGatewayProxyRequest, param string, idPrefix string) (string, error) {
	switch len(request.PathParameters) {
	case 0:
		return "", errors.New("no path parameters provided")
	case 1:
		paramPrefixURLEscaped := fmt.Sprintf("%s%s", idPrefix, idSeparatorUrlEscaped)
		dbPrefix := fmt.Sprintf("%s%s", idPrefix, idSeparator)
		idRegex := fmt.Sprintf(`^%s[a-zA-Z0-9-]+$`, dbPrefix)

		if id, found := request.PathParameters[param]; found {
			id = strings.Replace(id, paramPrefixURLEscaped, dbPrefix, 1)
			validFormat := regexp.MustCompile(idRegex).MatchString(id)
			if !validFormat {
				return "", errors.New("id is not formatted correctly")
			}
			return id, nil
		}
		return "", errors.New("invalid path parameters")
	default:
		return "", errors.New("too many path parameters provided")
	}
}

func validateQueryParameters(request events.APIGatewayProxyRequest, param string) (string, error) {
	if len(request.QueryStringParameters) == 0 {
		return "", errors.New("no query parameters provided")
	}

	if value, found := request.QueryStringParameters[param]; found {
		if value == "" {
			return "", errors.New("query parameter value is empty")
		}
		return value, nil
	}

	return "", fmt.Errorf("query parameter '%s' not found", param)
}

func readRequestBody(requestBody string, requestStruct interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader([]byte(requestBody)))
	decoder.DisallowUnknownFields()
	err := decoder.Decode(requestStruct)
	if err != nil {
		return err
	}

	validate := validator.New()
	err = validate.Struct(requestStruct)
	if err != nil {
		return err
	}

	return nil
}

func validateTimestamps(startTime, endTime string) error {
	st, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return fmt.Errorf("failed to RFC3339 parse start time of %s", startTime)
	}

	et, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return fmt.Errorf("failed to RFC3339 parse end time of %s", endTime)
	}

	if et.Before(st) {
		return fmt.Errorf("end time is before start time - %s is before %s", endTime, startTime)
	}

	return nil
}
