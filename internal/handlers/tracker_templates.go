package handlers

import (
	"context"
	"net/http"
	"net/url"

	awsevents "github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/response"
	"github.com/care-giver-app/care-giver-golang-common/pkg/tracker"
	"go.uber.org/zap"
)

const (
	listTrackerTemplates = "list tracker templates"
	getTrackerTemplate   = "get tracker template"
)

func HandleListTrackerTemplates(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, listTrackerTemplates)

	uid, ok := extractUID(params)
	if !ok {
		return response.CreateAccessDeniedResponse(), nil
	}

	templates, err := tracker.GetAllTemplates()
	if err != nil {
		params.AppCfg.Logger.Error("error loading tracker templates", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}

	params.AppCfg.Logger.Sugar().Infof(handlerSuccessful, listTrackerTemplates)
	_ = uid
	return response.FormatResponse(templates, http.StatusOK), nil
}

func HandleGetTrackerTemplate(ctx context.Context, params HandlerParams) (awsevents.APIGatewayProxyResponse, error) {
	params.AppCfg.Logger.Sugar().Infof(handlerStart, getTrackerTemplate)

	uid, ok := extractUID(params)
	if !ok {
		return response.CreateAccessDeniedResponse(), nil
	}

	rawName, found := params.Request.PathParameters["templateName"]
	if !found || rawName == "" {
		params.AppCfg.Logger.Error("missing templateName path parameter")
		return response.CreateBadRequestResponse(), nil
	}

	templateName, err := url.PathUnescape(rawName)
	if err != nil {
		params.AppCfg.Logger.Error("failed to decode templateName", zap.Error(err))
		return response.CreateBadRequestResponse(), nil
	}

	tmpl, err := tracker.GetTemplateByName(templateName)
	if err != nil {
		params.AppCfg.Logger.Error("error loading tracker template", zap.Error(err))
		return response.CreateInternalServerErrorResponse(), nil
	}
	if tmpl == nil {
		return response.CreateResourceNotFoundResponse(), nil
	}

	params.AppCfg.Logger.Sugar().Infof(handlerSuccessful, getTrackerTemplate)
	_ = uid
	return response.FormatResponse(tmpl, http.StatusOK), nil
}
