package handlers

import (
	"errors"

	"github.com/care-giver-app/care-giver-golang-common/pkg/event"
	"github.com/care-giver-app/care-giver-golang-common/pkg/receiver"
	"github.com/care-giver-app/care-giver-golang-common/pkg/relationship"
	"github.com/care-giver-app/care-giver-golang-common/pkg/user"
)

var (
	testUserRepo         = &MockUserRepo{}
	testReceiverRepo     = &MockReceiverRepo{}
	testEventRepo        = &MockEventRepo{}
	testRelationshipRepo = &MockRelationshipRepo{}
)

type MockUserRepo struct{}

func (mu *MockUserRepo) CreateUser(u user.User) error {
	switch u.Email {
	case "good@test.com":
		return nil
	case "error@test.com":
		return errors.New("error creating user")
	}
	return errors.New("unsupported mock")
}

func (mu *MockUserRepo) GetUser(uid string) (user.User, error) {
	switch uid {
	case "User#123":
		return user.User{
			UserID: "User#123",
		}, nil
	case "User#RelationshipError":
		return user.User{
			UserID: "User#RelationshipError",
		}, nil
	case "User#NotACareGiver":
		return user.User{
			UserID: "User#NotACareGiver",
		}, nil
	case "User#Error":
		return user.User{}, errors.New("error getting user from db")
	}
	return user.User{}, errors.New("unsupported mock")
}

func (mu *MockUserRepo) GetUserByEmail(email string) (user.User, error) {
	switch email {
	case "valid@example.com":
		return user.User{
			UserID: "User#123",
		}, nil
	case "error@example.com":
		return user.User{}, errors.New("error getting user from db")
	case "relationshiperror@example.com":
		return user.User{
			UserID: "User#RelationshipError",
		}, nil
	}
	return user.User{}, errors.New("unsupported mock")
}

func (mu *MockUserRepo) UpdateReceiverList(uid string, rid string, listName string) error {
	switch uid {
	case "User#123":
		return nil
	case "User#RelationshipError":
		return errors.New("error updating relationship")
	}
	return errors.New("unsupported mock")
}

type MockReceiverRepo struct{}

func (mr *MockReceiverRepo) CreateReceiver(r receiver.Receiver) error {
	switch r.FirstName {
	case "Good":
		return nil
	case "Error":
		return errors.New("error creating receiver")
	}
	return errors.New("unsupported mock")
}

func (mr *MockReceiverRepo) GetReceiver(rid string) (receiver.Receiver, error) {
	switch rid {
	case "Receiver#123":
		return receiver.Receiver{
			FirstName: "Success",
		}, nil
	case "Receiver#Error":
		return receiver.Receiver{}, errors.New("error retrieving from db")
	}
	return receiver.Receiver{}, errors.New("unsupported mock")
}

type MockEventRepo struct{}

func (me *MockEventRepo) AddEvent(e *event.Entry) error {
	switch e.ReceiverID {
	case "Receiver#123":
		return nil
	case "Receiver#Error":
		return errors.New("error adding event")
	}
	return errors.New("unsupported mock")
}
func (me *MockEventRepo) GetEvents(rid string) ([]event.Entry, error) {
	switch rid {
	case "Receiver#123":
		return []event.Entry{
			{
				EventID:    "Event#123",
				ReceiverID: "Receiver#123",
			},
		}, nil
	case "Receiver#Error":
		return nil, errors.New("error retrieving events")
	}
	return nil, errors.New("unsupported mock")
}

func (me *MockEventRepo) DeleteEvent(rid, eid string) error {
	switch rid {
	case "Receiver#123":
		if eid == "Event#123" {
			return nil
		}
		return errors.New("event not found")
	case "Receiver#Error":
		return errors.New("error deleting event")
	}
	return errors.New("unsupported mock")
}

type MockRelationshipRepo struct{}

func (mr *MockRelationshipRepo) GetRelationshipsByUser(uid string) ([]relationship.Relationship, error) {
	switch uid {
	case "User#123":
		return []relationship.Relationship{
			{
				UserID:             "User#123",
				ReceiverID:         "Receiver#123",
				PrimaryCareGiver:   true,
				EmailNotifications: true,
			},
			{
				UserID:             "User#123",
				ReceiverID:         "Receiver#Error",
				PrimaryCareGiver:   false,
				EmailNotifications: false,
			},
		}, nil
	case "User#NotACareGiver":
		return []relationship.Relationship{}, nil
	case "User#NotAPrimaryCareGiver":
		return []relationship.Relationship{
			{
				UserID:             "User#NotAPrimaryCareGiver",
				ReceiverID:         "Receiver#123",
				PrimaryCareGiver:   false,
				EmailNotifications: true,
			},
		}, nil
	case "User#RelationshipError":
		return nil, errors.New("error retrieving relationships from db")
	}
	return nil, errors.New("unsupported mock")
}

func (mr *MockRelationshipRepo) AddRelationship(r *relationship.Relationship) error {
	switch r.UserID {
	case "User#123":
		return nil
	case "User#RelationshipError":
		return errors.New("error adding relationship")
	}
	return nil
}

func (mr *MockRelationshipRepo) DeleteRelationship(uid, rid string) error {
	return nil
}

func (mr *MockRelationshipRepo) GetRelationship(userID string, receiverID string) (*relationship.Relationship, error) {
	return nil, errors.New("unsupported mock")
}

func (mr *MockRelationshipRepo) GetRelationshipsByEmailNotifications() ([]relationship.Relationship, error) {
	return nil, errors.New("unsupported mock")
}
