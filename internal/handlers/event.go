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
	"github.com/care-giver-app/care-giver-golang-common/pkg/repository"
	"github.com/care-giver-app/care-giver-golang-common/pkg/user"
	"go.uber.org/zap"
)

const (
	addReceiverEvent    = "add receiver event"
	deleteReceiverEvent = "delete receiver event"
	getReceiverEvents   = "get receiver events"
	getEventConfigs     = "get event configs"
)

type ReceiverEventRequest struct {
	ReceiverID string            `json:"receiverId" validate:"required"`
	UserID     string            `json:"userId" validate:"required"`
	Type       string            `json:"type" validate:"required"`
	StartTime  string            `json:"startTime" validate:"required"`
	EndTime    string            `json:"endTime" validate:"required"`
	Data       []event.DataPoint `json:"data"`
	Note       string            `json:"note"`
}

type ReceiverEventResponse struct {
	ReceiverID string `json:"receiverId"`
	EventID    string `json:"eventId"`
	Status     string `json:"status"`
}

// @Summary Create a care event for a receiver
// @Tags events
// @Security BearerAuth
// @Param body body ReceiverEventRequest true "Event details. startTime and endTime must be RFC3339. type must match a valid event config."
// @Success 200 {object} ReceiverEventResponse
// @Failure 400 {string} string "Bad request"
// @Failure 403 {string} string "Access denied"
// @Failure 500 {string} string "Internal server error"
// @Router /event [post]
func HandleReceiverEvent(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, addReceiverEvent)
	params.AppCfg.Logger.Info("handling add receiver event")

	var rer ReceiverEventRequest
	err := readRequestBody(params.Request.Body, &rer)
	if err != nil {
		params.AppCfg.Logger.Error(requestBodyError, zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	err = validateTimestamps(rer.StartTime, rer.EndTime)
	if err != nil {
		params.AppCfg.Logger.Error(requestBodyError, zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	u, err := params.UserRepo.GetUser(rer.UserID)
	if err != nil {
		params.AppCfg.Logger.Error(userDatabaseError, zap.String(log.UserIDLogKey, rer.UserID), zap.Error(err))
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
	if len(rer.Data) > 0 {
		opts = append(opts, event.WithData(rer.Data))
	}

	if rer.Note != "" {
		opts = append(opts, event.WithNote(rer.Note))
	}

	newEvent, err := event.NewEntry(rer.ReceiverID, u.UserID, rer.Type, rer.StartTime, rer.EndTime, opts...)
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

// @Summary Delete a care event
// @Tags events
// @Security BearerAuth
// @Param eventId path string true "Event ID (format: Event#<uuid>)"
// @Param receiverId query string true "Receiver ID (format: Receiver#<uuid>)"
// @Param userId query string true "ID of the requesting user (format: User#<uuid>)"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "Bad request"
// @Failure 403 {string} string "Access denied"
// @Failure 500 {string} string "Internal server error"
// @Router /event/{eventId} [delete]
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
		params.AppCfg.Logger.Error(userDatabaseError, zap.String(log.UserIDLogKey, uid), zap.Error(err))
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

// @Summary Get all care events for a receiver
// @Tags events
// @Security BearerAuth
// @Param receiverId path string true "Receiver ID (format: Receiver#<uuid>)"
// @Param userId query string true "ID of the requesting user (format: User#<uuid>)"
// @Param startTime query string false "Start of time range (RFC3339). Must be provided with endTime."
// @Param endTime query string false "End of time range (RFC3339). Must be provided with startTime."
// @Success 200 {object} []event.Entry
// @Failure 400 {string} string "Bad request"
// @Failure 403 {string} string "Access denied"
// @Failure 500 {string} string "Internal server error"
// @Router /events/{receiverId} [get]
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
		params.AppCfg.Logger.Error(userDatabaseError, zap.String(log.UserIDLogKey, uid), zap.Error(err))
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

	bound := repository.TimestampBound{}
	startTime := params.Request.QueryStringParameters["startTime"]
	endTime := params.Request.QueryStringParameters["endTime"]
	if startTime != "" && endTime != "" {
		if err := validateTimestamps(startTime, endTime); err != nil {
			params.AppCfg.Logger.Error("invalid date bound query params", zap.Error(err))
			return response.CreateBadRequestResponse(), nil
		}
		bound = repository.TimestampBound{Lower: startTime, Upper: endTime}
	}
	eventsList, err := params.EventRepo.GetEvents(rid, bound)
	if err != nil {
		params.AppCfg.Logger.Error("error retrieving events from db", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	params.AppCfg.Logger.Sugar().Infof(handlerSuccessful, getReceiverEvents)
	return response.FormatResponse(eventsList, http.StatusOK), nil
}

// @Summary Get all available event type configurations
// @Tags events
// @Security BearerAuth
// @Success 200 {object} []event.EventConfig
// @Failure 500 {string} string "Internal server error"
// @Router /events/configs [get]
func HandleGetEventConfigs(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, getEventConfigs)

	eventConfigs, err := event.GetAllConfigs()
	if err != nil {
		params.AppCfg.Logger.Error("error retrieving event configs", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	return response.FormatResponse(eventConfigs, http.StatusOK), nil
}
