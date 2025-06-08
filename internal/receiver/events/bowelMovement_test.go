package events

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessBowelMovementEvent(t *testing.T) {
	tests := map[string]struct {
		event         map[string]interface{}
		receiverID    string
		userID        string
		expectedEvent BowelMovementEvent
		expectErr     bool
	}{
		"Happy Path": {
			event:      map[string]interface{}(nil),
			receiverID: "Receiver#Test",
			userID:     "User#Test",
			expectedEvent: BowelMovementEvent{
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
			bme := NewBowelMovementEvent(tc.receiverID, tc.userID)
			tc.expectedEvent.Timestamp = bme.Timestamp
			tc.expectedEvent.EventID = bme.EventID

			err := bme.ProcessEvent(tc.event)
			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedEvent, *bme)
				assert.True(t, strings.HasPrefix(bme.EventID, eventPrefixBowelMovement))
			}
		})
	}
}
