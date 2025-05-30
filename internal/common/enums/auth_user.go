package enums

// ==========================================================================
// Fontes de Credenciais
// ==========================================================================

// BasicAuthCredentialsSource define as fontes de credenciais para autenticação básica no Schema Registry
// Controla de onde virão as credenciais para autenticação com o Schema Registry
type BasicAuthCredentialsSource string

const (
	// BASIC_AUTH_CREDENTIALS_SOURCE_NONE não utiliza credenciais de autenticação
	// Utilizado quando o Schema Registry não requer autenticação
	// Equivalente a acesso anônimo ou quando a segurança é implementada em outra camada
	BASIC_AUTH_CREDENTIALS_SOURCE_NONE BasicAuthCredentialsSource = ""

	// BASIC_AUTH_CREDENTIALS_SOURCE_USER_INFO utiliza credenciais específicas para o Schema Registry
	// As credenciais são fornecidas explicitamente na configuração do Schema Registry
	// Permite que o Schema Registry tenha credenciais diferentes das usadas nos brokers Kafka
	BASIC_AUTH_CREDENTIALS_SOURCE_USER_INFO BasicAuthCredentialsSource = "USER_INFO"

	// BASIC_AUTH_CREDENTIALS_SOURCE_SASL_INHERIT herda as credenciais SASL do broker Kafka
	// Reutiliza o mesmo usuário e senha configurados para autenticação com os brokers Kafka
	// Útil quando o Schema Registry e Kafka compartilham o mesmo sistema de autenticação
	BASIC_AUTH_CREDENTIALS_SOURCE_SASL_INHERIT BasicAuthCredentialsSource = "SASL_INHERIT"
)
