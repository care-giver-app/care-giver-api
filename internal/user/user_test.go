package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testFirstName = "Demo"
	testLastName  = "Daniel"
	testEmail     = "Demo.Daniel@email.com"
)

func TestNewUser(t *testing.T) {
	expectedUser := &User{
		FirstName:               testFirstName,
		LastName:                testLastName,
		Email:                   testEmail,
		PrimaryCareReceivers:    []string{},
		AdditionalCareReceivers: []string{},
	}

	user, err := NewUser(testEmail, testFirstName, testLastName)
	assert.Nil(t, err)

	expectedUser.UserID = user.UserID
	assert.Equal(t, expectedUser, user)
}

func TestIsACareGiver(t *testing.T) {
	user, err := NewUser(testEmail, testFirstName, testLastName)
	assert.Nil(t, err)

	tests := map[string]struct {
		primaryCareReceivers    []string
		additionalCareReceivers []string
		receiver                string
		expected                bool
	}{
		"Happy Path - Receiver in Primary List": {
			primaryCareReceivers: []string{"Receiver#1"},
			receiver:             "Receiver#1",
			expected:             true,
		},
		"Happy Path - Receiver in Additional List": {
			additionalCareReceivers: []string{"Receiver#1"},
			receiver:                "Receiver#1",
			expected:                true,
		},
		"Happy Path - Receiver in Primary List With Additional Entries": {
			primaryCareReceivers:    []string{"Receiver#1", "Receiver#2"},
			additionalCareReceivers: []string{"Receiver#3", "Receiver#4"},
			receiver:                "Receiver#1",
			expected:                true,
		},
		"Happy Path - Receiver in Additional List With Additional Entries": {
			primaryCareReceivers:    []string{"Receiver#1", "Receiver#2"},
			additionalCareReceivers: []string{"Receiver#3", "Receiver#4"},
			receiver:                "Receiver#3",
			expected:                true,
		},
		"Sad Path - Receiver Not in Lists": {
			primaryCareReceivers:    []string{"Receiver#1", "Receiver#2"},
			additionalCareReceivers: []string{"Receiver#3", "Receiver#4"},
			receiver:                "Receiver#5",
			expected:                false,
		},
		"Sad Path - Lists Are Empty": {
			receiver: "Receiver#5",
			expected: false,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			user.PrimaryCareReceivers = tc.primaryCareReceivers
			user.AdditionalCareReceivers = tc.additionalCareReceivers

			isCareGiver := user.IsACareGiver(tc.receiver)
			assert.Equal(t, tc.expected, isCareGiver)
		})
	}
}
