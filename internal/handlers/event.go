package handlers

import (
	"context"
	"net/http"

	awsevents "github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/event"
	"github.com/care-giver-app/care-giver-api/internal/log"
	"github.com/care-giver-app/care-giver-api/internal/receiver"
	"github.com/care-giver-app/care-giver-api/internal/response"
	"github.com/care-giver-app/care-giver-api/internal/user"
	"go.uber.org/zap"
)

type ReceiverEventRequest struct {
	ReceiverID string            `json:"receiverId" validate:"required"`
	UserID     string            `json:"userId" validate:"required"`
	Type       string            `json:"type" validate:"required"`
	Timestamp  string            `json:"timestamp"`
	Data       []event.DataPoint `json:"data"`
}

type ReceiverEventResponse struct {
	ReceiverID string `json:"receiverId"`
	EventID    string `json:"eventId"`
	Status     string `json:"status"`
}

func HandleReceiverEvent(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Info("handling add receiver event")

	var rer ReceiverEventRequest
	err := readRequestBody(params.Request.Body, &rer)
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

	opts := []event.EntryOption{}
	if rer.Timestamp != "" {
		opts = append(opts, event.WithTimestamp(rer.Timestamp))
	}

	if len(rer.Data) > 0 {
		opts = append(opts, event.WithData(rer.Data))
	}

	newEvent, err := event.NewEntry(rer.ReceiverID, u.UserID, rer.Type, opts...)
	if err != nil {
		params.AppCfg.Logger.Error("error creating new event entry", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	err = params.EventRepo.AddEvent(newEvent)
	if err != nil {
		params.AppCfg.Logger.Error("error adding event to db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	params.AppCfg.Logger.Info("processed add receiver event successfully")
	return response.FormatResponse(ReceiverEventResponse{
		ReceiverID: rer.ReceiverID,
		EventID:    newEvent.EventID,
		Status:     "Success",
	}, http.StatusOK), nil
}

func HandleDeleteReceiverEvent(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Info("handling delete receiver event")

	eid, err := validatePathParameters(params.Request, event.ParamID, event.DBPrefix)
	if err != nil {
		params.AppCfg.Logger.Error("error validating path parameters", zap.String(log.ParamIDLogKey, event.ParamID), zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	rid, err := validateQueryParameters(params.Request, receiver.ParamID)
	if err != nil {
		params.AppCfg.Logger.Error("error validating query parameters", zap.String(log.ParamIDLogKey, receiver.ParamID), zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	uid, err := validateQueryParameters(params.Request, user.ParamID)
	if err != nil {
		params.AppCfg.Logger.Error("error validating query parameters", zap.String(log.ParamIDLogKey, user.ParamID), zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	u, err := params.UserRepo.GetUser(uid)
	if err != nil {
		params.AppCfg.Logger.Error("error retrieving user from db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	if !u.IsACareGiver(rid) {
		params.AppCfg.Logger.Sugar().Errorf("user %s is unauthorized to delete event for receiver %s", u.UserID, rid)
		return response.CreateAccessDeniedResponse(), nil
	}

	err = params.EventRepo.DeleteEvent(rid, eid)
	if err != nil {
		params.AppCfg.Logger.Error("error deleting event from db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	params.AppCfg.Logger.Info("processed delete receiver event successfully")
	return response.FormatResponse(
		map[string]string{
			"status": "success",
		}, http.StatusOK,
	), nil
}

func HandleGetReceiverEvents(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Info("handling get receiver events")

	rid, err := validatePathParameters(params.Request, receiver.ParamID, receiver.DBPrefix)
	if err != nil {
		params.AppCfg.Logger.Error("error validating path parameters", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	uid, err := validateQueryParameters(params.Request, user.ParamID)
	if err != nil {
		params.AppCfg.Logger.Error("error validating query parameters", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	u, err := params.UserRepo.GetUser(uid)
	if err != nil {
		params.AppCfg.Logger.Error("error retrieving user from db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	if !u.IsACareGiver(rid) {
		params.AppCfg.Logger.Sugar().Errorf("user %s is unauthorized to access events for receiver %s", u.UserID, rid)
		return response.CreateAccessDeniedResponse(), nil
	}

	eventsList, err := params.EventRepo.GetEvents(rid)
	if err != nil {
		params.AppCfg.Logger.Error("error retrieving events from db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	params.AppCfg.Logger.Info("processed get receiver events successfully")
	return response.FormatResponse(eventsList, http.StatusOK), nil
}
