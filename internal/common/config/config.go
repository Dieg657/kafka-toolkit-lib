package config

import (
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/enums"
)

// ==========================================================================
// Tipos
// ==========================================================================

// kafkaOptions implementa a interface IKafkaOptions
type kafkaOptions struct {
	Brokers          string
	GroupId          string
	Offset           enums.AutoOffsetReset
	UserName         string
	Password         string
	SecurityProtocol enums.SecurityProtocol
	SaslMechanisms   enums.SaslMechanisms
	SchemaRegistry   ISchemaRegistryOptions
	RequestTimeout   int
	ProducerPriority enums.ProducerOrderPriority
	ConsumerPriority enums.ConsumerOrderPriority
	build            bool
}

// SchemaRegistryOptions implementa a interface ISchemaRegistryOptions
type schemaRegistryOptions struct {
	url                        string
	basicAuthUser              string
	basicAuthSecret            string
	autoRegisterSchemas        bool
	requestTimeout             int
	basicAuthCredentialsSource enums.BasicAuthCredentialsSource
	build                      bool
}

// ==========================================================================
// Construtores
// ==========================================================================

// NewKafkaOptions cria uma nova instância de configurações do Kafka
func NewKafkaOptions() *kafkaOptions {
	options := &kafkaOptions{}
	options.build = false
	return options
}

// NewSchemaRegistryOptions cria uma nova instância de configurações do Schema Registry
func NewSchemaRegistryOptions() *schemaRegistryOptions {
	options := &schemaRegistryOptions{}
	return options
}

// ==========================================================================
// Métodos KafkaOptions (Setters)
// ==========================================================================

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

func (k *kafkaOptions) SetSaslMechanisms(saslMechanisms enums.SaslMechanisms) {
	k.SaslMechanisms = saslMechanisms
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

// ==========================================================================
// Métodos KafkaOptions (Getters e validação)
// ==========================================================================

// Validate verifica se todos os campos obrigatórios foram configurados
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

	if k.SaslMechanisms == "" {
		panic("SaslMechanisms is required")
	}

	if k.RequestTimeout == 0 {
		k.RequestTimeout = 5000
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

func (k *kafkaOptions) GetSaslMechanisms() string {
	return string(k.SaslMechanisms)
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

// ==========================================================================
// Métodos SchemaRegistryOptions (Setters)
// ==========================================================================

func (s *schemaRegistryOptions) SetUrl(url string) {
	s.url = url
}

func (s *schemaRegistryOptions) SetBasicAuthUser(basicAuthUser string) {
	s.basicAuthUser = basicAuthUser
}

func (s *schemaRegistryOptions) SetAutoRegisterSchemas(autoRegisterSchemas bool) {
	s.autoRegisterSchemas = autoRegisterSchemas
}

func (s *schemaRegistryOptions) SetRequestTimeout(requestTimeout int) {
	s.requestTimeout = requestTimeout
}

func (s *schemaRegistryOptions) SetBasicAuthSecret(basicAuthSecret string) {
	s.basicAuthSecret = basicAuthSecret
}

func (s *schemaRegistryOptions) SetBasicAuthCredentialsSource(basicAuthCredentialsSource enums.BasicAuthCredentialsSource) {
	s.basicAuthCredentialsSource = basicAuthCredentialsSource
}

// ==========================================================================
// Métodos SchemaRegistryOptions (Getters e validação)
// ==========================================================================

// Validate verifica se todos os campos obrigatórios foram configurados
func (s *schemaRegistryOptions) Validate() {
	if s.build {
		return
	}

	// Validação do URL
	if s.url == "" {
		panic("Url is required")
	}

	// Validação do BasicAuthUser e BasicAuthSecret
	if s.basicAuthCredentialsSource != enums.BASIC_AUTH_CREDENTIALS_SOURCE_NONE && s.basicAuthUser == "" && s.basicAuthSecret == "" {
		panic("BasicAuthUser and BasicAuthSecret are required")
	}

	// Validação do timeout
	if s.requestTimeout == 0 {
		s.requestTimeout = 5000
	}

	s.build = true
}

func (s *schemaRegistryOptions) GetUrl() string {
	return s.url
}

func (s *schemaRegistryOptions) GetBasicAuthUser() string {
	return s.basicAuthUser
}

func (s *schemaRegistryOptions) GetBasicAuthSecret() string {
	return s.basicAuthSecret
}

func (s *schemaRegistryOptions) GetBasicAuthCredentialsSource() enums.BasicAuthCredentialsSource {
	return s.basicAuthCredentialsSource
}

func (s *schemaRegistryOptions) GetAutoRegisterSchemas() bool {
	return s.autoRegisterSchemas
}

func (s *schemaRegistryOptions) GetRequestTimeout() int {
	return s.requestTimeout
}
