package events

import (
	"bytes"
	"encoding/json"

	"github.com/go-playground/validator/v10"
)

type Event interface {
	ProcessEvent(event map[string]interface{}) error
}

func readEvent(event map[string]interface{}, eventStruct interface{}) error {
	jsonString, err := json.Marshal(event)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(bytes.NewReader(jsonString))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(eventStruct)
	if err != nil {
		return err
	}

	validate := validator.New()
	err = validate.Struct(eventStruct)
	if err != nil {
		return err
	}

	return nil
}
