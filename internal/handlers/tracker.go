package handlers

import (
	"context"
	"net/http"
	"strings"

	awsevents "github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/response"
	"github.com/care-giver-app/care-giver-golang-common/pkg/log"
	"github.com/care-giver-app/care-giver-golang-common/pkg/receiver"
	"github.com/care-giver-app/care-giver-golang-common/pkg/relationship"
	pkgtracker "github.com/care-giver-app/care-giver-golang-common/pkg/tracker"
	"github.com/care-giver-app/care-giver-golang-common/pkg/user"
	"go.uber.org/zap"
)

const (
	createTracker = "create tracker"
	listTrackers  = "list trackers"
	getTracker    = "get tracker"
	updateTracker = "update tracker"
	deleteTracker = "delete tracker"
)

type CreateTrackerRequest struct {
	ReceiverID      string                    `json:"receiverId"        validate:"required"`
	Name            string                    `json:"name"              validate:"required"`
	Kind            pkgtracker.TrackerKind    `json:"kind"              validate:"required"`
	Fields          []pkgtracker.TrackerField `json:"fields"`
	AlertThresholds []pkgtracker.AlertThreshold `json:"alertThresholds,omitempty"`
	Icon            string                    `json:"icon"              validate:"required"`
	Color           pkgtracker.ColorConfig    `json:"color"`
	IsActive        *bool                     `json:"isActive,omitempty"`
}

type UpdateTrackerRequest struct {
	Name            *string                    `json:"name,omitempty"`
	Fields          *[]pkgtracker.TrackerField  `json:"fields,omitempty"`
	AlertThresholds *[]pkgtracker.AlertThreshold `json:"alertThresholds,omitempty"`
	Icon            *string                    `json:"icon,omitempty"`
	Color           *pkgtracker.ColorConfig    `json:"color,omitempty"`
	IsActive        *bool                      `json:"isActive,omitempty"`
}

func isCareGiver(params HandlerParams, uid, rid string) (bool, error) {
	rels, err := params.RelationshipRepo.GetRelationshipsByUser(uid)
	if err != nil {
		return false, err
	}
	return relationship.IsACareGiver(uid, rid, rels), nil
}

