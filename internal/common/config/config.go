package config

import (
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/enums"
)

// kafkaOptions implementa a interface IKafkaOptions
type kafkaOptions struct {
	Brokers          string
	GroupId          string
	Offset           enums.AutoOffsetReset
	UserName         string
	Password         string
	SecurityProtocol enums.SecurityProtocol
	SaslMechanism    enums.SaslMechanims
	SchemaRegistry   ISchemaRegistryOptions
	RequestTimeout   int
	ProducerPriority enums.ProducerOrderPriority
	ConsumerPriority enums.ConsumerOrderPriority
	build            bool
}

// SchemaRegistryOptions implementa a interface ISchemaRegistryOptions
type schemaRegistryOptions struct {
	Url                        string
	BasicAuthUser              string
	BasicAuthSecret            string
	AutoRegisterSchemas        bool
	RequestTimeout             int
	BasicAuthCredentialsSource enums.BasicAuthCredentialsSource
	build                      bool
}

func NewKafkaOptions() *kafkaOptions {
	options := &kafkaOptions{}
	options.build = false
	return options
}

func (k *kafkaOptions) SetBrokers(brokers string) {
	k.Brokers = brokers
}

func (k *kafkaOptions) SetGroupId(groupId string) {
	k.GroupId = groupId
}

func (k *kafkaOptions) SetOffset(offset enums.AutoOffsetReset) {
	k.Offset = offset
}

func (k *kafkaOptions) SetUserName(userName string) {
	k.UserName = userName
}

func (k *kafkaOptions) SetPassword(password string) {
	k.Password = password
}

func (k *kafkaOptions) SetSecurityProtocol(securityProtocol enums.SecurityProtocol) {
	k.SecurityProtocol = securityProtocol
}

func (k *kafkaOptions) SetSaslMechanism(saslMechanism enums.SaslMechanims) {
	k.SaslMechanism = saslMechanism
}

func (k *kafkaOptions) SetSchemaRegistry(schemaRegistry ISchemaRegistryOptions) {
	k.SchemaRegistry = schemaRegistry
}

func (k *kafkaOptions) SetRequestTimeout(requestTimeout int) {
	k.RequestTimeout = requestTimeout
}

func (k *kafkaOptions) SetProducerPriority(producerPriority enums.ProducerOrderPriority) {
	k.ProducerPriority = producerPriority
}

func (k *kafkaOptions) SetConsumerPriority(consumerPriority enums.ConsumerOrderPriority) {
	k.ConsumerPriority = consumerPriority
}

func (k *kafkaOptions) Validate() {
	if k.build {
		return
	}

	if k.Brokers == "" {
		panic("Brokers is required")
	}

	if k.GroupId == "" {
		panic("GroupId is required")
	}

	if k.Offset == "" {
		panic("Offset is required")
	}

	if k.UserName == "" && k.SecurityProtocol != enums.SECURITY_PROTOCOL_PLAINTEXT {
		panic("UserName is required")
	}

	if k.Password == "" && k.SecurityProtocol != enums.SECURITY_PROTOCOL_PLAINTEXT {
		panic("Password is required")
	}

	if k.SecurityProtocol == "" {
		panic("SecurityProtocol is required")
	}

	if k.SaslMechanism == "" {
		panic("SaslMechanism is required")
	}

	if k.SchemaRegistry.GetUrl() == "" {
		panic("SchemaRegistry.Url is required")
	}

	if k.SchemaRegistry.GetBasicAuthUser() == "" {
		panic("SchemaRegistry.BasicAuthUser is required")
	}

	if k.SchemaRegistry.GetBasicAuthSecret() == "" {
		panic("SchemaRegistry.de is required")
	}

	if k.SchemaRegistry.GetRequestTimeout() == 0 {
		panic("SchemaRegistry.RequestTimeout is required")
	}

	if k.SchemaRegistry.GetBasicAuthCredentialsSource() == "" {
		panic("SchemaRegistry.BasicAuthCredentialsSource is required")
	}

	if k.RequestTimeout == 0 {
		panic("RequestTimeout is required")
	}

	if k.ProducerPriority == "" {
		panic("ProducerPriority is required")
	}

	if k.ConsumerPriority == "" {
		panic("ConsumerPriority is required")
	}

	k.build = true
}

func (k *kafkaOptions) GetBrokers() string {
	return k.Brokers
}

