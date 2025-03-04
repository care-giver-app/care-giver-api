package handlers

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/log"
	"github.com/care-giver-app/care-giver-api/internal/receiver"
	"github.com/care-giver-app/care-giver-api/internal/response"
	"go.uber.org/zap"
)

type ReceiverEventRequest struct {
	ReceiverID string                 `json:"receiverId" validate:"required"`
	UserID     string                 `json:"userId" validate:"required"`
	EventName  string                 `json:"eventName" validate:"required"`
	EventData  map[string]interface{} `json:"event"`
}

func HandleReceiver(ctx context.Context, params HandlerParams) (events.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Info("handling get receiver event")
	err := validateMethod(params.Request, http.MethodGet)
	if err != nil {
		params.AppCfg.Logger.Error("error validating request method", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	rid, err := validatePathParameters(params.Request, receiver.ParamId, receiver.DBPrefix)
	if err != nil {
		params.AppCfg.Logger.Error("error validating path parameters", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	r, err := params.ReceiverRepo.GetReceiver(receiver.ReceiverID(rid))
	if err != nil {
		params.AppCfg.Logger.Error("error retrieving receiver from db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	params.AppCfg.Logger.Info("processed get receiver event successfully")
	return response.FormatResponse(r, http.StatusOK), nil
}

func HandleReceiverEvent(ctx context.Context, params HandlerParams) (events.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Info("handling add receiver event")
	err := validateMethod(params.Request, http.MethodPost)
	if err != nil {
		params.AppCfg.Logger.Error("error validating request method", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	var receiverEventRequest ReceiverEventRequest
	err = readRequestBody(params.Request.Body, &receiverEventRequest)
	if err != nil {
		params.AppCfg.Logger.Error("error reading request body", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	u, err := params.UserRepo.GetUser(receiverEventRequest.UserID)
	if err != nil {
		params.AppCfg.Logger.Error("error retrieving user from db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	if !u.IsACareGiver(receiver.ReceiverID(receiverEventRequest.ReceiverID)) {
		params.AppCfg.Logger.Sugar().Errorf("user %s is unauthorized to access receiver %s", u.UserID, receiverEventRequest.ReceiverID)
		return response.CreateAccessDeniedResponse(), nil
	}

	newEvent, ok := receiver.NewEventMap[receiverEventRequest.EventName]
	if !ok {
		params.AppCfg.Logger.Error("unsupported event name provided", zap.Any(log.EventLogKey, receiverEventRequest.EventName), zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	err = newEvent.ProcessEvent(receiverEventRequest.EventData)
	if err != nil {
		params.AppCfg.Logger.Error("error processing event data", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	err = params.ReceiverRepo.UpdateReceiver(receiverEventRequest.ReceiverID, newEvent, receiverEventRequest.EventName)
	if err != nil {
		params.AppCfg.Logger.Error("error updating receiver in db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	params.AppCfg.Logger.Info("processed add receiver event successfully")
	return response.FormatResponse(map[string]string{
		"status": "success",
	}, http.StatusOK), nil
}
