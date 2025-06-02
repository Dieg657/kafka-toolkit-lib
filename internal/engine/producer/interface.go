package engine

import (
	"github.com/Dieg657/kafka-toolkit-lib/pkg/common/enums"
	"github.com/Dieg657/kafka-toolkit-lib/pkg/common/message"
)

// ==========================================================================
// Interfaces
// ==========================================================================

// IKafkaProducer define a interface pública para produção de mensagens Kafka
type IKafkaProducer[TData any] interface {
	// Publish publica uma mensagem no tópico Kafka especificado
	Publish(topic string, message message.Message[TData], serialization enums.Serialization) error
}
