package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/care-giver-app/care-giver-api/internal/appconfig"
	"github.com/care-giver-app/care-giver-api/internal/user"
	"github.com/stretchr/testify/assert"
)

type MockUserRepository struct{}

func (m *MockUserRepository) PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
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

func (m *MockUserRepository) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if av, found := params.Key["user_id"]; found {
		if id, ok := av.(*types.AttributeValueMemberS); ok {
			switch id.Value {
			case "User#123":
				return &dynamodb.GetItemOutput{
					Item: map[string]types.AttributeValue{
						"user_id":                &types.AttributeValueMemberS{Value: id.Value},
						"first_name":             &types.AttributeValueMemberS{Value: "testFirstName"},
						"last_name":              &types.AttributeValueMemberS{Value: "testLastName"},
						"primary_care_receivers": &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
					},
				}, nil
			case "Get Item Error":
				return nil, errors.New("An error occured during Get Item")
			case "Unmarshal Error":
				return &dynamodb.GetItemOutput{
					Item: map[string]types.AttributeValue{
						"user_id":                &types.AttributeValueMemberS{Value: id.Value},
						"first_name":             &types.AttributeValueMemberS{Value: "testFirstName"},
						"last_name":              &types.AttributeValueMemberS{Value: "testLastName"},
						"primary_care_receivers": &types.AttributeValueMemberBOOL{Value: false},
					},
				}, nil
			}
		}
	}
	return nil, errors.New("unsupported mock")
}

func (m *MockUserRepository) UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	if av, found := params.Key["user_id"]; found {
		if id, ok := av.(*types.AttributeValueMemberS); ok {
			switch id.Value {
			case "User#123":
				return &dynamodb.UpdateItemOutput{}, nil
			case "Update Item Error":
				return nil, errors.New("An error occured during Update Item")
			}
		}
	}
	return nil, errors.New("unsupported mock")
}

func (m *MockUserRepository) Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	return nil, errors.New("unsupported mock")
}

func (m *MockUserRepository) DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error) {
	if av, found := params.Key["user_id"]; found {
		if id, ok := av.(*types.AttributeValueMemberS); ok {
			switch id.Value {
			case "User#123":
				return &dynamodb.DeleteItemOutput{}, nil
			case "Delete Item Error":
				return nil, errors.New("An error occured during Delete Item")
			}
		}
	}
	return nil, errors.New("unsupported mock")
}

func TestCreateUser(t *testing.T) {
	appCfg := appconfig.NewAppConfig()
	testUserRepo := NewUserRespository(context.Background(), appCfg, &MockUserRepository{})

	tests := map[string]struct {
		user        user.User
		expectError bool
	}{
		"Happy Path - User Created": {
			user: user.User{
				UserID:    "User#123",
				FirstName: "testName",
				LastName:  "testLastName",
			},
		},
		"Sad Path - Error Putting Item": {
			user: user.User{
				UserID: "Error",
			},
			expectError: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := testUserRepo.CreateUser(tc.user)

			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	appCfg := appconfig.NewAppConfig()
	testUserRepo := NewUserRespository(context.Background(), appCfg, &MockUserRepository{})

	tests := map[string]struct {
		userID       string
		expectedUser user.User
		expectError  bool
	}{
		"Happy Path - Got User": {
			userID: "User#123",
			expectedUser: user.User{
				UserID:               "User#123",
				FirstName:            "testFirstName",
				LastName:             "testLastName",
				PrimaryCareReceivers: []string{},
			},
		},
		"Sad Path - Error Getting Item": {
			userID:      "Get Item Error",
			expectError: true,
		},
		"Sad Path - Error Unmarshalling Item": {
			userID:      "Unmarshal Error",
			expectError: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			user, err := testUserRepo.GetUser(tc.userID)

			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedUser, user)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	appCfg := appconfig.NewAppConfig()
	testUserRepo := NewUserRespository(context.Background(), appCfg, &MockUserRepository{})

	tests := map[string]struct {
		userID      string
		expectError bool
	}{
		"Happy Path - User Updated": {
			userID: "User#123",
		},
		"Sad Path - Error Updating Item": {
			userID:      "Update Item Error",
			expectError: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := testUserRepo.UpdateReceiverList(tc.userID, "", "")

			if tc.expectError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}
