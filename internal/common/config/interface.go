package config

import "github.com/Dieg657/kafka-toolkit-lib/internal/common/enums"

// ==========================================================================
// Interfaces
// ==========================================================================

// IKafkaOptions define a interface para as configurações do Kafka
type IKafkaOptions interface {
	// GetBrokers retorna o endereço dos brokers Kafka
	GetBrokers() string

	// GetGroupId retorna o ID do grupo de consumidores
	GetGroupId() string

	// GetOffset retorna a configuração de offset inicial
	GetOffset() string

	// GetSaslMechanisms retorna o mecanismo SASL configurado
	GetSaslMechanisms() string

	// GetSecurityProtocol retorna o protocolo de segurança configurado
	GetSecurityProtocol() string

	// GetUserName retorna o nome de usuário para autenticação
	GetUserName() string

	// GetPassword retorna a senha para autenticação
	GetPassword() string

	// GetRequestTimeout retorna o timeout para requisições em milissegundos
	GetRequestTimeout() int

	// GetProducerPriority retorna a prioridade configurada para o produtor
	GetProducerPriority() string

	// GetConsumerPriority retorna a prioridade configurada para o consumidor
	GetConsumerPriority() string

	// GetSchemaRegistry retorna as configurações do Schema Registry
	GetSchemaRegistry() ISchemaRegistryOptions

	// Validate valida as configurações, garantindo que todos os valores obrigatórios estão presentes
	Validate()
}

// ISchemaRegistryOptions define a interface para as configurações do Schema Registry
type ISchemaRegistryOptions interface {
	// GetUrl retorna a URL do Schema Registry
	GetUrl() string

	// GetBasicAuthUser retorna o usuário para autenticação básica
	GetBasicAuthUser() string

	// GetBasicAuthSecret retorna a senha para autenticação básica
	GetBasicAuthSecret() string

	// GetBasicAuthCredentialsSource retorna a fonte de credenciais para autenticação
	GetBasicAuthCredentialsSource() enums.BasicAuthCredentialsSource

	// GetAutoRegisterSchemas retorna se schemas devem ser registrados automaticamente
	GetAutoRegisterSchemas() bool

	// GetRequestTimeout retorna o timeout para requisições em milissegundos
	GetRequestTimeout() int

	// Validate valida as configurações, garantindo que todos os valores obrigatórios estão presentes
	Validate()
}
