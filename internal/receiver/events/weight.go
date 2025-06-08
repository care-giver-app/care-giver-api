package events

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type WeightEvent struct {
	EventID    string  `json:"eventId" dynamodbav:"event_id"`
	ReceiverID string  `json:"receiverId" dynamodbav:"receiver_id"`
	UserID     string  `json:"userId" dynamodbav:"user_id"`
	Timestamp  string  `json:"timestamp" dynamodbav:"timestamp"`
	Weight     float32 `json:"weight" dynamodbav:"weight" validate:"required"`
}

const eventPrefixWeight = "Weight#"

func NewWeightEvent(receiverID, userID string) *WeightEvent {
	return &WeightEvent{
		EventID:    eventPrefixWeight + uuid.New().String(),
		ReceiverID: receiverID,
		UserID:     userID,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}
}

func (we *WeightEvent) ProcessEvent(event map[string]interface{}) error {
	if event == nil {
		return errors.New("no weight event provided")
	}

	err := readEvent(event, we)
	if err != nil {
		return err
	}

	return nil
}
