package events

import "time"

type MedicationEvent struct {
	Timestamp string `json:"timestamp" dynamodbav:"timestamp"`
	UserID    string `json:"userId" dynamodbav:"user_id"`
}

func NewMedicationEvent() *MedicationEvent {
	return &MedicationEvent{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

func (me *MedicationEvent) ProcessEvent(event map[string]interface{}, userId string) error {
	err := readEvent(event, me)
	if err != nil {
		return err
	}

	me.UserID = userId

	return nil
}
