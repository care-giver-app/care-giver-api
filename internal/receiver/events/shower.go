package events

import "time"

type ShowerEvent struct {
	Timestamp string `json:"timestamp" dynamodbav:"timestamp"`
	UserID    string `json:"userId" dynamodbav:"user_id"`
}

func NewShowerEvent() *ShowerEvent {
	return &ShowerEvent{
		Timestamp: time.Now().UTC().String(),
	}
}

func (se *ShowerEvent) ProcessEvent(event map[string]interface{}, userId string) error {
	err := readEvent(event, se)
	if err != nil {
		return err
	}

	se.UserID = userId

	return nil
}
