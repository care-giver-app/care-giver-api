package repository

import (
	"context"
	"errors"

	"github.com/care-giver-app/care-giver-api/internal/appconfig"
	"github.com/care-giver-app/care-giver-api/internal/log"
	"go.uber.org/zap"
)

type ReceiverRepositoryContextKey string

const (
	receiverRepoContextKey ReceiverRepositoryContextKey = "ReceiverRepository"
)

type ReceiverRepository struct {
	Ctx       context.Context
	Client    DynamodbClientProvider
	TableName string
	logger    *zap.Logger
}

func NewReceiverRespository(ctx context.Context, cfg *appconfig.AppConfig, client DynamodbClientProvider) *ReceiverRepository {
	return &ReceiverRepository{
		Ctx:       ctx,
		Client:    client,
		TableName: cfg.UserTableName,
		logger:    cfg.Logger.With(zap.String(log.TableNameLogKey, cfg.UserTableName)),
	}
}

func (rr *ReceiverRepository) GetReceiver() error {
	return nil
}

func (rr *ReceiverRepository) AddEvent() error {
	return nil
}

func ContextWithReceiverRespository(ctx context.Context, rr *ReceiverRepository) context.Context {
	return context.WithValue(ctx, receiverRepoContextKey, rr)
}

func ReceiverRepositoryFromContext(ctx context.Context) (*ReceiverRepository, error) {
	rr, ok := ctx.Value(receiverRepoContextKey).(*ReceiverRepository)
	if !ok {
		return nil, errors.New("unable to get user repostiory from context")
	}
	return rr, nil
}
