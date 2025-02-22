package handlers

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
	"github.com/care-giver-app/care-giver-api/internal/log"
	"github.com/care-giver-app/care-giver-api/internal/receiver"
	"github.com/care-giver-app/care-giver-api/internal/repository"
	"github.com/care-giver-app/care-giver-api/internal/response"
	"github.com/care-giver-app/care-giver-api/internal/user"
	"go.uber.org/zap"
)

type CreateUserRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type PrimaryReceiverRequest struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type AdditionalReceiverRequest struct {
	UserID     string `json:"userId"`
	ReceiverID string `json:"receiverId"`
}

func HandleCreateUser(ctx context.Context, appCfg *appconfig.AppConfig, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	appCfg.Logger.Info("handling create user event")

	err := validateMethod(req, http.MethodPost)
	if err != nil {
		appCfg.Logger.Error("error validating request method", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	var createUserRequest CreateUserRequest
	err = readRequestBody(req.Body, &createUserRequest)
	if err != nil {
		appCfg.Logger.Error("error reading request body", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	user, err := user.NewUser(createUserRequest.Email, createUserRequest.Password, createUserRequest.FirstName, createUserRequest.LastName)
	if err != nil {
		appCfg.Logger.Error("error creating new user", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}
	appCfg.Logger = appCfg.Logger.With(zap.Any(log.UserIDLogKey, user.UserID))

	ur, err := repository.UserRepositoryFromContext(ctx)
	if err != nil {
		appCfg.Logger.Error("error retrieving user repo from context", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	err = ur.CreateUser(*user)
	if err != nil {
		appCfg.Logger.Error("error creating new user in db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	appCfg.Logger.Info("processed create user event successfully")
	return response.FormatResponse(map[string]string{
		"status": "success",
	}, http.StatusOK), nil
}

func HandleGetUser(ctx context.Context, appCfg *appconfig.AppConfig, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	appCfg.Logger.Info("handling get user event")

	err := validateMethod(req, http.MethodGet)
	if err != nil {
		appCfg.Logger.Error("error validating request method", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	uid, err := validatePathParameters(req, user.ParamId, user.DBPrefix)
	if err != nil {
		appCfg.Logger.Error("error validating path parameters", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	ur, err := repository.UserRepositoryFromContext(ctx)
	if err != nil {
		appCfg.Logger.Error("error retrieving user repo from context", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	u, err := ur.GetUser(uid)
	if err != nil {
		appCfg.Logger.Error("error retrieving user from db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	appCfg.Logger.Info("processed get user event successfully")
	return response.FormatResponse(u, http.StatusOK), nil
}

func HandleUserPrimaryReceiver(ctx context.Context, appCfg *appconfig.AppConfig, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	appCfg.Logger.Info("handling add primary receiver event")

	err := validateMethod(req, http.MethodPost)
	if err != nil {
		appCfg.Logger.Error("error validating request method", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	var primaryReceiverRequest PrimaryReceiverRequest
	err = readRequestBody(req.Body, &primaryReceiverRequest)
	if err != nil {
		appCfg.Logger.Error("error reading request body", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	rr, err := repository.ReceiverRepositoryFromContext(ctx)
	if err != nil {
		appCfg.Logger.Error("error retrieving receiver repo from context", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	ur, err := repository.UserRepositoryFromContext(ctx)
	if err != nil {
		appCfg.Logger.Error("error retrieving user repo from context", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	receiver := receiver.NewReceiver(primaryReceiverRequest.FirstName, primaryReceiverRequest.LastName)
	appCfg.Logger = appCfg.Logger.With(zap.Any(log.ReceiverIDLogKey, receiver.ReceiverID))

	err = rr.CreateReceiver(*receiver)
	if err != nil {
		appCfg.Logger.Error("error creating receiver in db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	err = ur.UpdateReceiverList(primaryReceiverRequest.UserID, string(receiver.ReceiverID), repository.PrimaryReceiverList)
	if err != nil {
		appCfg.Logger.Error("error updating user primary receiver list", zap.Error(err))
		// TODO: delete newly created receiver item
		return response.CreateInternalServerErrorResponse(), nil
	}

	appCfg.Logger.Info("processed add primary receiver event successfully")
	return response.FormatResponse(map[string]string{
		"status": "success",
	}, http.StatusOK), nil
}

func HandleUserAdditionalReceiver(ctx context.Context, appCfg *appconfig.AppConfig, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	appCfg.Logger.Info("handling add additional receiver event")

	err := validateMethod(req, http.MethodPost)
	if err != nil {
		appCfg.Logger.Error("error validating request method", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	var additionalReceiverRequest AdditionalReceiverRequest
	err = readRequestBody(req.Body, &additionalReceiverRequest)
	if err != nil {
		appCfg.Logger.Error("error reading request body", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	ur, err := repository.UserRepositoryFromContext(ctx)
	if err != nil {
		appCfg.Logger.Error("error retrieving user repo from context", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	err = ur.UpdateReceiverList(additionalReceiverRequest.UserID, additionalReceiverRequest.ReceiverID, repository.AdditionalReceiverList)
	if err != nil {
		appCfg.Logger.Error("error updating user additional receiver list", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	appCfg.Logger.Info("processed add additional receiver event successfully")
	return response.FormatResponse(map[string]string{
		"status": "success",
	}, http.StatusOK), nil
}
