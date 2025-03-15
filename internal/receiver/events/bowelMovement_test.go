package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessBowelMovementEvent(t *testing.T) {
	tests := map[string]struct {
		event         map[string]interface{}
		userId        string
		expectedEvent BowelMovementEvent
		expectErr     bool
	}{
		"Happy Path": {
			event:  map[string]interface{}(nil),
			userId: "User#Test",
			expectedEvent: BowelMovementEvent{
				UserID: "User#Test",
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
			bme := NewBowelMovementEvent()
			tc.expectedEvent.Timestamp = bme.Timestamp

			err := bme.ProcessEvent(tc.event, tc.userId)
			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedEvent, *bme)
			}
		})
	}
}
