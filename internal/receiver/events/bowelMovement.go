package events

import (
	"time"

	"github.com/google/uuid"
)

type BowelMovementEvent struct {
	EventID    string `json:"eventId" dynamodbav:"event_id"`
	ReceiverID string `json:"receiverId" dynamodbav:"receiver_id"`
	UserID     string `json:"userId" dynamodbav:"user_id"`
	Timestamp  string `json:"timestamp" dynamodbav:"timestamp"`
}

const eventPrefixBowelMovement = "BowelMovement#"

func NewBowelMovementEvent(receiverID, userID string) *BowelMovementEvent {
	return &BowelMovementEvent{
		EventID:    eventPrefixBowelMovement + uuid.New().String(),
		ReceiverID: receiverID,
		UserID:     userID,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}
}

func (bme *BowelMovementEvent) ProcessEvent(event map[string]interface{}) error {
	err := readEvent(event, bme)
	if err != nil {
		return err
	}

	return nil
}
