package enums

// Deserialization define os formatos de deserialização suportados
// Determina como os dados recebidos do Kafka serão convertidos para objetos
type Deserialization int

const (
	// JsonDeserialization representa deserialização em formato JSON puro
	// Converte dados JSON para objetos sem validação de schema
	// Contraparte da JsonSerialization
	JsonDeserialization Deserialization = 1

	// JsonSchemaDeserialization representa deserialização em formato JSON com schema
	// Converte dados JSON validando contra schema registrado
	// Contraparte da JsonSchemaSerialization
	JsonSchemaDeserialization Deserialization = 2

	// AvroDeserialization representa deserialização em formato Avro
	// Converte dados binários Avro para objetos usando schema registrado
	// Contraparte da AvroSerialization
	AvroDeserialization Deserialization = 3

	// ProtobufDeserialization representa deserialização em formato Protocol Buffers
	// Converte dados binários Protobuf para objetos usando schema registrado
	// Contraparte da ProtobufSerialization
	ProtobufDeserialization Deserialization = 4
)
