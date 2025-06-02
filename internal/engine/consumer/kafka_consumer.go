package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"

	internalEnums "github.com/Dieg657/kafka-toolkit-lib/internal/common/enums"
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/setup"
	"github.com/Dieg657/kafka-toolkit-lib/internal/engine/adapter"
	"github.com/Dieg657/kafka-toolkit-lib/pkg/common/constants"
	"github.com/Dieg657/kafka-toolkit-lib/pkg/common/enums"
	"github.com/Dieg657/kafka-toolkit-lib/pkg/common/ioc"
	"github.com/Dieg657/kafka-toolkit-lib/pkg/common/message"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
)

// ==========================================================================
// Tipos e Propriedades Estáticas
// ==========================================================================

// KafkaConsumer encapsula um consumidor Kafka fortemente tipado.
// Gerencia a conexão com o Kafka e deserialização de mensagens.
type kafkaConsumer[TData any] struct {
	ctx              context.Context
	client           *kafka.Consumer
	priority         internalEnums.ConsumerOrderPriority
	registry         setup.ISchemaRegistrySetup
	deserializerMap  map[enums.Deserialization]func(topic string, payload []byte, data *TData) error
	protobufAdapter  *adapter.ProtobufAdapter
	consumerPriority internalEnums.ConsumerOrderPriority
}

var (
	deserializationStrategyMap = map[enums.DeserializationStrategy]func(err error){
		enums.OnDeserializationFailedStopHost: func(err error) {
			if err != nil {
				fmt.Printf("Error on deserialize message: %s\n", err)
				panic(err)
			}
		},
		enums.OnDeserializationIgnoreMessage: func(err error) {
			if err != nil {
				fmt.Printf("Error on deserialize message: %s\n", err)
			}
		},
	}
)

// ==========================================================================
// Construtores
// ==========================================================================

// NewKafkaConsumer cria uma nova instância de IKafkaConsumer
//
// Parâmetros:
//   - ctx: Contexto contendo as dependências e configurações
//
// Retorno:
//   - IKafkaConsumer: Interface do consumidor
//   - error: Erro caso a inicialização falhe
func NewKafkaConsumer[TData any](ctx context.Context) (IKafkaConsumer[TData], error) {
	container := ctx.Value(constants.IocKey).(ioc.IContainer)
	if container == nil {
		return nil, errors.New("IoC do Kafka não encontrado no contexto")
	}

	consumerSetup, schemaRegistry := container.GetConsumer()
	consumer := &kafkaConsumer[TData]{}
	consumer.ctx = ctx
	consumer.client = consumerSetup.GetKafkaConsumer()
	consumer.priority = container.GetConsumerPriority()
	consumer.registry = schemaRegistry
	consumer.protobufAdapter = adapter.NewProtobufAdapter() // Inicializa o adaptador protobuf

	// Obtém a prioridade do consumidor
	consumer.consumerPriority = container.GetConsumerPriority()

	// Inicializa os deserializadores suportados
	consumer.initializeDeserializers()

	return consumer, nil
}

// ==========================================================================
// Métodos Públicos
// ==========================================================================

// Consume inicia o consumo de mensagens de um tópico Kafka.
// Deserializa as mensagens com o formato especificado e chama o handler para cada uma delas.
//
// Parâmetros:
//   - parallelWorkers: Número de workers paralelos para consumir mensagens
//   - batchSize: Tamanho do lote de mensagens a ser consumido
//   - topic: Nome do tópico Kafka para consumo
//   - deserialization: Formato de deserialização a ser usado
//   - strategy: Estratégia de tratamento de erro de deserialização
//   - handler: Função a ser chamada para cada mensagem consumida
//
// Retorno:
//   - error: Erro caso ocorra falha no consumo
func (c *kafkaConsumer[TData]) Consume(topic string, deserialization enums.Deserialization, strategy enums.DeserializationStrategy, handler func(message message.Message[TData]) error) error {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	err := c.client.Subscribe(topic, rebalanceCallback)

	if err != nil {
		return err
	}

	run := true

	for run {
		select {
		case sig := <-sigchan:
			fmt.Printf("Terminating consumer signal %v cought!\n", sig)
			run = false
		case <-c.ctx.Done():
			fmt.Printf("Terminating consumer context done!\n")
			run = false
		default:
			ev := c.client.Poll(100)
			if ev == nil {
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				baseMessage := message.Message[TData]{
					Metadata: make(map[string][]byte, len(e.Headers)),
				}
				c.fillHeader(&baseMessage, e.Headers)

				err = c.deserializeValue(e, deserialization, &baseMessage.Data)
				deserializationStrategyMap[strategy](err)

				err = handler(baseMessage)
				if err != nil {
					fmt.Println("Error on handle message")
				}

				if c.priority == internalEnums.CONSUMER_ORDER_PRIORITY_HIGH_PERFORMANCE || c.priority == internalEnums.CONSUMER_ORDER_PRIORITY_RISKY {
					continue
				}

				err = commitMessage(c.client, e.TopicPartition)
				if err != nil {
					fmt.Println("Error on commit message")
				}
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
			default:
				fmt.Printf("Ignored %v\n", e)
			}
		}
	}

	return nil
}

