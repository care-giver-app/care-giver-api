package models

type UserID string
type ReceiverID string

type TimestampEvent struct {
	timestamp string
}

type WeightEvent struct {
	timestamp string
	weight    string
}

type User struct {
	userID                  UserID
	email                   string
	password                string
	primaryCareReceivers    []ReceiverID
	additionalCareReceivers []ReceiverID
}

type Receiver struct {
	receiverID     ReceiverID
	firstName      string
	lastName       string
	medications    []TimestampEvent
	showers        []TimestampEvent
	urinations     []TimestampEvent
	bowelMovements []TimestampEvent
	Weight         []WeightEvent
}
