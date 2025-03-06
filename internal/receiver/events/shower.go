package events

import "time"

type ShowerEvent struct {
	Timestamp string `json:"timestamp" dynamodbav:"timestamp"`
}

func NewShowerEvent() *ShowerEvent {
	return &ShowerEvent{
		Timestamp: time.Now().UTC().String(),
	}
}

func (se *ShowerEvent) ProcessEvent(event map[string]interface{}) error {
	err := readEvent(event, se)
	if err != nil {
		return err
	}
	return nil
}
