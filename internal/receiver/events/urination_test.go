package events

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessUrinationEvent(t *testing.T) {
	tests := map[string]struct {
		event         map[string]interface{}
		receiverID    string
		userID        string
		expectedEvent UrinationEvent
		expectErr     bool
	}{
		"Happy Path": {
			event:      map[string]interface{}(nil),
			receiverID: "Receiver#Test",
			userID:     "User#Test",
			expectedEvent: UrinationEvent{
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
			ue := NewUrinationEvent(tc.receiverID, tc.userID)
			tc.expectedEvent.Timestamp = ue.Timestamp
			tc.expectedEvent.EventID = ue.EventID

			err := ue.ProcessEvent(tc.event)
			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedEvent, *ue)
				assert.True(t, strings.HasPrefix(ue.EventID, eventPrefixUrination))
			}
		})
	}
}
