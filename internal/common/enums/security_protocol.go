package enums

type SecurityProtocol string

const (
	SECURITY_PROTOCOL_PLAINTEXT      SecurityProtocol = "plaintext"
	SECURITY_PROTOCOL_SASL_PLAINTEXT SecurityProtocol = "sasl_plaintext"
	SECURITY_PROTOCOL_SSL            SecurityProtocol = "ssl"
	SECURITY_PROTOCOL_SASL_SSL       SecurityProtocol = "sasl_ssl"
)
