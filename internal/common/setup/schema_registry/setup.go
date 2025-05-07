package setup

import (
	"fmt"
	"os"
	"reflect"

	"github.com/Dieg657/kafka-toolkit-lib/internal/common/config"
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/setup"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/protobuf"
	"google.golang.org/protobuf/proto"
)

//
// Tipos e Estruturas
//

// schemaRegistrySetup implementação concreta privada
type schemaRegistrySetup struct {
	schemaRegistry           schemaregistry.Client
	avroSpecificSerializer   *avro.SpecificSerializer
	avroSpecificDeserializer *avro.SpecificDeserializer
	jsonSerializer           *jsonschema.Serializer
	jsonDeserializer         *jsonschema.Deserializer
	protobufSerializer       *protobuf.Serializer
	protobufDeserializer     *protobuf.Deserializer
	protoTypes               map[string]proto.Message // Registro de tipos protobuf para adaptação
}

// NewSchemaRegistrySetup cria uma nova instância da interface ISchemaRegistrySetup
func NewSchemaRegistrySetup(options config.IKafkaOptions) (setup.ISchemaRegistrySetup, error) {
	registry := &schemaRegistrySetup{
		protoTypes: make(map[string]proto.Message),
	}
	err := registry.New(options)
	if err != nil {
		return nil, err
	}
	return registry, nil
}

//
// Construtores
//

// New inicializa a configuração do schema registry com as opções fornecidas.
//
// Parâmetros:
//   - options: Configurações para conexão ao schema registry
//
// Retorno:
//   - error: Erro caso a inicialização falhe
func (registry *schemaRegistrySetup) New(options config.IKafkaOptions) error {
	configuration := &schemaregistry.Config{}

	configuration.SchemaRegistryURL = options.GetSchemaRegistry().GetUrl()
	configuration.RequestTimeoutMs = options.GetSchemaRegistry().GetRequestTimeout()
	configuration.ConnectionTimeoutMs = 5000

	if string(options.GetSchemaRegistry().GetBasicAuthCredentialsSource()) != "" {
		configuration.BasicAuthCredentialsSource = options.GetSchemaRegistry().GetBasicAuthCredentialsSource()
		configuration.BasicAuthUserInfo = fmt.Sprintf("%s:%s", options.GetSchemaRegistry().GetBasicAuthUser(), options.GetSchemaRegistry().GetBasicAuthSecret())
	}

	schemaRegistry, err := schemaregistry.NewClient(configuration)

	if err != nil {
		return err
	}

	registry.schemaRegistry = schemaRegistry
	registry.setSerializers()
	registry.setDeserializers()

	return nil
}

//
// Métodos Públicos
//

// GetAvroSerializer retorna o serializador Avro específico.
func (sc *schemaRegistrySetup) GetAvroSerializer() *avro.SpecificSerializer {
	return sc.avroSpecificSerializer
}

// GetAvroDeserializer retorna o deserializador Avro específico.
func (sc *schemaRegistrySetup) GetAvroDeserializer() *avro.SpecificDeserializer {
	return sc.avroSpecificDeserializer
}

// GetJsonSerializer retorna o serializador JSON.
func (sc *schemaRegistrySetup) GetJsonSerializer() *jsonschema.Serializer {
	return sc.jsonSerializer
}

// GetJsonDeserializer retorna o deserializador JSON.
func (sc *schemaRegistrySetup) GetJsonDeserializer() *jsonschema.Deserializer {
	return sc.jsonDeserializer
}

// GetProtobufSerializer retorna o serializador Protobuf.
func (sc *schemaRegistrySetup) GetProtobufSerializer() *protobuf.Serializer {
	return sc.protobufSerializer
}

// GetProtobufDeserializer retorna o deserializador Protobuf.
func (sc *schemaRegistrySetup) GetProtobufDeserializer() *protobuf.Deserializer {
	return sc.protobufDeserializer
}

// RegisterProtoType permite que o cliente registre um tipo protobuf a ser usado
// durante a serialização/deserialização de mensagens.
func (registry *schemaRegistrySetup) RegisterProtoType(targetType reflect.Type, protoMsgInstance proto.Message) error {
	if protoMsgInstance == nil {
		return fmt.Errorf("a instância protobuf não pode ser nula")
	}

	registry.protoTypes[targetType.String()] = protoMsgInstance
	return nil
}

