package enums

// ==========================================================================
// Mecanismos SASL
// ==========================================================================

// SaslMechanisms define os mecanismos SASL (Simple Authentication and Security Layer)
// suportados para autenticação com o Kafka
type SaslMechanisms string

const (
	// SASL_MECHANISM_NONE indica que nenhum mecanismo SASL será utilizado
	SASL_MECHANISM_NONE SaslMechanisms = "NONE"

	// SASL_MECHANISM_GSSAPI implementa autenticação Kerberos/GSSAPI
	// Adequado para ambientes com infraestrutura Kerberos existente
	SASL_MECHANISM_GSSAPI SaslMechanisms = "GSSAPI"

	// SASL_MECHANISM_PLAIN implementa autenticação simples com usuário e senha
	// Atenção: credenciais são transmitidas em texto claro, use com SSL/TLS
	SASL_MECHANISM_PLAIN SaslMechanisms = "PLAIN"

	// SASL_MECHANISM_SCRAM_SHA256 implementa SCRAM (Salted Challenge Response Authentication Mechanism)
	// com hash SHA-256, oferecendo melhor segurança que PLAIN
	SASL_MECHANISM_SCRAM_SHA256 SaslMechanisms = "SCRAM-SHA-256"

	// SASL_MECHANISM_SCRAM_SHA512 implementa SCRAM com hash SHA-512
	// Oferece segurança superior ao SHA-256, mas maior sobrecarga computacional
	SASL_MECHANISM_SCRAM_SHA512 SaslMechanisms = "SCRAM-SHA-512"

	// SASL_MECHANISM_OAUTHBEARER implementa autenticação baseada em OAuth 2.0
	// Adequado para integração com provedores de identidade externos
	SASL_MECHANISM_OAUTHBEARER SaslMechanisms = "OAUTHBEARER"
)
