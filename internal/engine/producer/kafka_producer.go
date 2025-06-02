package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Dieg657/kafka-toolkit-lib/internal/common/setup"
	"github.com/Dieg657/kafka-toolkit-lib/internal/engine/adapter"
	"github.com/Dieg657/kafka-toolkit-lib/pkg/common/constants"
	"github.com/Dieg657/kafka-toolkit-lib/pkg/common/enums"
	"github.com/Dieg657/kafka-toolkit-lib/pkg/common/ioc"
	"github.com/Dieg657/kafka-toolkit-lib/pkg/common/message"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// ==========================================================================
// Tipos e Propriedades Estáticas
// ==========================================================================

// Producer encapsula um produtor Kafka fortemente tipado.
// Gerencia a conexão com o Kafka e serialização de mensagens.
type kafkaProducer[TData any] struct {
	client          *kafka.Producer
	registry        setup.ISchemaRegistrySetup
	serializers     map[enums.Serialization]func(topic string, payload *TData) ([]byte, error)
	protobufAdapter *adapter.ProtobufAdapter // Adaptador para mensagens protobuf
}

// ==========================================================================
// Construtores
// ==========================================================================

// New inicializa um novo produtor Kafka.
// Configura os serializadores padrão (Avro e JSON) e obtém as dependências necessárias.
//
// Parâmetros:
//   - ctx: Contexto contendo as dependências e configurações
//
// Retorno:
//   - error: Erro caso a inicialização falhe
func NewKafkaProducer[TData any](ctx context.Context) (IKafkaProducer[TData], error) {
	container := ctx.Value(constants.IocKey).(ioc.IContainer)
	if container == nil {
		return nil, errors.New("IoC do Kafka não encontrado no contexto")
	}

	producerSetup, schemaRegistry := container.GetProducer()
	producer := &kafkaProducer[TData]{}
	producer.client = producerSetup.GetKafkaProducer()
	producer.registry = schemaRegistry
	producer.protobufAdapter = adapter.NewProtobufAdapter() // Inicializa o adaptador protobuf

	// Inicializa os serializadores suportados
	producer.initializeSerializers()
	return producer, nil
}

// ==========================================================================
// Métodos Públicos
// ==========================================================================

// Publish publica uma mensagem fortemente tipada em um tópico Kafka.
// Serializa a mensagem usando o formato especificado e gerencia metadados e chaves.
//
// Parâmetros:
//   - topic: Nome do tópico Kafka para publicação
//   - message: Mensagem tipada a ser publicada
//   - serialization: Formato de serialização a ser usado
//
// Retorno:
//   - error: Erro caso a publicação falhe
func (producer *kafkaProducer[TData]) Publish(topic string, message message.Message[TData], serialization enums.Serialization) error {
	var messageKey string

	if keyBytes, exists := message.Metadata["key"]; exists {
		messageKey = string(keyBytes)
	} else {
		messageKey = message.CorrelationId.String()
	}

	key, err := producer.serializeKey(messageKey)
	if err != nil {
		fmt.Println("Failed attempt to serialize key")
		return err
	}

	// Serializar apenas o campo Data, não a estrutura Message inteira
	payload, err := producer.serializePayload(topic, message.Data, serialization)
	if err != nil {
		fmt.Println("Failed attempt to serialize value")
		return err
	}

	// Convert metadata to Kafka headers
	var headers []kafka.Header
	// Add all metadata items as headers
	for key, value := range message.Metadata {
		if key != "key" {
			headers = append(headers, kafka.Header{
				Key:   key,
				Value: value,
			})
		}
	}

	kafkaMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          payload,
		Headers:        headers,
		Key:            key,
	}

	err = producer.client.Produce(kafkaMessage, nil)
	if err != nil {
		fmt.Printf("Failed when produce message: %v\n", err)
		return err
	}

	producer.client.Flush(10)
	return nil
}

// ProduceMessage publica uma mensagem Kafka pré-montada.
// Útil para casos onde o usuário precisa controle total sobre a configuração da mensagem.
//
// Parâmetros:
//   - msg: Mensagem Kafka pré-configurada
//
// Retorno:
//   - error: Erro caso a publicação falhe
func (producer *kafkaProducer[TData]) ProduceMessage(msg *kafka.Message) error {
	err := producer.client.Produce(msg, nil)
	if err != nil {
		return err
	}
	producer.client.Flush(10)
	return nil
}

// ==========================================================================
// Métodos Privados
// ==========================================================================

// initializeSerializers configura todos os serializadores suportados pela biblioteca.
// Centraliza a criação e registro dos serializadores padrão.
func (producer *kafkaProducer[TData]) initializeSerializers() {
	// Inicializa o mapa de serializadores
	producer.serializers = make(map[enums.Serialization]func(topic string, payload *TData) ([]byte, error))

	// Registra os serializadores disponíveis
	producer.serializers[enums.JsonSerialization] = func(topic string, payload *TData) ([]byte, error) {
		return json.Marshal(payload)
	}

	producer.serializers[enums.AvroSerialization] = func(topic string, payload *TData) ([]byte, error) {
		return producer.registry.GetAvroSerializer().Serialize(topic, payload)
	}

	producer.serializers[enums.JsonSchemaSerialization] = func(topic string, payload *TData) ([]byte, error) {
		return producer.registry.GetJsonSerializer().Serialize(topic, payload)
	}

	producer.serializers[enums.ProtobufSerialization] = func(topic string, payload *TData) ([]byte, error) {
		// Usa o adaptador para converter a mensagem para o formato esperado pelo serializador Protobuf
		adaptedPayload, err := producer.adaptProtobufPayload(payload)
		if err != nil {
			return nil, err
		}
		return producer.registry.GetProtobufSerializer().Serialize(topic, adaptedPayload)
	}
}

// adaptProtobufPayload adapta o payload para o formato correto do protobuf
func (producer *kafkaProducer[TData]) adaptProtobufPayload(payload *TData) (interface{}, error) {
	// Usamos o adaptador para converter a mensagem para o formato correto
	return producer.protobufAdapter.AdaptMessage(payload)
}

// serializePayload serializa o payload da mensagem usando o serializador apropriado.
// Seleciona o serializador com base no tipo de serialização especificado.
func (producer *kafkaProducer[TData]) serializePayload(topic string, payload TData, serialization enums.Serialization) ([]byte, error) {
	serializerFunc, exists := producer.serializers[serialization]
	if !exists {
		return nil, fmt.Errorf("serializador não registrado para o tipo: %v", serialization)
	}

	return serializerFunc(topic, &payload)
}

// serializeKey serializa a chave da mensagem para o formato usado pelo Kafka.
func (producer *kafkaProducer[TData]) serializeKey(key string) ([]byte, error) {
	return []byte(key), nil
}