func HandleCreateTracker(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, createTracker)

	uid, err := validateQueryParameters(params.Request, user.ParamID)
	if err != nil {
		params.AppCfg.Logger.Error(queryParamsError, zap.Error(err))
		return response.CreateAccessDeniedResponse(), nil
	}

	var req CreateTrackerRequest
	if err = readRequestBody(params.Request.Body, &req); err != nil {
		params.AppCfg.Logger.Error(requestBodyError, zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	authorized, err := isCareGiver(params, uid, req.ReceiverID)
	if err != nil {
		params.AppCfg.Logger.Error(relationshipDatabaseError, zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}
	if !authorized {
		params.AppCfg.Logger.Error(userNotCareGiverError, zap.String(log.ReceiverIDLogKey, req.ReceiverID), zap.String(log.UserIDLogKey, uid))
		return response.CreateAccessDeniedResponse(), nil
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	t := pkgtracker.NewTracker(req.ReceiverID, req.Name, req.Kind, req.Fields, req.AlertThresholds, req.Icon, req.Color, isActive)
	if err := t.Validate(); err != nil {
		params.AppCfg.Logger.Error("tracker validation failed", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	if err := params.TrackerRepo.CreateTracker(t); err != nil {
		if strings.Contains(err.Error(), "name already in use") {
			return response.FormatResponse(response.ErrorResponse{Status: "Conflict"}, http.StatusConflict), nil
		}
		params.AppCfg.Logger.Error("error creating tracker", zap.String(log.ReceiverIDLogKey, req.ReceiverID), zap.String(log.UserIDLogKey, uid))
		return response.CreateInternalServerErrorResponse(), nil
	}

	params.AppCfg.Logger.Info(handlerSuccessful, zap.String("action", createTracker), zap.String(log.UserIDLogKey, uid), zap.String(log.ReceiverIDLogKey, t.ReceiverID))
	return response.FormatResponse(t, http.StatusOK), nil
}

func HandleListTrackers(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, listTrackers)

	uid, err := validateQueryParameters(params.Request, user.ParamID)
	if err != nil {
		params.AppCfg.Logger.Error(queryParamsError, zap.Error(err))
		return response.CreateAccessDeniedResponse(), nil
	}

	rid, err := validatePathParameters(params.Request, receiver.ParamID, receiver.DBPrefix)
	if err != nil {
		params.AppCfg.Logger.Error(pathParametersError, zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	authorized, err := isCareGiver(params, uid, rid)
	if err != nil {
		params.AppCfg.Logger.Error(relationshipDatabaseError, zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}
	if !authorized {
		params.AppCfg.Logger.Error(userNotCareGiverError, zap.String(log.ReceiverIDLogKey, rid), zap.String(log.UserIDLogKey, uid))
		return response.CreateAccessDeniedResponse(), nil
	}

	trackers, err := params.TrackerRepo.ListTrackers(rid)
	if err != nil {
		params.AppCfg.Logger.Error("error listing trackers", zap.String(log.ReceiverIDLogKey, rid), zap.String(log.UserIDLogKey, uid))
		return response.CreateInternalServerErrorResponse(), nil
	}

	if trackers == nil {
		trackers = []pkgtracker.Tracker{}
	}

	params.AppCfg.Logger.Info(handlerSuccessful, zap.String("action", listTrackers), zap.String(log.UserIDLogKey, uid), zap.String(log.ReceiverIDLogKey, rid))
	return response.FormatResponse(trackers, http.StatusOK), nil
}

func HandleGetTracker(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, getTracker)

	uid, err := validateQueryParameters(params.Request, user.ParamID)
	if err != nil {
		params.AppCfg.Logger.Error(queryParamsError, zap.Error(err))
		return response.CreateAccessDeniedResponse(), nil
	}

	tid, err := validatePathParameters(params.Request, pkgtracker.ParamID, pkgtracker.DBPrefix)
	if err != nil {
		params.AppCfg.Logger.Error(pathParametersError, zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	rid, err := validateQueryParameters(params.Request, receiver.ParamID)
	if err != nil {
		params.AppCfg.Logger.Error(queryParamsError, zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	// authz before existence check — prevents probing (FR-014)
	authorized, err := isCareGiver(params, uid, rid)
	if err != nil {
		params.AppCfg.Logger.Error(relationshipDatabaseError, zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}
	if !authorized {
		params.AppCfg.Logger.Error(userNotCareGiverError, zap.String(log.ReceiverIDLogKey, rid), zap.String(log.UserIDLogKey, uid))
		return response.CreateAccessDeniedResponse(), nil
	}

	t, err := params.TrackerRepo.GetTracker(rid, tid)
	if err != nil {
		params.AppCfg.Logger.Error("error getting tracker", zap.String(log.ReceiverIDLogKey, rid), zap.String(log.UserIDLogKey, uid))
		return response.CreateInternalServerErrorResponse(), nil
	}
	if t == nil {
		return response.CreateResourceNotFoundResponse(), nil
	}

	params.AppCfg.Logger.Info(handlerSuccessful, zap.String("action", getTracker), zap.String(log.UserIDLogKey, uid), zap.String(log.ReceiverIDLogKey, rid))
	return response.FormatResponse(t, http.StatusOK), nil
}

func HandleUpdateTracker(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, updateTracker)

	uid, err := validateQueryParameters(params.Request, user.ParamID)
	if err != nil {
		params.AppCfg.Logger.Error(queryParamsError, zap.Error(err))
		return response.CreateAccessDeniedResponse(), nil
	}

	tid, err := validatePathParameters(params.Request, pkgtracker.ParamID, pkgtracker.DBPrefix)
	if err != nil {
		params.AppCfg.Logger.Error(pathParametersError, zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	rid, err := validateQueryParameters(params.Request, receiver.ParamID)
	if err != nil {
		params.AppCfg.Logger.Error(queryParamsError, zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	authorized, err := isCareGiver(params, uid, rid)
	if err != nil {
		params.AppCfg.Logger.Error(relationshipDatabaseError, zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}
	if !authorized {
		params.AppCfg.Logger.Error(userNotCareGiverError, zap.String(log.ReceiverIDLogKey, rid), zap.String(log.UserIDLogKey, uid))
		return response.CreateAccessDeniedResponse(), nil
	}

	var req UpdateTrackerRequest
	if err := readRequestBody(params.Request.Body, &req); err != nil {
		params.AppCfg.Logger.Error(requestBodyError, zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	existing, err := params.TrackerRepo.GetTracker(rid, tid)
	if err != nil {
		params.AppCfg.Logger.Error("error getting tracker for update", zap.String(log.ReceiverIDLogKey, rid), zap.String(log.UserIDLogKey, uid))
		return response.CreateInternalServerErrorResponse(), nil
	}
	if existing == nil {
		return response.CreateResourceNotFoundResponse(), nil
	}

	// name uniqueness check when name is being changed
	if req.Name != nil && !strings.EqualFold(*req.Name, existing.Name) {
		all, err := params.TrackerRepo.ListTrackers(rid)
		if err != nil {
			params.AppCfg.Logger.Error("error listing trackers for name check", zap.String(log.ReceiverIDLogKey, rid), zap.String(log.UserIDLogKey, uid))
			return response.CreateInternalServerErrorResponse(), nil
		}
		newNameLower := strings.ToLower(*req.Name)
		for _, tr := range all {
			if tr.TrackerID != tid && strings.ToLower(tr.Name) == newNameLower {
				return response.FormatResponse(response.ErrorResponse{Status: "Conflict"}, http.StatusConflict), nil
			}
		}
	}

	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.Fields != nil {
		existing.Fields = *req.Fields
	}
	if req.AlertThresholds != nil {
		existing.AlertThresholds = *req.AlertThresholds
	}
	if req.Icon != nil {
		existing.Icon = *req.Icon
	}
	if req.Color != nil {
		existing.Color = *req.Color
	}
	if req.IsActive != nil {
		existing.IsActive = *req.IsActive
	}

	if err := existing.Validate(); err != nil {
		params.AppCfg.Logger.Error("tracker validation failed after update", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	if err := params.TrackerRepo.UpdateTracker(existing); err != nil {
		params.AppCfg.Logger.Error("error updating tracker", zap.String(log.ReceiverIDLogKey, rid), zap.String(log.UserIDLogKey, uid))
		return response.CreateInternalServerErrorResponse(), nil
	}

	params.AppCfg.Logger.Info(handlerSuccessful, zap.String("action", updateTracker), zap.String(log.UserIDLogKey, uid), zap.String(log.ReceiverIDLogKey, rid))
	return response.FormatResponse(existing, http.StatusOK), nil
}

func HandleDeleteTracker(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, deleteTracker)

	uid, err := validateQueryParameters(params.Request, user.ParamID)
	if err != nil {
		params.AppCfg.Logger.Error(queryParamsError, zap.Error(err))
		return response.CreateAccessDeniedResponse(), nil
	}

	tid, err := validatePathParameters(params.Request, pkgtracker.ParamID, pkgtracker.DBPrefix)
	if err != nil {
		params.AppCfg.Logger.Error(pathParametersError, zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	rid, err := validateQueryParameters(params.Request, receiver.ParamID)
	if err != nil {
		params.AppCfg.Logger.Error(queryParamsError, zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	authorized, err := isCareGiver(params, uid, rid)
	if err != nil {
		params.AppCfg.Logger.Error(relationshipDatabaseError, zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}
	if !authorized {
		params.AppCfg.Logger.Error(userNotCareGiverError, zap.String(log.ReceiverIDLogKey, rid), zap.String(log.UserIDLogKey, uid))
		return response.CreateAccessDeniedResponse(), nil
	}

	existing, err := params.TrackerRepo.GetTracker(rid, tid)
	if err != nil {
		params.AppCfg.Logger.Error("error getting tracker for delete", zap.String(log.ReceiverIDLogKey, rid), zap.String(log.UserIDLogKey, uid))
		return response.CreateInternalServerErrorResponse(), nil
	}
	if existing == nil {
		return response.CreateResourceNotFoundResponse(), nil
	}

	hasEvents, err := params.EventRepo.HasEventsForTracker(rid, tid)
	if err != nil {
		params.AppCfg.Logger.Error("error checking events for tracker", zap.String(log.ReceiverIDLogKey, rid), zap.String("tracker_id", tid))
		return response.CreateInternalServerErrorResponse(), nil
	}
	if hasEvents {
		return response.FormatResponse(
			response.ErrorResponse{Status: "Conflict", DeveloperText: "This tracker has associated events. Deactivate it instead of deleting it."},
			http.StatusConflict,
		), nil
	}

	if err := params.TrackerRepo.DeleteTracker(rid, tid); err != nil {
		params.AppCfg.Logger.Error("error deleting tracker", zap.String(log.ReceiverIDLogKey, rid), zap.String(log.UserIDLogKey, uid))
		return response.CreateInternalServerErrorResponse(), nil
	}

	params.AppCfg.Logger.Info(handlerSuccessful, zap.String("action", deleteTracker), zap.String(log.UserIDLogKey, uid), zap.String(log.ReceiverIDLogKey, rid), zap.String("tracker_id", tid))
	return response.FormatResponse(map[string]string{"status": response.Success}, http.StatusOK), nil
}
