package consumer

import (
	"context"

	"github.com/Dieg657/kafka-toolkit-lib/internal/common/enums"
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/message"
)

type IConsumer[TData any] interface {
	ConsumeMessage(ctx context.Context, topic string, format enums.Deserialization, handler func(message message.Message[TData]) error) error
}
