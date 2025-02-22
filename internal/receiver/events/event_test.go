package events

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestEventStruct struct {
	FieldOne   string  `json:"fieldOne" validate:"required"`
	FieldTwo   float32 `json:"fieldTwo" validate:"required"`
	FieldThree string  `json:"fieldThree"`
}

func TestReadEvent(t *testing.T) {
	tests := map[string]struct {
		event         map[string]interface{}
		expectedEvent TestEventStruct
		expectErr     bool
	}{
		"Happy Path - Only Required Fields": {
			event: map[string]interface{}{
				"fieldOne": "test val",
				"fieldTwo": 120.32,
			},
			expectedEvent: TestEventStruct{
				FieldOne: "test val",
				FieldTwo: 120.32,
			},
		},
		"Happy Path - All Fields": {
			event: map[string]interface{}{
				"fieldOne":   "test val",
				"fieldTwo":   120.32,
				"fieldThree": "optional field",
			},
			expectedEvent: TestEventStruct{
				FieldOne:   "test val",
				FieldTwo:   120.32,
				FieldThree: "optional field",
			},
		},
		"Sad Path - Unknown Fields": {
			event: map[string]interface{}{
				"fieldOne":  "test val",
				"fieldTwo":  120.32,
				"fieldFour": "test",
			},
			expectErr: true,
		},
		"Sad Path - Fields Are Wrong Type": {
			event: map[string]interface{}{
				"fieldOne": "test val",
				"fieldTwo": "NAN",
			},
			expectErr: true,
		},
		"Sad Path - Missing Fields": {
			event: map[string]interface{}{
				"fieldOne": "test val",
			},
			expectErr: true,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var testEventStruct TestEventStruct
			err := readEvent(tc.event, &testEventStruct)

			if tc.expectErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedEvent, testEventStruct)
			}
		})
	}
}
