package events

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessMedicationEvent(t *testing.T) {
	tests := map[string]struct {
		event         map[string]interface{}
		receiverID    string
		userID        string
		expectedEvent MedicationEvent
		expectErr     bool
	}{
		"Happy Path": {
			event:      map[string]interface{}(nil),
			receiverID: "Receiver#Test",
			userID:     "User#Test",
			expectedEvent: MedicationEvent{
				ReceiverID: "Receiver#Test",
				UserID:     "User#Test",
			},
		},
		"Sad Path - Unknown Fields": {
			event: map[string]interface{}{
				"field1": "test",
			},
			expectErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			me := NewMedicationEvent(tc.receiverID, tc.userID)
			tc.expectedEvent.Timestamp = me.Timestamp
			tc.expectedEvent.EventID = me.EventID

			err := me.ProcessEvent(tc.event)
			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedEvent, *me)
				assert.True(t, strings.HasPrefix(me.EventID, eventPrefixMedication))
			}
		})
	}
}