// GetProtoTypeForTarget retorna um tipo protobuf registrado para o tipo de destino especificado
func (registry *schemaRegistrySetup) GetProtoTypeForTarget(targetType reflect.Type) (proto.Message, bool) {
	protoType, exists := registry.protoTypes[targetType.String()]
	return protoType, exists
}

//
// Métodos Privados
//

// setSerializers inicializa todos os serializadores.
func (sc *schemaRegistrySetup) setSerializers() {
	sc.configAvroSerializer()
	sc.configJsonSerializer()
	sc.configProtobufSerializer()
}

// setDeserializers inicializa todos os deserializadores.
func (sc *schemaRegistrySetup) setDeserializers() {
	sc.configAvroDeserializer()
	sc.configJsonDeserializer()
	sc.configProtobufDeserializer()
}

// configAvroSerializer configura o serializador Avro.
func (registry *schemaRegistrySetup) configAvroSerializer() {
	configSerializer := avro.NewSerializerConfig()
	configSerializer.AutoRegisterSchemas = true

	avroValueSerializer, errValue := avro.NewSpecificSerializer(registry.schemaRegistry, serde.ValueSerde, configSerializer)

	if errValue != nil {
		fmt.Println("Error on create AVRO Schema Serializer")
		os.Exit(-1)
	}

	registry.avroSpecificSerializer = avroValueSerializer
	// Não deve retornar nada
}

// configAvroDeserializer configura o deserializador Avro.
func (sc *schemaRegistrySetup) configAvroDeserializer() {
	configDeserializer := avro.NewDeserializerConfig()
	avroValueDeserializer, errValue := avro.NewSpecificDeserializer(sc.schemaRegistry, serde.ValueSerde, configDeserializer)

	if errValue != nil {
		fmt.Println("Error on create AVRO Schema Deserializer")
		os.Exit(-1)
	}

	sc.avroSpecificDeserializer = avroValueDeserializer
}

// configJsonSerializer configura o serializador JSON.
func (sc *schemaRegistrySetup) configJsonSerializer() {
	configSerializer := jsonschema.NewSerializerConfig()
	configSerializer.AutoRegisterSchemas = true
	jsonValueSerializer, errValue := jsonschema.NewSerializer(sc.schemaRegistry, serde.ValueSerde, configSerializer)

	if errValue != nil {
		fmt.Println("Error on create JSON Schema Serializer")
		os.Exit(-1)
	}

	sc.jsonSerializer = jsonValueSerializer
}

// configJsonDeserializer configura o deserializador JSON.
func (sc *schemaRegistrySetup) configJsonDeserializer() {
	configDeserializer := jsonschema.NewDeserializerConfig()
	jsonValueDeserializer, errValue := jsonschema.NewDeserializer(sc.schemaRegistry, serde.ValueSerde, configDeserializer)

	if errValue != nil {
		fmt.Println("Error on create JSON Schema Deserializer")
		os.Exit(-1)
	}

	sc.jsonDeserializer = jsonValueDeserializer
}

func (sc *schemaRegistrySetup) configProtobufSerializer() {
	configSerializer := protobuf.NewSerializerConfig()
	configSerializer.AutoRegisterSchemas = true
	protobufValueSerializer, errValue := protobuf.NewSerializer(sc.schemaRegistry, serde.ValueSerde, configSerializer)

	if errValue != nil {
		fmt.Println("Error on create Protobuf Schema Serializer")
		os.Exit(-1)
	}

	sc.protobufSerializer = protobufValueSerializer
}

func (sc *schemaRegistrySetup) configProtobufDeserializer() {
	configDeserializer := protobuf.NewDeserializerConfig()
	protobufValueDeserializer, errValue := protobuf.NewDeserializer(sc.schemaRegistry, serde.ValueSerde, configDeserializer)

	if errValue != nil {
		fmt.Println("Error on create Protobuf Schema Deserializer")
		os.Exit(-1)
	}

	sc.protobufDeserializer = protobufValueDeserializer
}
