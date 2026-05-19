package handlers

import (
	"errors"
	"fmt"
	"strings"

	"github.com/care-giver-app/care-giver-golang-common/pkg/event"
	"github.com/care-giver-app/care-giver-golang-common/pkg/receiver"
	"github.com/care-giver-app/care-giver-golang-common/pkg/relationship"
	"github.com/care-giver-app/care-giver-golang-common/pkg/repository"
	pkgtracker "github.com/care-giver-app/care-giver-golang-common/pkg/tracker"
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
			UserID:    "User#123",
			FirstName: "John",
			LastName:  "Doe",
		}, nil
	case "User#RelationshipError":
		return user.User{
			UserID: "User#RelationshipError",
		}, nil
	case "User#NotACareGiver":
		return user.User{
			UserID: "User#NotACareGiver",
		}, nil
	case "User#456":
		return user.User{
			UserID:    "User#456",
			FirstName: "Jane",
			LastName:  "Smith",
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
func (me *MockEventRepo) GetEvents(rid string, bound repository.TimestampBound) ([]event.Entry, error) {
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

func (me *MockEventRepo) HasEventsForTracker(receiverID, trackerID string) (bool, error) {
	switch trackerID {
	case "Tracker#123":
		return true, nil
	case "Tracker#NoEvents":
		return false, nil
	case "Tracker#EventCheckError":
		return false, errors.New("error checking events for tracker")
	case "Tracker#NotFound":
		return false, nil
	}
	return false, nil
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
				ReceiverID:         "Receiver#Empty",
				PrimaryCareGiver:   true,
				EmailNotifications: false,
			},
			{
				UserID:             "User#123",
				ReceiverID:         "Receiver#Error",
				PrimaryCareGiver:   false,
				EmailNotifications: false,
			},
			{
				UserID:             "User#123",
				ReceiverID:         "Receiver#RelationshipError",
				PrimaryCareGiver:   true,
				EmailNotifications: false,
			},
			{
				UserID:             "User#123",
				ReceiverID:         "Receiver#UserError",
				PrimaryCareGiver:   true,
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

var testTrackerRepo = &MockTrackerRepo{}

type MockTrackerRepo struct{}

func (m *MockTrackerRepo) CreateTracker(t *pkgtracker.Tracker) error {
	switch t.ReceiverID {
	case "Receiver#123":
		if strings.EqualFold(t.Name, "Duplicate Name") {
			return fmt.Errorf("name already in use for this receiver")
		}
		return nil
	case "Receiver#Error":
		return errors.New("error creating tracker")
	}
	return errors.New("unsupported mock")
}

func (m *MockTrackerRepo) GetTracker(receiverID, trackerID string) (*pkgtracker.Tracker, error) {
	switch trackerID {
	case "Tracker#123":
		return &pkgtracker.Tracker{
			TrackerID:  "Tracker#123",
			ReceiverID: receiverID,
			Name:       "Blood Pressure",
			Kind:       pkgtracker.KindMeasurement,
			Fields:     []pkgtracker.TrackerField{{Name: "systolic", InputType: "number", Required: true}},
			Icon:       "assets/icon.svg",
			Color:      pkgtracker.ColorConfig{Primary: "#000", Secondary: "#fff"},
			IsActive:   true,
			CreatedAt:  "2026-01-01T00:00:00Z",
			UpdatedAt:  "2026-01-01T00:00:00Z",
		}, nil
	case "Tracker#456":
		return &pkgtracker.Tracker{
			TrackerID:  "Tracker#456",
			ReceiverID: receiverID,
			Name:       "Shower",
			Kind:       pkgtracker.KindEvent,
			Fields:     []pkgtracker.TrackerField{},
			Icon:       "assets/shower-icon.svg",
			Color:      pkgtracker.ColorConfig{Primary: "#1E90FF", Secondary: "#D1E8FF"},
			IsActive:   true,
			CreatedAt:  "2026-01-01T00:00:00Z",
			UpdatedAt:  "2026-01-01T00:00:00Z",
		}, nil
	case "Tracker#EventCheckError":
		return &pkgtracker.Tracker{
			TrackerID:  "Tracker#EventCheckError",
			ReceiverID: receiverID,
			Name:       "Walk",
			Kind:       pkgtracker.KindEvent,
			Fields:     []pkgtracker.TrackerField{},
			Icon:       "assets/walk-icon.svg",
			Color:      pkgtracker.ColorConfig{Primary: "#990000", Secondary: "#ff6666"},
			IsActive:   true,
			CreatedAt:  "2026-01-01T00:00:00Z",
			UpdatedAt:  "2026-01-01T00:00:00Z",
		}, nil
	case "Tracker#NotFound":
		return nil, nil
	case "Tracker#Error":
		return nil, errors.New("error getting tracker")
	}
	return nil, errors.New("unsupported mock")
}

func (m *MockTrackerRepo) ListTrackers(receiverID string) ([]pkgtracker.Tracker, error) {
	switch receiverID {
	case "Receiver#123":
		return []pkgtracker.Tracker{
			{
				TrackerID:  "Tracker#123",
				ReceiverID: "Receiver#123",
				Name:       "Blood Pressure",
				Kind:       pkgtracker.KindMeasurement,
				Fields:     []pkgtracker.TrackerField{},
				IsActive:   true,
				CreatedAt:  "2026-01-01T00:00:00Z",
				UpdatedAt:  "2026-01-01T00:00:00Z",
			},
		}, nil
	case "Receiver#Empty":
		return []pkgtracker.Tracker{}, nil
	case "Receiver#Error":
		return nil, errors.New("error listing trackers")
	}
	return nil, errors.New("unsupported mock")
}

func (m *MockTrackerRepo) UpdateTracker(t *pkgtracker.Tracker) error {
	switch t.TrackerID {
	case "Tracker#123":
		return nil
	case "Tracker#Error":
		return errors.New("error updating tracker")
	}
	return errors.New("unsupported mock")
}

func (m *MockTrackerRepo) DeleteTracker(receiverID, trackerID string) error {
	switch trackerID {
	case "Tracker#123", "Tracker#456":
		return nil
	case "Tracker#Error":
		return errors.New("error deleting tracker")
	}
	return errors.New("unsupported mock")
}

func (mr *MockRelationshipRepo) GetRelationshipsByReceiver(rid string) ([]relationship.Relationship, error) {
	switch rid {
	case "Receiver#123":
		return []relationship.Relationship{
			{
				UserID:           "User#123",
				ReceiverID:       "Receiver#123",
				PrimaryCareGiver: true,
			},
			{
				UserID:           "User#456",
				ReceiverID:       "Receiver#123",
				PrimaryCareGiver: false,
			},
		}, nil
	case "Receiver#RelationshipError":
		return nil, errors.New("error retrieving relationships from db")
	case "Receiver#UserError":
		return []relationship.Relationship{
			{
				UserID:           "User#Error",
				ReceiverID:       "Receiver#UserError",
				PrimaryCareGiver: true,
			},
		}, nil
	}
	return nil, errors.New("unsupported mock")
}
