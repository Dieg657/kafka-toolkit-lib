package enums

// ==========================================================================
// Protocolos de Segurança
// ==========================================================================

// SecurityProtocol define os protocolos de segurança suportados para comunicação com o Kafka
// Determina como os clientes se comunicam com os brokers em termos de segurança
type SecurityProtocol string

const (
	// SECURITY_PROTOCOL_PLAINTEXT protocolo sem criptografia ou autenticação
	// Recomendado apenas para ambientes de desenvolvimento isolados
	// NÃO UTILIZE em ambientes de produção expostos à rede pública
	SECURITY_PROTOCOL_PLAINTEXT SecurityProtocol = "plaintext"

	// SECURITY_PROTOCOL_SASL_PLAINTEXT protocolo com autenticação SASL mas sem criptografia
	// Fornece autenticação, mas as mensagens trafegam sem criptografia
	// Utilize apenas em redes internas seguras
	SECURITY_PROTOCOL_SASL_PLAINTEXT SecurityProtocol = "sasl_plaintext"

	// SECURITY_PROTOCOL_SSL protocolo com criptografia SSL/TLS
	// Fornece criptografia do tráfego, mas sem autenticação SASL
	// A autenticação pode ser feita via certificados de cliente
	SECURITY_PROTOCOL_SSL SecurityProtocol = "ssl"

	// SECURITY_PROTOCOL_SASL_SSL protocolo com autenticação SASL e criptografia SSL/TLS
	// Oferece máxima segurança com autenticação e criptografia
	// Recomendado para ambientes de produção e comunicação através de redes públicas
	SECURITY_PROTOCOL_SASL_SSL SecurityProtocol = "sasl_ssl"
)
