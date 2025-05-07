package consumer

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/Dieg657/kafka-toolkit-lib/internal/common/enums"
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/message"
	engine "github.com/Dieg657/kafka-toolkit-lib/internal/engine/consumer"
)

// ==========================================================================
// Tipos e Propriedades
// ==========================================================================

// concreteConsumer implementa um consumidor thread-safe para mensagens Kafka
// utilizando padrão singleton para cada combinação de tipo de dados e tópico
type concreteConsumer[TData any] struct {
	ctx       context.Context // Contexto para operações assíncronas
	consumers sync.Map        // Cache de consumidores por tópico
}

// Mapa global thread-safe para armazenar instâncias singleton por tipo+tópico
var consumerRegistry sync.Map

// ==========================================================================
// Métodos Públicos
// ==========================================================================

// ConsumeMessage implementa o método da interface IConsumer
// Recebe um tópico e formato, criando ou reutilizando um consumidor existente
// para o par tipo+tópico, e inicia o consumo de mensagens.
func ConsumeMessage[TData any](ctx context.Context, topic string, format enums.Deserialization, handler func(message message.Message[TData]) error) error {
	// Gera uma chave única baseada no tipo e tópico
	key := fmt.Sprintf("%s:%s", getTypeName[TData](), topic)

	// Verifica se já existe um consumer para este par tipo+tópico
	var consumer *concreteConsumer[TData]

	if existingConsumer, exists := consumerRegistry.Load(key); exists {
		if typedConsumer, ok := existingConsumer.(*concreteConsumer[TData]); ok {
			consumer = typedConsumer
		}
	}

	// Se não existe, cria um novo
	if consumer == nil {
		consumer = &concreteConsumer[TData]{
			ctx: ctx,
		}

		// Armazena o consumidor no registro global, garantindo concorrência segura
		actualConsumer, loaded := consumerRegistry.LoadOrStore(key, consumer)
		if loaded {
			// Se outro thread criou o consumidor enquanto estávamos criando, usa o existente
			if typedConsumer, ok := actualConsumer.(*concreteConsumer[TData]); ok {
				consumer = typedConsumer
			}
		}
	}

	// Obtém ou cria um engine.Consumer específico para este tópico
	engineConsumer, err := consumer.getOrCreateConsumer(topic)
	if err != nil {
		return fmt.Errorf("falha ao preparar consumidor para tópico %s: %w", topic, err)
	}

	// Inicia o consumo usando o engine consumer
	return engineConsumer.Consume(topic, format, handler)
}

// ==========================================================================
// Métodos Privados
// ==========================================================================

// getOrCreateConsumer obtém ou cria um consumidor específico para um tópico
func (c *concreteConsumer[TData]) getOrCreateConsumer(topic string) (engine.IKafkaConsumer[TData], error) {
	// Tenta obter do cache
	if instance, exists := c.consumers.Load(topic); exists {
		if consumer, ok := instance.(engine.IKafkaConsumer[TData]); ok {
			return consumer, nil
		}
	}

	// Cria novo consumidor
	consumer, err := engine.NewKafkaConsumer[TData](c.ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar consumidor para tópico %s: %w", topic, err)
	}

	// Armazena no cache de forma thread-safe
	actual, loaded := c.consumers.LoadOrStore(topic, consumer)
	if loaded {
		if typedConsumer, ok := actual.(engine.IKafkaConsumer[TData]); ok {
			return typedConsumer, nil
		}
	}

	return consumer, nil
}

// getTypeName retorna o nome do tipo genérico
func getTypeName[TData any]() string {
	return reflect.TypeOf((*TData)(nil)).Elem().String()
}
