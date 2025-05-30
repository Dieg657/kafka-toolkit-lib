package enums

// ==========================================================================
// Tipos de Serialização
// ==========================================================================

// Serialization define os formatos de serialização suportados
// Determina como os dados serão convertidos antes de serem enviados ao Kafka
type Serialization int

const (
	// JsonSerialization representa serialização em formato JSON puro
	// Dados são convertidos para JSON sem validação de schema
	// Adequado para mensagens simples ou quando flexibilidade é mais importante que validação
	JsonSerialization Serialization = 1

	// JsonSchemaSerialization representa serialização em formato JSON com schema
	// Dados JSON são validados contra um schema JSON definido
	// Oferece validação estrutural mantendo a legibilidade do JSON
	JsonSchemaSerialization Serialization = 2

	// AvroSerialization representa serialização em formato Avro
	// Formato binário compacto com schemas definidos
	// Oferece melhor desempenho, economia de espaço e evolução de schema controlada
	AvroSerialization Serialization = 3

	// ProtobufSerialization representa serialização em formato Protocol Buffers
	// Formato binário compacto e eficiente baseado em schemas (.proto)
	// Excelente para alta performance, baixa latência e forte tipagem
	ProtobufSerialization Serialization = 4
)
