package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessShowerEvent(t *testing.T) {
	tests := map[string]struct {
		event         map[string]interface{}
		expectedEvent ShowerEvent
		expectErr     bool
	}{
		"Happy Path": {
			event:         map[string]interface{}(nil),
			expectedEvent: ShowerEvent{},
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
			se := NewShowerEvent()
			tc.expectedEvent.Timestamp = se.Timestamp

			err := se.ProcessEvent(tc.event)
			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedEvent, *se)
			}
		})
	}
}