// ==========================================================================
// Métodos Privados
// ==========================================================================

// initializeDeserializers configura todos os deserializadores suportados pela biblioteca.
// Centraliza a criação e registro dos deserializadores padrão.
func (c *kafkaConsumer[TData]) initializeDeserializers() {
	c.deserializerMap = map[enums.Deserialization]func(topic string, payload []byte, data *TData) error{
		enums.JsonDeserialization: func(topic string, payload []byte, data *TData) error {
			err := json.Unmarshal(payload, data)
			if err != nil {
				return err
			}
			return nil
		},
		enums.JsonSchemaDeserialization: func(topic string, payload []byte, data *TData) error {
			return c.registry.GetJsonDeserializer().DeserializeInto(topic, payload, data)
		},
		enums.AvroDeserialization: func(topic string, payload []byte, data *TData) error {
			return c.registry.GetAvroDeserializer().DeserializeInto(topic, payload, data)
		},
		enums.ProtobufDeserialization: func(topic string, payload []byte, data *TData) error {
			// Tenta deserializar diretamente para o tipo alvo
			err := c.registry.GetProtobufDeserializer().DeserializeInto(topic, payload, data)
			if err == nil {
				return nil
			}

			// Obtém o tipo do dado TData
			dataType := reflect.TypeOf(*data).Elem()

			// Tenta criar uma instância de protobuf apropriada para esse tipo
			protoInstance, err := c.protobufAdapter.CreateProtoInstance(dataType)
			if err != nil {
				// Se não conseguir criar uma instância específica, tenta um método mais genérico
				// Isso pode falhar se o tipo não for adequado
				return fmt.Errorf("não foi possível encontrar um tipo protobuf compatível: %v", err)
			}

			// Usa o deserializador com o tipo protobuf concreto
			err = c.registry.GetProtobufDeserializer().DeserializeInto(topic, payload, protoInstance)
			if err != nil {
				return fmt.Errorf("falha na deserialização protobuf: %v", err)
			}

			// Adapta a mensagem deserializada para o tipo alvo
			return c.protobufAdapter.AdaptDeserializedMessage(protoInstance, data)
		},
	}
}

// deserializeValue deserializa o payload de uma mensagem usando o deserializador apropriado.
// Seleciona o deserializador com base no tipo de deserialização especificado.
func (c *kafkaConsumer[TData]) deserializeValue(e *kafka.Message, deserialization enums.Deserialization, data *TData) error {
	if deserializerFunc, exists := c.deserializerMap[deserialization]; exists {
		return deserializerFunc(*e.TopicPartition.Topic, e.Value, data)
	}
	return errors.New("invalid deserializer")
}

// fillHeader preenche os metadados da mensagem com base nos cabeçalhos Kafka.
// Extrai também o correlationId, se disponível.
func (c *kafkaConsumer[TData]) fillHeader(message *message.Message[TData], headers []kafka.Header) {
	for _, header := range headers {
		message.Metadata[header.Key] = header.Value

		if strings.ToUpper(header.Key) == "CORRELATIONID" {
			correlationID, err := uuid.ParseBytes(header.Value)
			if err == nil {
				message.CorrelationId = correlationID
			}
		}
	}

	if message.CorrelationId.String() == "00000000-0000-0000-0000-000000000000" {
		message.CorrelationId = uuid.New()
	}
}

// Callback do cliente Kafka que recebe eventos de atribuição/revogação de partições.
// Gerencia atribuição e revogação de partições durante o rebalanceamento.
func rebalanceCallback(c *kafka.Consumer, event kafka.Event) error {
	switch ev := event.(type) {
	case kafka.AssignedPartitions:
		fmt.Printf("Rebalance Protocol: %s rebalance: %d new partition(s) assigned: %v\n",
			c.GetRebalanceProtocol(), len(ev.Partitions), ev.Partitions)
		err := c.Assign(ev.Partitions)
		if err != nil {
			return err
		}
	case kafka.RevokedPartitions:
		fmt.Printf("Rebalance Protocol: %s rebalance: %d partition(s) revoked: %v\n",
			c.GetRebalanceProtocol(), len(ev.Partitions), ev.Partitions)
		if c.AssignmentLost() {
			fmt.Fprintln(os.Stderr, "Assignment lost involuntarily, commit may fail")
		}
	default:
		fmt.Fprintf(os.Stderr, "Rebalance Protocol: %s Unxpected event type: %v\n", c.GetRebalanceProtocol(), event)
	}

	return nil
}

// commitMessage realiza o commit do offset da mensagem no Kafka.
// Faz commit a cada 10 mensagens para melhorar a performance.
func commitMessage(c *kafka.Consumer, topicPartition kafka.TopicPartition) error {
	if topicPartition.Offset%10 != 0 {
		return nil
	}

	_, err := c.Commit()

	if err != nil && err.(kafka.Error).Code() != kafka.ErrNoOffset {
		return err
	}

	return nil
}
