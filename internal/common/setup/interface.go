package setup

import (
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/protobuf"
)

// ==========================================================================
// Interfaces
// ==========================================================================

// IKafkaConsumerSetup define a interface para configuração do consumidor Kafka
type IKafkaConsumerSetup interface {
	// GetKafkaConsumer retorna a instância do consumidor Kafka configurado
	GetKafkaConsumer() *kafka.Consumer
}

// IKafkaProducerSetup define a interface para configuração do produtor Kafka
type IKafkaProducerSetup interface {
	// GetKafkaProducer retorna a instância do produtor Kafka configurado
	GetKafkaProducer() *kafka.Producer
}

// ISchemaRegistrySetup define a interface para interação com o Schema Registry
type ISchemaRegistrySetup interface {
	// GetAvroSerializer retorna o serializador Avro específico
	GetAvroSerializer() *avro.SpecificSerializer

	// GetAvroDeserializer retorna o deserializador Avro específico
	GetAvroDeserializer() *avro.SpecificDeserializer

	// GetJsonSerializer retorna o serializador JSON
	GetJsonSerializer() *jsonschema.Serializer

	// GetJsonDeserializer retorna o deserializador JSON
	GetJsonDeserializer() *jsonschema.Deserializer

	// GetProtobufSerializer retorna o serializador Protobuf
	GetProtobufSerializer() *protobuf.Serializer

	// GetProtobufDeserializer retorna o deserializador Protobuf
	GetProtobufDeserializer() *protobuf.Deserializer
}

// ISerializer define a interface para serialização de mensagens em diferentes formatos
type ISerializer interface {
	// SerializeAvro serializa uma mensagem para o formato Avro
	SerializeAvro(topic string, msg any) ([]byte, error)

	// SerializeJson serializa uma mensagem para o formato JSON
	SerializeJson(topic string, msg any) ([]byte, error)

	// SerializeProtobuf serializa uma mensagem para o formato Protobuf
	SerializeProtobuf(topic string, msg any) ([]byte, error)
}

// IDeserializer define a interface para deserialização de mensagens de diferentes formatos
type IDeserializer interface {
	// DeserializeAvro deserializa uma mensagem do formato Avro
	DeserializeAvro(topic string, payload []byte) (any, error)

	// DeserializeJson deserializa uma mensagem do formato JSON
	DeserializeJson(topic string, payload []byte) (any, error)

	// DeserializeProtobuf deserializa uma mensagem do formato Protobuf
	DeserializeProtobuf(topic string, payload []byte) (any, error)
}
