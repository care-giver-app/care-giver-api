package handlers

import (
	"context"
	"net/http"

	awsevents "github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/response"
	"github.com/care-giver-app/care-giver-golang-common/pkg/log"
	"github.com/care-giver-app/care-giver-golang-common/pkg/receiver"
	"github.com/care-giver-app/care-giver-golang-common/pkg/relationship"
	"github.com/care-giver-app/care-giver-golang-common/pkg/user"
	"go.uber.org/zap"
)

const (
	getReceiver = "get receiver"
)

func HandleReceiver(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, getReceiver)

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

	relationships, err := params.RelationshipRepo.GetRelationshipsByUser(uid)
	if err != nil {
		params.AppCfg.Logger.Error(relationshipDatabaseError, zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	if !relationship.IsACareGiver(uid, rid, relationships) {
		params.AppCfg.Logger.Error(userNotCareGiverError, zap.String(log.ReceiverIDLogKey, rid), zap.String(log.UserIDLogKey, uid))
		return response.CreateAccessDeniedResponse(), nil
	}

	r, err := params.ReceiverRepo.GetReceiver(rid)
	if err != nil {
		params.AppCfg.Logger.Error(receiverDatabaseError, zap.String(log.ReceiverIDLogKey, rid), zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	params.AppCfg.Logger.Sugar().Infof(handlerSuccessful, getReceiver)
	return response.FormatResponse(r, http.StatusOK), nil
}

type CareGiverResponse struct {
	UserID    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	IsPrimary bool   `json:"isPrimary"`
}

type GetReceiverCareGiversResponse struct {
	CareGivers []CareGiverResponse `json:"careGivers"`
}

func HandleGetReceiverCareGivers(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, "get receiver care givers")

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

	userRelationships, err := params.RelationshipRepo.GetRelationshipsByUser(uid)
	if err != nil {
		params.AppCfg.Logger.Error(relationshipDatabaseError, zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	if !relationship.IsACareGiver(uid, rid, userRelationships) {
		params.AppCfg.Logger.Error(userNotCareGiverError, zap.String(log.ReceiverIDLogKey, rid), zap.String(log.UserIDLogKey, uid))
		return response.CreateAccessDeniedResponse(), nil
	}

	receiverRelationships, err := params.RelationshipRepo.GetRelationshipsByReceiver(rid)
	if err != nil {
		params.AppCfg.Logger.Error(relationshipDatabaseError, zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	careGivers := make([]CareGiverResponse, 0, len(receiverRelationships))
	for _, rel := range receiverRelationships {
		u, err := params.UserRepo.GetUser(rel.UserID)
		if err != nil {
			params.AppCfg.Logger.Error(userDatbaseError, zap.String(log.UserIDLogKey, rel.UserID), zap.Error(err))
			return response.CreateInternalServerErrorResponse(), nil
		}
		careGivers = append(careGivers, CareGiverResponse{
			UserID:    u.UserID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			IsPrimary: rel.PrimaryCareGiver,
		})
	}

	params.AppCfg.Logger.Sugar().Infof(handlerSuccessful, "get receiver care givers")
	return response.FormatResponse(GetReceiverCareGiversResponse{
		CareGivers: careGivers,
	}, http.StatusOK), nil
}
