package receiver

import (
	"testing"

	"github.com/care-giver-app/care-giver-api/internal/receiver/events"
	"github.com/stretchr/testify/assert"
)

func TestNewReceiver(t *testing.T) {
	testFirstName := "Demo"
	testLastName := "Dan"
	expectedReceiver := &Receiver{
		FirstName:      testFirstName,
		LastName:       testLastName,
		Medications:    []events.MedicationEvent{},
		Showers:        []events.ShowerEvent{},
		Urinations:     []events.UrinationEvent{},
		BowelMovements: []events.BowelMovementEvent{},
		Weights:        []events.WeightEvent{},
	}

	receiver := NewReceiver(testFirstName, testLastName)
	expectedReceiver.ReceiverID = receiver.ReceiverID

	assert.Equal(t, expectedReceiver, receiver)
}
