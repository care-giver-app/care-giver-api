package repository

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
	"github.com/care-giver-app/care-giver-api/internal/log"
	"go.uber.org/zap"
)

type DynamodbClientProvider interface {
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

type UserRepositoryContextKey string

const (
	userRepoContextKey UserRepositoryContextKey = "UserRepository"
)

type UserRepository struct {
	Ctx       context.Context
	Client    DynamodbClientProvider
	TableName string
	logger    *zap.Logger
}

func NewUserRespository(ctx context.Context, cfg *appconfig.AppConfig, client DynamodbClientProvider) *UserRepository {
	return &UserRepository{
		Ctx:       ctx,
		Client:    client,
		TableName: cfg.UserTableName,
		logger:    cfg.Logger.With(zap.String(log.TableNameLogKey, cfg.UserTableName)),
	}
}

func (ur *UserRepository) PutUser() error {
	return nil
}

func (ur *UserRepository) GetReceivers() error {
	return nil
}

func (ur *UserRepository) CreatePrimaryReceiver() error {
	return nil
}

func (ur *UserRepository) AddAdditionalReceiver() error {
	return nil
}

func ContextWithUserRespository(ctx context.Context, ur *UserRepository) context.Context {
	return context.WithValue(ctx, userRepoContextKey, ur)
}

func UserRepositoryFromContext(ctx context.Context) (*UserRepository, error) {
	ur, ok := ctx.Value(userRepoContextKey).(*UserRepository)
	if !ok {
		return nil, errors.New("unable to get user repostiory from context")
	}
	return ur, nil
}
