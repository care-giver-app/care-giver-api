package events

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessWeightEvent(t *testing.T) {
	tests := map[string]struct {
		event         map[string]interface{}
		receiverID    string
		userID        string
		expectedEvent WeightEvent
		expectErr     bool
	}{
		"Happy Path": {
			event: map[string]interface{}{
				"weight": 130.3,
			},
			receiverID: "Receiver#Test",
			userID:     "User#Test",
			expectedEvent: WeightEvent{
				Weight:     130.3,
				ReceiverID: "Receiver#Test",
				UserID:     "User#Test",
			},
		},
		"Sad Path - No Event Provided": {
			event:     map[string]interface{}(nil),
			expectErr: true,
		},
		"Sad Path - Weight Not Provided": {
			event:     map[string]interface{}{},
			expectErr: true,
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
			we := NewWeightEvent(tc.receiverID, tc.userID)
			tc.expectedEvent.Timestamp = we.Timestamp
			tc.expectedEvent.EventID = we.EventID

			err := we.ProcessEvent(tc.event)
			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedEvent, *we)
				assert.True(t, strings.HasPrefix(we.EventID, eventPrefixWeight))
			}
		})
	}
}
