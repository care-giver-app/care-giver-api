package user

import (
	"crypto/sha256"
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
	UserID                  string   `json:"userId" dynamodbav:"user_id"`
	Email                   string   `json:"email" dynamodbav:"email"`
	Password                []byte   `json:"password" dynamodbav:"password"`
	FirstName               string   `json:"firstName" dynamodbav:"first_name"`
	LastName                string   `json:"lastName" dynamodbav:"last_name"`
	PrimaryCareReceivers    []string `json:"primaryCareReceivers" dynamodbav:"primary_care_receivers"`
	AdditionalCareReceivers []string `json:"additionalCareReceivers" dynamodbav:"additional_care_receivers"`
}

func NewUser(email string, password string, firstName string, lastName string) (*User, error) {
	hashFunc := sha256.New()
	_, err := hashFunc.Write([]byte(password))
	if err != nil {
		return nil, err
	}
	encryptedPassword := hashFunc.Sum(nil)

	return &User{
		UserID:                  fmt.Sprintf("%s#%s", DBPrefix, uuid.New()),
		Email:                   email,
		Password:                encryptedPassword,
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
