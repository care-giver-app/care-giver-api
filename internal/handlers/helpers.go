package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
	"github.com/care-giver-app/care-giver-api/internal/repository"
	"github.com/go-playground/validator/v10"
)

const (
	idSeparator           = "#"
	idSeparatorUrlEscaped = "%23"
)

type HandlerParams struct {
	AppCfg       *appconfig.AppConfig
	Request      events.APIGatewayProxyRequest
	UserRepo     repository.UserRepositoryProvider
	ReceiverRepo repository.ReceiverRepositoryProvider
}

func validateMethod(request events.APIGatewayProxyRequest, method string) error {
	if request.HTTPMethod != method {
		return fmt.Errorf("unsupported request method: expected %s", method)
	}
	return nil
}

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
