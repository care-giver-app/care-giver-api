package events

import "time"

type BowelMovementEvent struct {
	Timestamp string `json:"timestamp" dynamodbav:"timestamp"`
}

func NewBowelMovementEvent() *BowelMovementEvent {
	return &BowelMovementEvent{
		Timestamp: time.Now().UTC().String(),
	}
}

func (bme *BowelMovementEvent) ProcessEvent(event map[string]interface{}) error {
	err := readEvent(event, bme)
	if err != nil {
		return err
	}
	return nil
}