func (k *kafkaOptions) GetGroupId() string {
	return k.GroupId
}

func (k *kafkaOptions) GetOffset() string {
	return string(k.Offset)
}

func (k *kafkaOptions) GetSaslMechanism() string {
	return string(k.SaslMechanism)
}

func (k *kafkaOptions) GetSecurityProtocol() string {
	return string(k.SecurityProtocol)
}

func (k *kafkaOptions) GetUserName() string {
	return k.UserName
}

func (k *kafkaOptions) GetPassword() string {
	return k.Password
}

func (k *kafkaOptions) GetRequestTimeout() int {
	return k.RequestTimeout
}

func (k *kafkaOptions) GetProducerPriority() string {
	return string(k.ProducerPriority)
}

func (k *kafkaOptions) GetConsumerPriority() string {
	return string(k.ConsumerPriority)
}

func (k *kafkaOptions) GetSchemaRegistry() ISchemaRegistryOptions {
	return k.SchemaRegistry
}

func NewSchemaRegistryOptions() *schemaRegistryOptions {
	options := &schemaRegistryOptions{}
	options.BasicAuthCredentialsSource = enums.BASIC_AUTH_CREDENTIALS_SOURCE_USER_INFO
	return options
}

func (s *schemaRegistryOptions) SetUrl(url string) {
	s.Url = url
}

func (s *schemaRegistryOptions) SetBasicAuthUser(basicAuthUser string) {
	s.BasicAuthUser = basicAuthUser
}

func (s *schemaRegistryOptions) SetAutoRegisterSchemas(autoRegisterSchemas bool) {
	s.AutoRegisterSchemas = autoRegisterSchemas
}

func (s *schemaRegistryOptions) SetRequestTimeout(requestTimeout int) {
	s.RequestTimeout = requestTimeout
}

func (s *schemaRegistryOptions) SetBasicAuthSecret(basicAuthSecret string) {
	s.BasicAuthSecret = basicAuthSecret
}

func (s *schemaRegistryOptions) SetBasicAuthCredentialsSource(basicAuthCredentialsSource enums.BasicAuthCredentialsSource) {
	s.BasicAuthCredentialsSource = basicAuthCredentialsSource
}

func (s *schemaRegistryOptions) Validate() {
	if s.build {
		return
	}

	if s.Url == "" {
		panic("Url is required")
	}

	// Validação condicional para autenticação do Schema Registry
	if s.BasicAuthCredentialsSource == enums.BASIC_AUTH_CREDENTIALS_SOURCE_USER_INFO {
		// Se a fonte for USER_INFO, as credenciais são obrigatórias
		if s.BasicAuthUser == "" {
			panic("BasicAuthUser is required when BasicAuthCredentialsSource is USER_INFO")
		}
		if s.BasicAuthSecret == "" {
			panic("BasicAuthSecret is required when BasicAuthCredentialsSource is USER_INFO")
		}
	} else if s.BasicAuthCredentialsSource == enums.BASIC_AUTH_CREDENTIALS_SOURCE_SASL_INHERIT {
		// Se for SASL_INHERIT, as credenciais SASL do Kafka serão usadas
		// Não é necessário validar BasicAuthUser e BasicAuthSecret
	} else if s.BasicAuthCredentialsSource == enums.BASIC_AUTH_CREDENTIALS_SOURCE_NONE {
		// Se for vazio, não exige autenticação
	} else {
		// Valor inválido para BasicAuthCredentialsSource
		panic("Invalid BasicAuthCredentialsSource")
	}

	if s.RequestTimeout == 0 {
		panic("RequestTimeout is required")
	}

	s.build = true
}

func (s *schemaRegistryOptions) GetUrl() string {
	return s.Url
}

func (s *schemaRegistryOptions) GetBasicAuthUser() string {
	return s.BasicAuthUser
}

func (s *schemaRegistryOptions) GetBasicAuthSecret() string {
	return s.BasicAuthSecret
}

func (s *schemaRegistryOptions) GetBasicAuthCredentialsSource() string {
	return string(s.BasicAuthCredentialsSource)
}

func (s *schemaRegistryOptions) GetAutoRegisterSchemas() bool {
	return s.AutoRegisterSchemas
}

func (s *schemaRegistryOptions) GetRequestTimeout() int {
	return s.RequestTimeout
}
