package enums

type SaslMechanims string

const (
	SASL_MECHANISM_NONE         SaslMechanims = "NONE"
	SASL_MECHANISM_GSSAPI       SaslMechanims = "GSSAPI"
	SASL_MECHANISM_PLAIN        SaslMechanims = "PLAIN"
	SASL_MECHANISM_SCRAM_SHA256 SaslMechanims = "SCRAM-SHA-256"
	SASL_MECHANISM_SCRAM_SHA512 SaslMechanims = "SCRAM-SHA-2512"
	SASL_MECHANISM_OAUTHBEARER  SaslMechanims = "OAUTHBEARER"
)
