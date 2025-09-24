package repository

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
	"github.com/care-giver-app/care-giver-api/internal/log"
	"github.com/care-giver-app/care-giver-api/internal/relationship"
	"go.uber.org/zap"
)

type RelationshipRepositoryProvider interface {
	AddRelationship(r *relationship.Relationship) error
	GetRelationship(userID string, receiverID string) (*relationship.Relationship, error)
	GetRelationshipsByUser(userID string) ([]relationship.Relationship, error)
	DeleteRelationship(userID string, receiverID string) error
}

type RelationshipRepository struct {
	Ctx       context.Context
	Client    DynamodbClientProvider
	TableName string
	logger    *zap.Logger
}

func NewRelationshipRepository(ctx context.Context, cfg *appconfig.AppConfig, client DynamodbClientProvider) *RelationshipRepository {
	return &RelationshipRepository{
		Ctx:       ctx,
		Client:    client,
		TableName: cfg.RelationshipTableName,
		logger:    cfg.Logger.With(zap.String(log.TableNameLogKey, cfg.RelationshipTableName)),
	}
}

func (rr *RelationshipRepository) AddRelationship(r *relationship.Relationship) error {
	rr.logger.Info("adding user receiver relationship to db")

	rr.logger.Info("marshalling user receiver relationship struct")
	av, err := attributevalue.MarshalMap(r)
	if err != nil {
		return err
	}

	rr.logger.Info("inserting item into db", zap.Any("item", av))
	_, err = rr.Client.PutItem(rr.Ctx, &dynamodb.PutItemInput{
		TableName: aws.String(rr.TableName),
		Item:      av,
	})
	if err != nil {
		return err
	}
	rr.logger.Info("successfully inserted item")

	return nil
}

func (rr *RelationshipRepository) GetRelationship(userID string, receiverID string) (*relationship.Relationship, error) {
	rr.logger.Info("getting user receiver relationship from db", zap.String(log.UserIDLogKey, userID), zap.String(log.ReceiverIDLogKey, receiverID))

	result, err := rr.Client.GetItem(rr.Ctx, &dynamodb.GetItemInput{
		TableName: &rr.TableName,
		Key: map[string]types.AttributeValue{
			"user_id":     &types.AttributeValueMemberS{Value: userID},
			"receiver_id": &types.AttributeValueMemberS{Value: receiverID},
		},
	})
	if err != nil {
		return nil, err
	}

	var r relationship.Relationship
	err = attributevalue.UnmarshalMap(result.Item, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (rr *RelationshipRepository) GetRelationshipsByUser(userID string) ([]relationship.Relationship, error) {
	rr.logger.Info("getting user receiver relationships from db", zap.String(log.UserIDLogKey, userID))

	keyCondition := "user_id = :uid"
	expressionAttributeValues := map[string]types.AttributeValue{
		":uid": &types.AttributeValueMemberS{Value: userID},
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(rr.TableName),
		KeyConditionExpression:    aws.String(keyCondition),
		ExpressionAttributeValues: expressionAttributeValues,
	}

	result, err := rr.Client.Query(rr.Ctx, queryInput)
	if err != nil {
		return nil, err
	}

	var relationshipsList []relationship.Relationship
	err = attributevalue.UnmarshalListOfMaps(result.Items, &relationshipsList)
	if err != nil {
		rr.logger.Error("error unmarshalling relationships list", zap.Error(err))
		return nil, err
	}

	return relationshipsList, nil
}

func (rr *RelationshipRepository) DeleteRelationship(userID string, receiverID string) error {
	rr.logger.Info("deleting user receiver relationship from db", zap.String(log.UserIDLogKey, userID), zap.String(log.ReceiverIDLogKey, receiverID))

	_, err := rr.Client.DeleteItem(rr.Ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(rr.TableName),
		Key: map[string]types.AttributeValue{
			"user_id":     &types.AttributeValueMemberS{Value: userID},
			"receiver_id": &types.AttributeValueMemberS{Value: receiverID},
		},
	})

	if err != nil {
		return err
	}

	rr.logger.Info("successfully deleted relationship")
	return nil
}
