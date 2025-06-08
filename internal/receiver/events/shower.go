package events

import (
	"time"

	"github.com/google/uuid"
)

type ShowerEvent struct {
	EventID    string `json:"eventId" dynamodbav:"event_id"`
	ReceiverID string `json:"receiverId" dynamodbav:"receiver_id"`
	UserID     string `json:"userId" dynamodbav:"user_id"`
	Timestamp  string `json:"timestamp" dynamodbav:"timestamp"`
}

const eventPrefixShower = "Shower#"

func NewShowerEvent(receiverID, userID string) *ShowerEvent {
	return &ShowerEvent{
		EventID:    eventPrefixShower + uuid.New().String(),
		ReceiverID: receiverID,
		UserID:     userID,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}
}

func (se *ShowerEvent) ProcessEvent(event map[string]interface{}) error {
	err := readEvent(event, se)
	if err != nil {
		return err
	}

	return nil
}
