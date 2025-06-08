package events

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessShowerEvent(t *testing.T) {
	tests := map[string]struct {
		event         map[string]interface{}
		receiverID    string
		userID        string
		expectedEvent ShowerEvent
		expectErr     bool
	}{
		"Happy Path": {
			event:      map[string]interface{}(nil),
			receiverID: "Receiver#Test",
			userID:     "User#Test",
			expectedEvent: ShowerEvent{
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
			se := NewShowerEvent(tc.receiverID, tc.userID)
			tc.expectedEvent.Timestamp = se.Timestamp
			tc.expectedEvent.EventID = se.EventID

			err := se.ProcessEvent(tc.event)
			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedEvent, *se)
				assert.True(t, strings.HasPrefix(se.EventID, eventPrefixShower))
			}
		})
	}
}
