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
	"go.uber.org/zap"
)

type ReceiverEventRequest struct {
	ReceiverID string                 `json:"receiverId" validate:"required"`
	UserID     string                 `json:"userId" validate:"required"`
	EventName  string                 `json:"eventName" validate:"required"`
	EventData  map[string]interface{} `json:"event"`
}

func HandleReceiver(ctx context.Context, appCfg *appconfig.AppConfig, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	appCfg.Logger.Info("handling get receiver event")
	err := validateMethod(req, http.MethodGet)
	if err != nil {
		appCfg.Logger.Error("error validating request method", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	rid, err := validatePathParameters(req, receiver.ParamId, receiver.DBPrefix)
	if err != nil {
		appCfg.Logger.Error("error validating path parameters", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	rr, err := repository.ReceiverRepositoryFromContext(ctx)
	if err != nil {
		appCfg.Logger.Error("error retrieving receiver repo from context", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	r, err := rr.GetReceiver(receiver.ReceiverID(rid))
	if err != nil {
		appCfg.Logger.Error("error retrieving receiver from db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	appCfg.Logger.Info("processed get receiver event successfully")
	return response.FormatResponse(r, http.StatusOK), nil
}

func HandleReceiverEvent(ctx context.Context, appCfg *appconfig.AppConfig, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	appCfg.Logger.Info("handling add receiver event")
	err := validateMethod(req, http.MethodPost)
	if err != nil {
		appCfg.Logger.Error("error validating request method", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	var receiverEventRequest ReceiverEventRequest
	err = readRequestBody(req.Body, &receiverEventRequest)
	if err != nil {
		appCfg.Logger.Error("error reading request body", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	ur, err := repository.UserRepositoryFromContext(ctx)
	if err != nil {
		appCfg.Logger.Error("error retrieving user repo from context", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	u, err := ur.GetUser(receiverEventRequest.UserID)
	if err != nil {
		appCfg.Logger.Error("error retrieving user from db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	if !u.IsACareGiver(receiver.ReceiverID(receiverEventRequest.ReceiverID)) {
		appCfg.Logger.Sugar().Errorf("user %s is unauthorized to access receiver %s", u.UserID, receiverEventRequest.ReceiverID)
		return response.CreateAccessDeniedResponse(), nil
	}

	newEvent, ok := receiver.NewEventMap[receiverEventRequest.EventName]
	if !ok {
		appCfg.Logger.Error("unsupported event name provided", zap.Any(log.EventLogKey, receiverEventRequest.EventName), zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	err = newEvent.ProcessEvent(receiverEventRequest.EventData)
	if err != nil {
		appCfg.Logger.Error("error processing event data", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	rr, err := repository.ReceiverRepositoryFromContext(ctx)
	if err != nil {
		appCfg.Logger.Error("error retrieving receiver repo from context", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	err = rr.UpdateReceiver(receiverEventRequest.ReceiverID, newEvent, receiverEventRequest.EventName)
	if err != nil {
		appCfg.Logger.Error("error updating receiver in db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	appCfg.Logger.Info("processed add receiver event successfully")
	return response.FormatResponse(map[string]string{
		"status": "success",
	}, http.StatusOK), nil
}
