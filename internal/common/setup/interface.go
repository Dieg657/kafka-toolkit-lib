package setup

import (
	"reflect"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/protobuf"
	"google.golang.org/protobuf/proto"
)

type IKafkaConsumerSetup interface {
	GetKafkaConsumer() *kafka.Consumer
}

type IKafkaProducerSetup interface {
	GetKafkaProducer() *kafka.Producer
}

type ISchemaRegistrySetup interface {
	GetAvroSerializer() *avro.SpecificSerializer
	GetAvroDeserializer() *avro.SpecificDeserializer
	GetJsonSerializer() *jsonschema.Serializer
	GetJsonDeserializer() *jsonschema.Deserializer
	GetProtobufSerializer() *protobuf.Serializer
	GetProtobufDeserializer() *protobuf.Deserializer

	// RegisterProtoType permite que o cliente registre um tipo protobuf a ser usado
	// durante a serialização/deserialização de mensagens.
	// targetType: é o tipo para o qual você pretende adaptar (por ex: reflect.TypeOf((*MeuTipo)(nil)).Elem())
	// protoMsgInstance: é uma instância do tipo protobuf (por ex: &pb.ProtoPayload{})
	RegisterProtoType(targetType reflect.Type, protoMsgInstance proto.Message) error
}

type ISerializer interface {
	SerializeAvro(topic string, msg any) ([]byte, error)
	SerializeJson(topic string, msg any) ([]byte, error)
	SerializeProtobuf(topic string, msg any) ([]byte, error)
}

type IDeserializer interface {
	DeserializeAvro(topic string, payload []byte) (any, error)
	DeserializeJson(topic string, payload []byte) (any, error)
	DeserializeProtobuf(topic string, payload []byte) (any, error)
}
