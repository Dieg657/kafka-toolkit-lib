package enums

type BasicAuthCredentialsSource string

const (
	BASIC_AUTH_CREDENTIALS_SOURCE_NONE         BasicAuthCredentialsSource = ""
	BASIC_AUTH_CREDENTIALS_SOURCE_USER_INFO    BasicAuthCredentialsSource = "USER_INFO"
	BASIC_AUTH_CREDENTIALS_SOURCE_SASL_INHERIT BasicAuthCredentialsSource = "SASL_INHERIT"
)
