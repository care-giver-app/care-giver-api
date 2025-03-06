package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessMedicationEvent(t *testing.T) {
	tests := map[string]struct {
		event         map[string]interface{}
		expectedEvent MedicationEvent
		expectErr     bool
	}{
		"Happy Path": {
			event:         map[string]interface{}(nil),
			expectedEvent: MedicationEvent{},
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
			me := NewMedicationEvent()
			tc.expectedEvent.Timestamp = me.Timestamp

			err := me.ProcessEvent(tc.event)
			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedEvent, *me)
			}
		})
	}
}
