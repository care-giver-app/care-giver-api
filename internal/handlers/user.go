package handlers

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/log"
	"github.com/care-giver-app/care-giver-api/internal/receiver"
	"github.com/care-giver-app/care-giver-api/internal/repository"
	"github.com/care-giver-app/care-giver-api/internal/response"
	"github.com/care-giver-app/care-giver-api/internal/user"
	"go.uber.org/zap"
)

const (
	createUser            = "create user"
	getUser               = "get user"
	addPrimaryReceiver    = "add primary receiver"
	addAdditionalReceiver = "add additional receiver"
)

type CreateUserRequest struct {
	Email     string `json:"email" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type CreateUserResponse struct {
	UserID string `json:"userId"`
	Status string `json:"status"`
}

type PrimaryReceiverRequest struct {
	UserID    string `json:"userId" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type PrimaryReceiverResponse struct {
	ReceiverID string `json:"receiverId"`
	Status     string `json:"status"`
}

type AdditionalReceiverRequest struct {
	UserID     string `json:"userId" validate:"required"`
	ReceiverID string `json:"receiverId" validate:"required"`
}

func HandleCreateUser(ctx context.Context, params HandlerParams) (events.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, createUser)

	var createUserRequest CreateUserRequest
	err := readRequestBody(params.Request.Body, &createUserRequest)
	if err != nil {
		params.AppCfg.Logger.Error(requestBodyError, zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	user, err := user.NewUser(createUserRequest.Email, createUserRequest.FirstName, createUserRequest.LastName)
	if err != nil {
		params.AppCfg.Logger.Error("error creating new user", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}
	params.AppCfg.Logger = params.AppCfg.Logger.With(zap.Any(log.UserIDLogKey, user.UserID))

	err = params.UserRepo.CreateUser(*user)
	if err != nil {
		params.AppCfg.Logger.Error("error creating new user in db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	resp := CreateUserResponse{
		UserID: user.UserID,
		Status: response.Success,
	}

	params.AppCfg.Logger.Sugar().Infof(handlerSuccessful, createUser)
	return response.FormatResponse(resp, http.StatusOK), nil
}

func HandleGetUser(ctx context.Context, params HandlerParams) (events.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, getUser)

	uid, err := validatePathParameters(params.Request, user.ParamID, user.DBPrefix)
	if err != nil {
		params.AppCfg.Logger.Error(pathParametersError, zap.String(log.ParamIDLogKey, user.ParamID), zap.Any(log.PathParametersLogKey, params.Request.PathParameters), zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	u, err := params.UserRepo.GetUser(uid)
	if err != nil {
		params.AppCfg.Logger.Error(userDatbaseError, zap.String(log.UserIDLogKey, uid), zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	params.AppCfg.Logger.Sugar().Infof(handlerSuccessful, getUser)
	return response.FormatResponse(u, http.StatusOK), nil
}

func HandleUserPrimaryReceiver(ctx context.Context, params HandlerParams) (events.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, addPrimaryReceiver)

	var primaryReceiverRequest PrimaryReceiverRequest
	err := readRequestBody(params.Request.Body, &primaryReceiverRequest)
	if err != nil {
		params.AppCfg.Logger.Error(requestBodyError, zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	receiver := receiver.NewReceiver(primaryReceiverRequest.FirstName, primaryReceiverRequest.LastName)
	params.AppCfg.Logger = params.AppCfg.Logger.With(zap.Any(log.ReceiverIDLogKey, receiver.ReceiverID))

	err = params.ReceiverRepo.CreateReceiver(*receiver)
	if err != nil {
		params.AppCfg.Logger.Error("error creating receiver in db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	err = params.UserRepo.UpdateReceiverList(primaryReceiverRequest.UserID, receiver.ReceiverID, repository.PrimaryReceiverList)
	if err != nil {
		params.AppCfg.Logger.Error("error updating user primary receiver list", zap.Error(err))
		// TODO: delete newly created receiver item
		return response.CreateInternalServerErrorResponse(), nil
	}

	resp := PrimaryReceiverResponse{
		ReceiverID: string(receiver.ReceiverID),
		Status:     response.Success,
	}

	params.AppCfg.Logger.Sugar().Infof(handlerSuccessful, addPrimaryReceiver)
	return response.FormatResponse(resp, http.StatusOK), nil
}

func HandleUserAdditionalReceiver(ctx context.Context, params HandlerParams) (events.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, addAdditionalReceiver)

	var additionalReceiverRequest AdditionalReceiverRequest
	err := readRequestBody(params.Request.Body, &additionalReceiverRequest)
	if err != nil {
		params.AppCfg.Logger.Error(requestBodyError, zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	err = params.UserRepo.UpdateReceiverList(additionalReceiverRequest.UserID, additionalReceiverRequest.ReceiverID, repository.AdditionalReceiverList)
	if err != nil {
		params.AppCfg.Logger.Error("error updating user additional receiver list", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	params.AppCfg.Logger.Sugar().Infof(handlerSuccessful, addAdditionalReceiver)
	return response.FormatResponse(map[string]string{
		"status": response.Success,
	}, http.StatusOK), nil
}
