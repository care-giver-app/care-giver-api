package handlers

import (
	"context"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
	"github.com/care-giver-app/care-giver-api/internal/repository"
)

const (
	handlerStart          = "handling %s"
	handlerSuccessful     = "processed %s successfully"
	requestBodyError      = "error reading request body"
	pathParametersError   = "error validating path parameters"
	queryParamsError      = "error validating query parameters"
	userDatbaseError      = "error retrieving user from db"
	receiverDatabaseError = "error retrieving receiver from db"
	userNotCareGiverError = "user is not a caregiver for the receiver"
)

type HandlerParams struct {
	AppCfg       *appconfig.AppConfig
	Request      events.APIGatewayProxyRequest
	UserRepo     repository.UserRepositoryProvider
	ReceiverRepo repository.ReceiverRepositoryProvider
	EventRepo    repository.EventRepositoryProvider
}

type Endpoint struct {
	Path   string
	Method string
}

type HandlerFunc func(ctx context.Context, params HandlerParams) (events.APIGatewayProxyResponse, error)

var handlersMap = map[Endpoint]HandlerFunc{
	Endpoint{"/user", http.MethodPost}:                     HandleCreateUser,
	Endpoint{"/user/{userId}", http.MethodGet}:             HandleGetUser,
	Endpoint{"/user/primary-receiver", http.MethodPost}:    HandleUserPrimaryReceiver,
	Endpoint{"/user/additional-receiver", http.MethodPost}: HandleUserAdditionalReceiver,
	Endpoint{"/receiver/{receiverId}", http.MethodGet}:     HandleReceiver,
	Endpoint{"/event", http.MethodPost}:                    HandleReceiverEvent,
	Endpoint{"/event/{eventId}", http.MethodDelete}:        HandleDeleteReceiverEvent,
	Endpoint{"/events/{receiverId}", http.MethodGet}:       HandleGetReceiverEvents,
}

type RegistryProvider interface {
	GetHandler(request events.APIGatewayProxyRequest) (HandlerFunc, bool)
	RunHandler(ctx context.Context, handler HandlerFunc, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type Registry struct {
	AppCfg       *appconfig.AppConfig
	UserRepo     repository.UserRepositoryProvider
	ReceiverRepo repository.ReceiverRepositoryProvider
	EventRepo    repository.EventRepositoryProvider
}

func NewRegistry(appCfg *appconfig.AppConfig, userRepo repository.UserRepositoryProvider, receiverRepo repository.ReceiverRepositoryProvider, eventRepo repository.EventRepositoryProvider) *Registry {
	return &Registry{
		AppCfg:       appCfg,
		UserRepo:     userRepo,
		ReceiverRepo: receiverRepo,
		EventRepo:    eventRepo,
	}
}

func (r *Registry) GetHandler(request events.APIGatewayProxyRequest) (HandlerFunc, bool) {
	endpoint := Endpoint{
		Path:   removePathPrefix(request.RequestContext.ResourcePath),
		Method: request.HTTPMethod,
	}

	handler, exists := handlersMap[endpoint]
	return handler, exists
}

func (r *Registry) RunHandler(ctx context.Context, handler HandlerFunc, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	params := HandlerParams{
		AppCfg:       r.AppCfg,
		Request:      request,
		UserRepo:     r.UserRepo,
		ReceiverRepo: r.ReceiverRepo,
		EventRepo:    r.EventRepo,
	}

	return handler(ctx, params)
}

func removePathPrefix(path string) string {
	pathPrefixes := []string{"/Stage", "/Prod"}
	for _, prefix := range pathPrefixes {
		path = strings.TrimPrefix(path, prefix)
	}
	return path
}
