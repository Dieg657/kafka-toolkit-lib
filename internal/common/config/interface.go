package config

type IKafkaOptions interface {
	GetBrokers() string
	GetGroupId() string
	GetOffset() string
	GetSaslMechanism() string
	GetSecurityProtocol() string
	GetUserName() string
	GetPassword() string
	GetRequestTimeout() int
	GetProducerPriority() string
	GetConsumerPriority() string
	GetSchemaRegistry() ISchemaRegistryOptions
	Validate()
}

type ISchemaRegistryOptions interface {
	GetUrl() string
	GetBasicAuthUser() string
	GetBasicAuthSecret() string
	GetBasicAuthCredentialsSource() string
	GetAutoRegisterSchemas() bool
	GetRequestTimeout() int
	Validate()
}
