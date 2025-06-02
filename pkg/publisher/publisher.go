package publisher

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	engine "github.com/Dieg657/kafka-toolkit-lib/internal/engine/producer"
	"github.com/Dieg657/kafka-toolkit-lib/pkg/common/enums"
	"github.com/Dieg657/kafka-toolkit-lib/pkg/common/message"
)

// ==========================================================================
// Tipos e Propriedades Estáticas
// ==========================================================================

// Publisher implementa a interface IPublisher e fornece uma implementação
// thread-safe para publicação de mensagens em tópicos Kafka.
// Utiliza um padrão singleton e mantém um cache de produtores por tópico.
type concretePublisher[TData any] struct {
	ctx       context.Context // Contexto usado para criar produtores
	producers sync.Map        // Mapa thread-safe de produtores indexados por tópico
}

// Mapa global thread-safe para armazenar instâncias singleton por tipo
var publisherRegistry sync.Map

// ==========================================================================
// Construtores
// ==========================================================================

// New cria ou retorna a instância singleton do Publisher.
// Implementa o padrão singleton garantindo que apenas uma instância
// do Publisher exista em toda a aplicação para cada tipo TData.
//
// Parâmetros:
//   - ctx: Contexto que será usado pelos produtores para operações assíncronas
//
// Retorno:
//   - *Publisher[TData]: Instância do Publisher para o tipo TData
func New[TData any](ctx context.Context) *concretePublisher[TData] {
	// Usar o nome do tipo como chave para o mapa de instâncias
	typeName := getTypeName[TData]()

	// Verifica se já existe uma instância para este tipo
	if existingInstance, exists := publisherRegistry.Load(typeName); exists {
		if instance, ok := existingInstance.(*concretePublisher[TData]); ok {
			return instance
		}
	}

	// Criar uma nova instância para este tipo
	newInstance := &concretePublisher[TData]{
		ctx: ctx,
	}

	// Armazenar a instância no mapa (usa LoadOrStore para evitar condição de corrida)
	actualInstance, _ := publisherRegistry.LoadOrStore(typeName, newInstance)

	// Retorna a instância que foi realmente armazenada (pode ser a nossa ou outra criada concorrentemente)
	if instance, ok := actualInstance.(*concretePublisher[TData]); ok {
		return instance
	}

	// Este caso só aconteceria se houvesse um erro grave de tipo
	return newInstance
}

// ==========================================================================
// Métodos Públicos
// ==========================================================================

// PublishMessage publica uma mensagem em um tópico Kafka usando o formato de serialização especificado.
// Preserva o tipo específico TData durante todo o processo para garantir a serialização correta.
//
// Parâmetros:
//   - topic: Nome do tópico Kafka onde a mensagem será publicada
//   - message: Mensagem tipada a ser publicada
//   - serialization: Formato de serialização a ser usado (Avro, JSON, etc.)
//
// Retorno:
//   - error: Erro caso ocorra falha na publicação
func (p *concretePublisher[TData]) PublishMessage(topic string, message message.Message[TData], serialization enums.Serialization) error {
	// Obtém ou cria um produtor fortemente tipado para o tópico
	producer, err := p.getOrCreateProducer(topic)
	if err != nil {
		return err
	}

	// Usa diretamente o produtor com o tipo TData, sem conversões intermediárias
	// que poderiam fazer perder informações da interface Avro
	return producer.Publish(topic, message, serialization)
}

// PublishMessage é uma função estática que centraliza a instanciação e publicação em uma única chamada.
// A função cria ou reutiliza um Publisher do tipo apropriado e publica a mensagem.
//
// Parâmetros:
//   - ctx: Contexto usado para operações assíncronas
//   - topic: Nome do tópico Kafka onde a mensagem será publicada
//   - message: Mensagem tipada a ser publicada
//   - serialization: Formato de serialização a ser usado
//
// Retorno:
//   - error: Erro caso ocorra falha na publicação
func PublishMessage[TData any](ctx context.Context, topic string, message message.Message[TData], serialization enums.Serialization) error {
	// Usa a função New para obter ou criar uma instância do publisher
	publisher := New[TData](ctx)

	// Publica a mensagem usando o publisher obtido
	return publisher.PublishMessage(topic, message, serialization)
}

// ==========================================================================
// Métodos Privados
// ==========================================================================

// getTypeName retorna o nome do tipo TData como string sem alocar memória desnecessariamente
func getTypeName[TData any]() string {
	return reflect.TypeOf((*TData)(nil)).Elem().String()
}

// getOrCreateProducer obtém um produtor existente do cache ou cria um novo se necessário.
// Retorna um produtor fortemente tipado como *engine.Producer[TData] para preservar informações de tipo.
// Usa sync.Map para eliminar necessidade de locks manuais e melhorar concorrência.
func (p *concretePublisher[TData]) getOrCreateProducer(topic string) (engine.IKafkaProducer[TData], error) {
	// Tenta obter um produtor existente do mapa thread-safe
	if existingProducer, exists := p.producers.Load(topic); exists {
		// Tenta converter para o tipo Producer[TData]
		if typedProducer, ok := existingProducer.(engine.IKafkaProducer[TData]); ok {
			return typedProducer, nil
		}
	}

	// Se não existir ou não for do tipo correto, cria um novo
	newProducer, err := engine.NewKafkaProducer[TData](p.ctx)
	if err != nil {
		return nil, fmt.Errorf("falha ao criar novo produtor para tópico '%s': %w", topic, err)
	}

	// Armazena o novo produtor no cache de forma thread-safe
	// LoadOrStore evita que múltiplas goroutines criem produtores duplicados
	actualProducer, loaded := p.producers.LoadOrStore(topic, newProducer)

	// Se já existe um produtor, use o existente
	if loaded {
		if typedProducer, ok := actualProducer.(engine.IKafkaProducer[TData]); ok {
			return typedProducer, nil
		}
	}

	// Caso contrário, use o novo produtor criado
	return newProducer, nil
}
