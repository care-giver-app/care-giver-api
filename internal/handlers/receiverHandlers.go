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

	r, err := params.ReceiverRepo.GetReceiver(rid)
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

	var rer ReceiverEventRequest
	err = readRequestBody(params.Request.Body, &rer)
	if err != nil {
		params.AppCfg.Logger.Error("error reading request body", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	u, err := params.UserRepo.GetUser(rer.UserID)
	if err != nil {
		params.AppCfg.Logger.Error("error retrieving user from db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	if !u.IsACareGiver(rer.ReceiverID) {
		params.AppCfg.Logger.Sugar().Errorf("user %s is unauthorized to access receiver %s", u.UserID, rer.ReceiverID)
		return response.CreateAccessDeniedResponse(), nil
	}

	newEvent, err := receiver.GenerateEvent(rer.EventName, rer.ReceiverID, u.UserID)
	if err != nil {
		params.AppCfg.Logger.Error("unsupported event name provided", zap.Any(log.EventLogKey, rer.EventName), zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	err = newEvent.ProcessEvent(rer.EventData)
	if err != nil {
		params.AppCfg.Logger.Error("error processing event data", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	err = params.EventRepo.AddEvent(newEvent)
	if err != nil {
		params.AppCfg.Logger.Error("error adding event to db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	// err = params.ReceiverRepo.UpdateReceiver(rer.ReceiverID, newEvent, rer.EventName)
	// if err != nil {
	// 	params.AppCfg.Logger.Error("error updating receiver in db", zap.Error(err))
	// 	return response.CreateInternalServerErrorResponse(), nil
	// }

	params.AppCfg.Logger.Info("processed add receiver event successfully")
	return response.FormatResponse(map[string]string{
		"status": "success",
	}, http.StatusOK), nil
}
