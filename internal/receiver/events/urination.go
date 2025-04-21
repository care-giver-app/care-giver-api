package events

import "time"

type UrinationEvent struct {
	Timestamp string `json:"timestamp" dynamodbav:"timestamp"`
	UserID    string `json:"userId" dynamodbav:"user_id"`
}

func NewUrinationEvent() *UrinationEvent {
	return &UrinationEvent{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func (ue *UrinationEvent) ProcessEvent(event map[string]interface{}, userId string) error {
	err := readEvent(event, ue)
	if err != nil {
		return err
	}

	ue.UserID = userId

	return nil
}
