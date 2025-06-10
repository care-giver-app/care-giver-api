package user

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
)

const (
	ParamId  = "userId"
	DBPrefix = "User"
)

type UserID string

type User struct {
	UserID                  UserID                `json:"userId" dynamodbav:"user_id"`
	Email                   string                `json:"email" dynamodbav:"email"`
	FirstName               string                `json:"firstName" dynamodbav:"first_name"`
	LastName                string                `json:"lastName" dynamodbav:"last_name"`
	PrimaryCareReceivers    []receiver.ReceiverID `json:"primaryCareReceivers" dynamodbav:"primary_care_receivers"`
	AdditionalCareReceivers []receiver.ReceiverID `json:"additionalCareReceivers" dynamodbav:"additional_care_receivers"`
}

func NewUser(email string, firstName string, lastName string) (*User, error) {
	return &User{
		UserID:                  fmt.Sprintf("%s#%s", DBPrefix, uuid.New()),
		Email:                   email,
		FirstName:               firstName,
		LastName:                lastName,
		PrimaryCareReceivers:    []string{},
		AdditionalCareReceivers: []string{},
	}, nil
}

func (u *User) IsACareGiver(rid string) bool {
	receiverList := append(u.PrimaryCareReceivers, u.AdditionalCareReceivers...)
	return slices.Contains(receiverList, rid)
}
