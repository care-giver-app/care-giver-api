package events

import (
	"time"

	"github.com/google/uuid"
)

type UrinationEvent struct {
	EventID    string `json:"eventId" dynamodbav:"event_id"`
	ReceiverID string `json:"receiverId" dynamodbav:"receiver_id"`
	UserID     string `json:"userId" dynamodbav:"user_id"`
	Timestamp  string `json:"timestamp" dynamodbav:"timestamp"`
}

const eventPrefixUrination = "Urination#"

func NewUrinationEvent(receiverID, userID string) *UrinationEvent {
	return &UrinationEvent{
		EventID:    eventPrefixUrination + uuid.New().String(),
		ReceiverID: receiverID,
		UserID:     userID,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}
}

func (ue *UrinationEvent) ProcessEvent(event map[string]interface{}) error {
	err := readEvent(event, ue)
	if err != nil {
		return err
	}

	return nil
}
