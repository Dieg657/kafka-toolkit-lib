package message

import (
	"github.com/google/uuid"
)

type IMessage[TData any] interface {
	NewForData(correlationId uuid.UUID, data TData, metadata map[string][]byte) (Message[TData], error)
	NewForDataWithKey(correlationId uuid.UUID, data TData, key string, metadata map[string][]byte) (Message[TData], error)
}
