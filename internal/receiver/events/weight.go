package events

import (
	"errors"
	"time"
)

type WeightEvent struct {
	Timestamp string  `json:"timestamp" dynamodbav:"timestamp"`
	Weight    float32 `json:"weight" dynamodbav:"weight" validate:"required"`
}

func NewWeightEvent() *WeightEvent {
	return &WeightEvent{
		Timestamp: time.Now().UTC().String(),
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
