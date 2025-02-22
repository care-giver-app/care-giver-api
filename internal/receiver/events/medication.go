package events

import "time"

type MedicationEvent struct {
	Timestamp string `json:"timestamp" dynamodbav:"timestamp"`
}

func NewMedicationEvent() *MedicationEvent {
	return &MedicationEvent{
		Timestamp: time.Now().UTC().String(),
	}
}

func (me *MedicationEvent) ProcessEvent(event map[string]interface{}) error {
	err := readEvent(event, me)
	if err != nil {
		return err
	}
	return nil
}
