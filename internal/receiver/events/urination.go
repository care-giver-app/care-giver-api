package events

import "time"

type UrinationEvent struct {
	Timestamp string `json:"timestamp" dynamodbav:"timestamp"`
}

func NewUrinationEvent() *UrinationEvent {
	return &UrinationEvent{
		Timestamp: time.Now().UTC().String(),
	}
}

func (ue *UrinationEvent) ProcessEvent(event map[string]interface{}) error {
	err := readEvent(event, ue)
	if err != nil {
		return err
	}
	return nil
}
