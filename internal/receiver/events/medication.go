package events

import (
	"time"

	"github.com/google/uuid"
)

type MedicationEvent struct {
	EventID    string `json:"eventId" dynamodbav:"event_id"`
	ReceiverID string `json:"receiverId" dynamodbav:"receiver_id"`
	UserID     string `json:"userId" dynamodbav:"user_id"`
	Timestamp  string `json:"timestamp" dynamodbav:"timestamp"`
}

const eventPrefixMedication = "Medication#"

func NewMedicationEvent(receiverID, userID string) *MedicationEvent {
	return &MedicationEvent{
		EventID:    eventPrefixMedication + uuid.New().String(),
		ReceiverID: receiverID,
		UserID:     userID,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}
}

func (me *MedicationEvent) ProcessEvent(event map[string]interface{}) error {
	err := readEvent(event, me)
	if err != nil {
		return err
	}

	return nil
}
