package setup

import (
	"fmt"

	"github.com/Dieg657/kafka-toolkit-lib/internal/common/config"
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/enums"
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/setup"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/avro"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/jsonschema"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/protobuf"
	"google.golang.org/protobuf/proto"
)

// ==========================================================================
// Tipos e Propriedades
// ==========================================================================

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
	err                      error
}

// ==========================================================================
// Factory
// ==========================================================================

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

// ==========================================================================
// Construtores e Inicialização
// ==========================================================================

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

	if options.GetSchemaRegistry().GetBasicAuthCredentialsSource() != enums.BASIC_AUTH_CREDENTIALS_SOURCE_NONE {
		configuration.BasicAuthCredentialsSource = string(options.GetSchemaRegistry().GetBasicAuthCredentialsSource())
		configuration.BasicAuthUserInfo = fmt.Sprintf("%s:%s", options.GetSchemaRegistry().GetBasicAuthUser(), options.GetSchemaRegistry().GetBasicAuthSecret())
	}

	schemaRegistry, err := schemaregistry.NewClient(configuration)

	if err != nil {
		return err
	}

	registry.schemaRegistry = schemaRegistry
	registry.setSerializers()
	registry.setDeserializers()

	// Verifica se algum erro ocorreu durante a configuração
	if registry.err != nil {
		return registry.err
	}

	return nil
}

// ==========================================================================
// Métodos Públicos
// ==========================================================================

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

// // RegisterProtoType permite que o cliente registre um tipo protobuf a ser usado
// // durante a serialização/deserialização de mensagens.
// func (registry *schemaRegistrySetup) RegisterProtoType(targetType reflect.Type, protoMsgInstance proto.Message) error {
// if protoMsgInstance == nil {
// return fmt.Errorf("a instância protobuf não pode ser nula")
// }

// registry.protoTypes[targetType.String()] = protoMsgInstance
// return nil
// }

// // GetProtoTypeForTarget retorna um tipo protobuf registrado para o tipo de destino especificado
// func (registry *schemaRegistrySetup) GetProtoTypeForTarget(targetType reflect.Type) (proto.Message, bool) {
// protoType, exists := registry.protoTypes[targetType.String()]
// return protoType, exists
// }

// ==========================================================================
// Métodos Privados
// ==========================================================================

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

	avroValueSerializer, err := avro.NewSpecificSerializer(registry.schemaRegistry, serde.ValueSerde, configSerializer)
	if err != nil {
		registry.avroSpecificSerializer = nil
		registry.err = fmt.Errorf("falha ao criar serializador AVRO: %w", err)
		return
	}

	registry.avroSpecificSerializer = avroValueSerializer
}

// configAvroDeserializer configura o deserializador Avro.
func (sc *schemaRegistrySetup) configAvroDeserializer() {
	configDeserializer := avro.NewDeserializerConfig()
	avroValueDeserializer, err := avro.NewSpecificDeserializer(sc.schemaRegistry, serde.ValueSerde, configDeserializer)
	if err != nil {
		sc.avroSpecificDeserializer = nil
		sc.err = fmt.Errorf("falha ao criar deserializador AVRO: %w", err)
		return
	}

	sc.avroSpecificDeserializer = avroValueDeserializer
}

// configJsonSerializer configura o serializador JSON.
func (sc *schemaRegistrySetup) configJsonSerializer() {
	configSerializer := jsonschema.NewSerializerConfig()
	configSerializer.AutoRegisterSchemas = true
	jsonValueSerializer, err := jsonschema.NewSerializer(sc.schemaRegistry, serde.ValueSerde, configSerializer)
	if err != nil {
		sc.jsonSerializer = nil
		sc.err = fmt.Errorf("falha ao criar serializador JSON: %w", err)
		return
	}

	sc.jsonSerializer = jsonValueSerializer
}

// configJsonDeserializer configura o deserializador JSON.
func (sc *schemaRegistrySetup) configJsonDeserializer() {
	configDeserializer := jsonschema.NewDeserializerConfig()
	jsonValueDeserializer, err := jsonschema.NewDeserializer(sc.schemaRegistry, serde.ValueSerde, configDeserializer)
	if err != nil {
		sc.jsonDeserializer = nil
		sc.err = fmt.Errorf("falha ao criar deserializador JSON: %w", err)
		return
	}

	sc.jsonDeserializer = jsonValueDeserializer
}

// configProtobufSerializer configura o serializador Protobuf.
func (sc *schemaRegistrySetup) configProtobufSerializer() {
	configSerializer := protobuf.NewSerializerConfig()
	configSerializer.AutoRegisterSchemas = true
	protobufValueSerializer, err := protobuf.NewSerializer(sc.schemaRegistry, serde.ValueSerde, configSerializer)
	if err != nil {
		sc.protobufSerializer = nil
		sc.err = fmt.Errorf("falha ao criar serializador Protobuf: %w", err)
		return
	}

	sc.protobufSerializer = protobufValueSerializer
}

// configProtobufDeserializer configura o deserializador Protobuf.
func (sc *schemaRegistrySetup) configProtobufDeserializer() {
	configDeserializer := protobuf.NewDeserializerConfig()
	protobufValueDeserializer, err := protobuf.NewDeserializer(sc.schemaRegistry, serde.ValueSerde, configDeserializer)
	if err != nil {
		sc.protobufDeserializer = nil
		sc.err = fmt.Errorf("falha ao criar deserializador Protobuf: %w", err)
		return
	}

	sc.protobufDeserializer = protobufValueDeserializer
}
