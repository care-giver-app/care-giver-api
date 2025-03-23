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

type ReceiverID string

type Receiver struct {
	ReceiverID     ReceiverID                  `json:"receiverId" dynamodbav:"receiver_id"`
	FirstName      string                      `json:"firstName" dynamodbav:"first_name"`
	LastName       string                      `json:"lastName" dynamodbav:"last_name"`
	Medications    []events.MedicationEvent    `json:"medications" dynamodbav:"medications"`
	Showers        []events.ShowerEvent        `json:"showers" dynamodbav:"showers"`
	Urinations     []events.UrinationEvent     `json:"urinations" dynamodbav:"urinations"`
	BowelMovements []events.BowelMovementEvent `json:"bowelMovements" dynamodbav:"bowel_movements"`
	Weights        []events.WeightEvent        `json:"weights" dynamodbav:"weights"`
}

func EventFromName(eventName string) (events.Event, error) {
	switch eventName {
	case "bowel_movements":
		return events.NewBowelMovementEvent(), nil
	case "medications":
		return events.NewMedicationEvent(), nil
	case "showers":
		return events.NewShowerEvent(), nil
	case "urinations":
		return events.NewUrinationEvent(), nil
	case "weights":
		return events.NewWeightEvent(), nil
	default:
		return nil, fmt.Errorf("event name %s not supported", eventName)
	}
}

func NewReceiver(firstName string, lastName string) *Receiver {
	return &Receiver{
		ReceiverID:     ReceiverID(fmt.Sprintf("%s#%s", DBPrefix, uuid.New())),
		FirstName:      firstName,
		LastName:       lastName,
		Medications:    []events.MedicationEvent{},
		Showers:        []events.ShowerEvent{},
		Urinations:     []events.UrinationEvent{},
		BowelMovements: []events.BowelMovementEvent{},
		Weights:        []events.WeightEvent{},
	}
}
