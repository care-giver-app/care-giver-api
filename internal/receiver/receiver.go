package receiver

import (
	"fmt"

	"github.com/care-giver-app/care-giver-api/internal/receiver/events"
	"github.com/google/uuid"
)

const (
	DBPrefix = "Receiver"
	ParamId  = "receiverId"
)

type Receiver struct {
	ReceiverID string `json:"receiverId" dynamodbav:"receiver_id"`
	FirstName  string `json:"firstName" dynamodbav:"first_name"`
	LastName   string `json:"lastName" dynamodbav:"last_name"`
}

func GenerateEvent(eventName, receiverID, userID string) (events.Event, error) {
	switch eventName {
	case "bowel_movements":
		return events.NewBowelMovementEvent(receiverID, userID), nil
	case "medications":
		return events.NewMedicationEvent(receiverID, userID), nil
	case "showers":
		return events.NewShowerEvent(receiverID, userID), nil
	case "urinations":
		return events.NewUrinationEvent(receiverID, userID), nil
	case "weights":
		return events.NewWeightEvent(receiverID, userID), nil
	default:
		return nil, fmt.Errorf("event name %s not supported", eventName)
	}
}

func NewReceiver(firstName string, lastName string) *Receiver {
	return &Receiver{
		ReceiverID: fmt.Sprintf("%s#%s", DBPrefix, uuid.New()),
		FirstName:  firstName,
		LastName:   lastName,
	}
}
