package consumer

import (
	"context"

	"github.com/Dieg657/kafka-toolkit-lib/internal/common/enums"
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/message"
)

// ==========================================================================
// Interfaces
// ==========================================================================

// IConsumer define a interface para consumo de mensagens
// O parâmetro genérico TData define o tipo de dado que será consumido
type IConsumer[TData any] interface {
	// ConsumeMessage inicia o consumo de mensagens de um tópico, com o formato especificado pelo handler
	ConsumeMessage(ctx context.Context, topic string, format enums.Deserialization, strategy enums.DeserializationStrategy, handler func(message message.Message[TData]) error) error
	// ConsumeMessageInParallel inicia o consumo de mensagens de um tópico, utilizando múltiplos workers paralelos
	// e um tamanho de lote para otimizar o consumo de mensagens
	ConsumeMessageInParallel(ctx context.Context, parallelWorkers int, batchSize int, topic string, format enums.Deserialization, strategy enums.DeserializationStrategy, handler func(message message.Message[TData]) error) error
}
