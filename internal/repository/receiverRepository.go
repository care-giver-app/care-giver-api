package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
	"github.com/care-giver-app/care-giver-api/internal/log"
	"github.com/care-giver-app/care-giver-api/internal/receiver"
	"github.com/care-giver-app/care-giver-api/internal/receiver/events"
	"go.uber.org/zap"
)

var (
	updateEventExpression = "SET #evt = list_append(#evt, :val)"
	receiverID            = "receiver_id"
)

type ReceiverRepositoryProvider interface {
	CreateReceiver(r receiver.Receiver) error
	GetReceiver(rid receiver.ReceiverID) (receiver.Receiver, error)
	UpdateReceiver(rid string, evt events.Event, eventName string) error
}

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
		TableName: cfg.ReceiverTableName,
		logger:    cfg.Logger.With(zap.String(log.TableNameLogKey, cfg.ReceiverTableName)),
	}
}

func (rr *ReceiverRepository) CreateReceiver(r receiver.Receiver) error {
	rr.logger.Info("adding receiver to db", zap.Any(log.ReceiverIDLogKey, r.ReceiverID))

	rr.logger.Info("marshalling receiver struct")
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

func (rr *ReceiverRepository) GetReceiver(rid receiver.ReceiverID) (receiver.Receiver, error) {
	rr.logger.Info("getting receiver from db", zap.Any(log.ReceiverIDLogKey, rid))
	result, err := rr.Client.GetItem(rr.Ctx, &dynamodb.GetItemInput{
		TableName: &rr.TableName,
		Key: map[string]types.AttributeValue{
			receiverID: &types.AttributeValueMemberS{Value: string(rid)},
		},
	})
	if err != nil {
		return receiver.Receiver{}, err
	}

	var r receiver.Receiver
	err = attributevalue.UnmarshalMap(result.Item, &r)
	if err != nil {
		return receiver.Receiver{}, err
	}

	return r, nil
}

func (rr *ReceiverRepository) UpdateReceiver(rid string, evt events.Event, eventName string) error {
	rr.logger.Info("updating receiver in db", zap.Any(log.ReceiverIDLogKey, rid))

	rr.logger.Info("marshalling receiver struct")
	av, err := attributevalue.MarshalMap(evt)
	if err != nil {
		return err
	}

	rr.logger.Info("updating item in db", zap.Any("item", av))
	_, err = rr.Client.UpdateItem(rr.Ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(rr.TableName),
		Key: map[string]types.AttributeValue{
			receiverID: &types.AttributeValueMemberS{Value: string(rid)},
		},
		UpdateExpression: &updateEventExpression,
		ExpressionAttributeNames: map[string]string{
			"#evt": eventName,
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":val": &types.AttributeValueMemberL{Value: []types.AttributeValue{&types.AttributeValueMemberM{Value: av}}},
		},
		ConditionExpression: aws.String(fmt.Sprintf("attribute_exists(%s)", receiverID)),
	})
	if err != nil {
		return err
	}
	rr.logger.Info("successfully updated item")

	return nil
}
