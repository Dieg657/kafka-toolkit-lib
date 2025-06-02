package message

import (
	"errors"
	"reflect"

	"github.com/google/uuid"
)

type Message[TData any] struct {
	CorrelationId uuid.UUID
	Data          TData
	Metadata      map[string][]byte
}

func NewForData[TData any](correlationId uuid.UUID, data TData, metadata map[string][]byte) (Message[TData], error) {
	if reflect.ValueOf(data).IsZero() {
		return Message[TData]{}, errors.New("data cannot be empty")
	}

	message := Message[TData]{}
	message.Data = data
	message.CorrelationId = correlationId
	message.Metadata = make(map[string][]byte, len(metadata)+1)
	for key, value := range metadata {
		message.Metadata[key] = value
	}
	message.Metadata["correlationId"] = []byte(correlationId.String())
	return message, nil
}

func NewForDataWithKey[TData any](correlationId uuid.UUID, data TData, key string, metadata map[string][]byte) (Message[TData], error) {
	if reflect.ValueOf(data).IsZero() {
		return Message[TData]{}, errors.New("data cannot be empty")
	}

	message := Message[TData]{}
	message.Data = data
	message.CorrelationId = correlationId
	message.Metadata = make(map[string][]byte, len(metadata)+2)
	message.Metadata[key] = []byte(key)
	message.Metadata["correlationId"] = []byte(correlationId.String())
	for key, value := range metadata {
		message.Metadata[key] = value
	}
	return message, nil
}
