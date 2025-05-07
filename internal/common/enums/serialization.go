package enums

type Serialization int

const (
	JsonSerialization       Serialization = 1
	JsonSchemaSerialization Serialization = 2
	AvroSerialization       Serialization = 3
	ProtobufSerialization   Serialization = 4
)

type Deserialization int

const (
	JsonDeserialization       Deserialization = 1
	JsonSchemaDeserialization Deserialization = 2
	AvroDeserialization       Deserialization = 3
	ProtobufDeserialization   Deserialization = 4
)
