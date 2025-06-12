package handlers

import (
	"errors"

	"github.com/care-giver-app/care-giver-api/internal/event"
	"github.com/care-giver-app/care-giver-api/internal/receiver"
	"github.com/care-giver-app/care-giver-api/internal/user"
)

var (
	testUserRepo     = &MockUserRepo{}
	testReceiverRepo = &MockReceiverRepo{}
	testEventRepo    = &MockEventRepo{}
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
			PrimaryCareReceivers: []string{"Receiver#123", "Receiver#Error"},
		}, nil
	case "User#NotACareGiver":
		return user.User{}, nil
	case "User#Error":
		return user.User{}, errors.New("error getting user from db")
	}
	return user.User{}, errors.New("unsupported mock")
}

func (mu *MockUserRepo) UpdateReceiverList(uid string, rid string, listName string) error {
	switch uid {
	case "User#123":
		return nil
	case "User#ListError":
		return errors.New("error updating receiver list")
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
