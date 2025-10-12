package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
	"github.com/care-giver-app/care-giver-api/internal/relationship"
	"github.com/stretchr/testify/assert"
)

type MockRelationshipDB struct{}

func (m *MockRelationshipDB) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if av, found := params.Item["user_id"]; found {
		if id, ok := av.(*types.AttributeValueMemberS); ok {
			switch id.Value {
			case "User#123":
				return nil, nil
			case "Error":
				return nil, errors.New("An error occured during Put Item")
			}
		}
	}
	return nil, errors.New("unsupported mock")
}

func (m *MockRelationshipDB) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if av, found := params.Key["user_id"]; found {
		if id, ok := av.(*types.AttributeValueMemberS); ok {
			switch id.Value {
			case "User#123":
				return &dynamodb.GetItemOutput{
					Item: map[string]types.AttributeValue{
						"user_id":             &types.AttributeValueMemberS{Value: id.Value},
						"receiver_id":         &types.AttributeValueMemberS{Value: "Receiver#123"},
						"primary_care_giver":  &types.AttributeValueMemberBOOL{Value: true},
						"email_notifications": &types.AttributeValueMemberBOOL{Value: false},
					},
				}, nil
			case "Get Item Error":
				return nil, errors.New("An error occured during Get Item")
			case "Unmarshal Error":
				return &dynamodb.GetItemOutput{
					Item: map[string]types.AttributeValue{
						"user_id":            &types.AttributeValueMemberS{Value: id.Value},
						"receiver_id":        &types.AttributeValueMemberS{Value: "testFirstName"},
						"primary_care_giver": &types.AttributeValueMemberS{Value: "false"},
					},
				}, nil
			}
		}
	}
	return nil, errors.New("unsupported mock")
}

func (m *MockRelationshipDB) UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	return nil, errors.New("unsupported mock")
}

func (m *MockRelationshipDB) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	if av, found := params.ExpressionAttributeValues[":uid"]; found {
		if uid, ok := av.(*types.AttributeValueMemberS); ok {
			switch uid.Value {
			case "User#123":
				return &dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{
						{
							"user_id":             &types.AttributeValueMemberS{Value: "User#123"},
							"receiver_id":         &types.AttributeValueMemberS{Value: "Receiver#123"},
							"primary_care_giver":  &types.AttributeValueMemberBOOL{Value: true},
							"email_notifications": &types.AttributeValueMemberBOOL{Value: false},
						},
					},
				}, nil
			case "Error":
				return nil, errors.New("An error occured during Query")
			case "Unmarshal Error":
				return &dynamodb.QueryOutput{
					Items: []map[string]types.AttributeValue{
						{
							"user_id": &types.AttributeValueMemberBOOL{Value: false},
						},
					},
				}, nil
			}
		}
	}

	return nil, errors.New("unsupported mock")
}

func (m *MockRelationshipDB) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	if av, found := params.Key["user_id"]; found {
		if id, ok := av.(*types.AttributeValueMemberS); ok {
			switch id.Value {
			case "User#123":
				return &dynamodb.DeleteItemOutput{}, nil
			case "Error":
				return nil, errors.New("An error occured during Delete Item")
			}
		}
	}
	return nil, errors.New("unsupported mock")
}

func TestAddRelationship(t *testing.T) {
	appCfg := appconfig.NewAppConfig()
	testEventRepo := NewRelationshipRepository(context.Background(), appCfg, &MockRelationshipDB{})

	tests := map[string]struct {
		relationship *relationship.Relationship
		expectError  bool
	}{
		"Happy Path - Event Added": {
			relationship: &relationship.Relationship{
				UserID:             "User#123",
				ReceiverID:         "Receiver#123",
				PrimaryCareGiver:   true,
				EmailNotifications: false,
			},
		},
		"Sad Path - Put Item Error": {
			relationship: &relationship.Relationship{
				UserID:             "Error",
				ReceiverID:         "Receiver#123",
				PrimaryCareGiver:   true,
				EmailNotifications: false,
			},
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := testEventRepo.AddRelationship(tc.relationship)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetRelationship(t *testing.T) {
	appCfg := appconfig.NewAppConfig()
	testEventRepo := NewRelationshipRepository(context.Background(), appCfg, &MockRelationshipDB{})

	tests := map[string]struct {
		userID               string
		receiverID           string
		expectedRelationship relationship.Relationship
		expectError          bool
	}{
		"Happy Path": {
			userID:     "User#123",
			receiverID: "Receiver#123",
			expectedRelationship: relationship.Relationship{
				UserID:             "User#123",
				ReceiverID:         "Receiver#123",
				PrimaryCareGiver:   true,
				EmailNotifications: false,
			},
		},
		"Sad Path - Get Item Error": {
			userID:      "Error",
			receiverID:  "Receiver#123",
			expectError: true,
		},
		"Sad Path - Unmarshal Error": {
			userID:      "Unmarshal Error",
			receiverID:  "Receiver#123",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := testEventRepo.GetRelationship(tc.userID, tc.receiverID)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, r)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, r)
				assert.Equal(t, tc.expectedRelationship, *r)
			}
		})
	}
}

func TestGetRelationships(t *testing.T) {
	appCfg := appconfig.NewAppConfig()
	testEventRepo := NewRelationshipRepository(context.Background(), appCfg, &MockRelationshipDB{})

	tests := map[string]struct {
		userID        string
		expectedValue []relationship.Relationship
		expectError   bool
	}{
		"Happy Path": {
			userID: "User#123",
			expectedValue: []relationship.Relationship{
				{
					UserID:             "User#123",
					ReceiverID:         "Receiver#123",
					PrimaryCareGiver:   true,
					EmailNotifications: false,
				},
			},
		},
		"Sad Path - Get Item Error": {
			userID:      "Error",
			expectError: true,
		},
		"Sad Path - Unmarshal Error": {
			userID:      "Unmarshal Error",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			r, err := testEventRepo.GetRelationshipsByUser(tc.userID)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, r)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, r)
				assert.Equal(t, tc.expectedValue, r)
			}
		})
	}
}

func TestDeleteRelationship(t *testing.T) {
	appCfg := appconfig.NewAppConfig()
	testEventRepo := NewRelationshipRepository(context.Background(), appCfg, &MockRelationshipDB{})

	tests := map[string]struct {
		userID      string
		receiverID  string
		expectError bool
	}{
		"Happy Path": {
			userID:     "User#123",
			receiverID: "Receiver#123",
		},
		"Sad Path": {
			userID:      "Error",
			receiverID:  "Receiver#123",
			expectError: true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := testEventRepo.DeleteRelationship(tc.userID, tc.receiverID)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
