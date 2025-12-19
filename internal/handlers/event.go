package handlers

import (
	"context"
	"net/http"

	awsevents "github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/response"
	"github.com/care-giver-app/care-giver-golang-common/pkg/event"
	"github.com/care-giver-app/care-giver-golang-common/pkg/log"
	"github.com/care-giver-app/care-giver-golang-common/pkg/receiver"
	"github.com/care-giver-app/care-giver-golang-common/pkg/relationship"
	"github.com/care-giver-app/care-giver-golang-common/pkg/user"
	"go.uber.org/zap"
)

const (
	addReceiverEvent    = "add receiver event"
	deleteReceiverEvent = "delete receiver event"
	getReceiverEvents   = "get receiver events"
)

type ReceiverEventRequest struct {
	ReceiverID string            `json:"receiverId" validate:"required"`
	UserID     string            `json:"userId" validate:"required"`
	Type       string            `json:"type" validate:"required"`
	Timestamp  string            `json:"timestamp"`
	Data       []event.DataPoint `json:"data"`
	Note       string            `json:"note"`
}

type ReceiverEventResponse struct {
	ReceiverID string `json:"receiverId"`
	EventID    string `json:"eventId"`
	Status     string `json:"status"`
}

func HandleReceiverEvent(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, addReceiverEvent)
	params.AppCfg.Logger.Info("handling add receiver event")

	var rer ReceiverEventRequest
	err := readRequestBody(params.Request.Body, &rer)
	if err != nil {
		params.AppCfg.Logger.Error(requestBodyError, zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	u, err := params.UserRepo.GetUser(rer.UserID)
	if err != nil {
		params.AppCfg.Logger.Error(userDatbaseError, zap.String(log.UserIDLogKey, rer.UserID), zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	relationships, err := params.RelationshipRepo.GetRelationshipsByUser(u.UserID)
	if err != nil {
		params.AppCfg.Logger.Error(relationshipDatabaseError, zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	if !relationship.IsACareGiver(u.UserID, rer.ReceiverID, relationships) {
		params.AppCfg.Logger.Sugar().Error(userNotCareGiverError, zap.String(log.ReceiverIDLogKey, rer.ReceiverID), zap.String(log.UserIDLogKey, u.UserID))
		return response.CreateAccessDeniedResponse(), nil
	}

	opts := []event.EntryOption{}
	if rer.Timestamp != "" {
		opts = append(opts, event.WithTimestamp(rer.Timestamp))
	}

	if len(rer.Data) > 0 {
		opts = append(opts, event.WithData(rer.Data))
	}

	if rer.Note != "" {
		opts = append(opts, event.WithNote(rer.Note))
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

	params.AppCfg.Logger.Sugar().Infof(handlerSuccessful, addReceiverEvent)
	return response.FormatResponse(ReceiverEventResponse{
		ReceiverID: rer.ReceiverID,
		EventID:    newEvent.EventID,
		Status:     response.Success,
	}, http.StatusOK), nil
}

func HandleDeleteReceiverEvent(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, deleteReceiverEvent)

	eid, err := validatePathParameters(params.Request, event.ParamID, event.DBPrefix)
	if err != nil {
		params.AppCfg.Logger.Error(pathParametersError, zap.String(log.ParamIDLogKey, event.ParamID), zap.Any(log.PathParametersLogKey, params.Request.PathParameters), zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	rid, err := validateQueryParameters(params.Request, receiver.ParamID)
	if err != nil {
		params.AppCfg.Logger.Error(queryParamsError, zap.String(log.ParamIDLogKey, receiver.ParamID), zap.Any(log.QueryParametersLogKey, params.Request.QueryStringParameters), zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	uid, err := validateQueryParameters(params.Request, user.ParamID)
	if err != nil {
		params.AppCfg.Logger.Error(queryParamsError, zap.String(log.ParamIDLogKey, user.ParamID), zap.Any(log.QueryParametersLogKey, params.Request.QueryStringParameters), zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	u, err := params.UserRepo.GetUser(uid)
	if err != nil {
		params.AppCfg.Logger.Error(userDatbaseError, zap.String(log.UserIDLogKey, uid), zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	relationships, err := params.RelationshipRepo.GetRelationshipsByUser(uid)
	if err != nil {
		params.AppCfg.Logger.Error(relationshipDatabaseError, zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	if !relationship.IsACareGiver(uid, rid, relationships) {
		params.AppCfg.Logger.Sugar().Error(userNotCareGiverError, zap.String(log.ReceiverIDLogKey, rid), zap.String(log.UserIDLogKey, u.UserID))
		return response.CreateAccessDeniedResponse(), nil
	}

	err = params.EventRepo.DeleteEvent(rid, eid)
	if err != nil {
		params.AppCfg.Logger.Error("error deleting event from db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	params.AppCfg.Logger.Sugar().Infof(handlerSuccessful, deleteReceiverEvent)
	return response.FormatResponse(
		map[string]string{
			"status": response.Success,
		}, http.StatusOK,
	), nil
}

func HandleGetReceiverEvents(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, getReceiverEvents)

	rid, err := validatePathParameters(params.Request, receiver.ParamID, receiver.DBPrefix)
	if err != nil {
		params.AppCfg.Logger.Error(pathParametersError, zap.String(log.ParamIDLogKey, receiver.ParamID), zap.Any(log.PathParametersLogKey, params.Request.PathParameters), zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	uid, err := validateQueryParameters(params.Request, user.ParamID)
	if err != nil {
		params.AppCfg.Logger.Error(queryParamsError, zap.String(log.ParamIDLogKey, user.ParamID), zap.Any(log.QueryParametersLogKey, params.Request.QueryStringParameters), zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	u, err := params.UserRepo.GetUser(uid)
	if err != nil {
		params.AppCfg.Logger.Error(userDatbaseError, zap.String(log.UserIDLogKey, uid), zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	relationships, err := params.RelationshipRepo.GetRelationshipsByUser(uid)
	if err != nil {
		params.AppCfg.Logger.Error(relationshipDatabaseError, zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	if !relationship.IsACareGiver(uid, rid, relationships) {
		params.AppCfg.Logger.Sugar().Error(userNotCareGiverError, zap.String(log.ReceiverIDLogKey, rid), zap.String(log.UserIDLogKey, u.UserID))
		return response.CreateAccessDeniedResponse(), nil
	}

	eventsList, err := params.EventRepo.GetEvents(rid)
	if err != nil {
		params.AppCfg.Logger.Error("error retrieving events from db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	params.AppCfg.Logger.Sugar().Infof(handlerSuccessful, getReceiverEvents)
	return response.FormatResponse(eventsList, http.StatusOK), nil
}
