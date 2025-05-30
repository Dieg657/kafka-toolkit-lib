package engine

import (
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/enums"
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/message"
)

// ==========================================================================
// Interfaces
// ==========================================================================

// IKafkaConsumer define a interface pública para consumo de mensagens Kafka
type IKafkaConsumer[TData any] interface {
	// Consume inicia o consumo de mensagens de um tópico Kafka
	Consume(topic string, deserialization enums.Deserialization, strategy enums.DeserializationStrategy, handler func(message message.Message[TData]) error) error
}
