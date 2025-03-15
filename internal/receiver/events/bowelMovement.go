package events

import "time"

type BowelMovementEvent struct {
	Timestamp string `json:"timestamp" dynamodbav:"timestamp"`
	UserID    string `json:"userId" dynamodbav:"user_id"`
}

func NewBowelMovementEvent() *BowelMovementEvent {
	return &BowelMovementEvent{
		Timestamp: time.Now().UTC().String(),
	}
}

func (bme *BowelMovementEvent) ProcessEvent(event map[string]interface{}, userId string) error {
	err := readEvent(event, bme)
	if err != nil {
		return err
	}

	bme.UserID = userId

	return nil
}
