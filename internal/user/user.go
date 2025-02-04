package user

import "github.com/care-giver-app/care-giver-api/internal/receiver"

type UserID string

type User struct {
	userID                  UserID
	email                   string
	password                string
	primaryCareReceivers    []receiver.ReceiverID
	additionalCareReceivers []receiver.ReceiverID
}
