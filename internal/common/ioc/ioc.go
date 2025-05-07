package ioc

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/viper"

	"github.com/Dieg657/kafka-toolkit-lib/internal/common/config"
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/enums"
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/setup"
	consumerModule "github.com/Dieg657/kafka-toolkit-lib/internal/common/setup/consumer"
	producerModule "github.com/Dieg657/kafka-toolkit-lib/internal/common/setup/producer"
	registryModule "github.com/Dieg657/kafka-toolkit-lib/internal/common/setup/schema_registry"
)

// ==========================================================================
// Interfaces
// ==========================================================================

// IContainer define a interface pública para o container de dependências
type IContainer interface {
	// GetConsumer retorna interfaces para consumer e schema registry
	GetConsumer() (setup.IKafkaConsumerSetup, setup.ISchemaRegistrySetup)

	// GetProducer retorna interfaces para producer e schema registry
	GetProducer() (setup.IKafkaProducerSetup, setup.ISchemaRegistrySetup)

	// GetSchemaRegistry retorna a interface para o schema registry
	GetSchemaRegistry() setup.ISchemaRegistrySetup
}

// ==========================================================================
// Tipos e Propriedades
// ==========================================================================

// kafkaIoC implementação concreta do container de dependências
// Nota: agora é privado (letra minúscula) para esconder a implementação
type kafkaIoC struct {
	consumerSetup       setup.IKafkaConsumerSetup
	producerSetup       setup.IKafkaProducerSetup
	schemaRegistrySetup setup.ISchemaRegistrySetup
}

// ==========================================================================
// Factory
// ==========================================================================

// NewKafkaIoC cria e retorna uma nova instância de IContainer
func NewKafkaIoC(ctx context.Context) (IContainer, error) {
	ioc := &kafkaIoC{}
	err := ioc.initialize()
	if err != nil {
		return nil, err
	}
	return ioc, nil
}

// ==========================================================================
// Métodos Privados
// ==========================================================================

func (ioc *kafkaIoC) initialize() error {
	viper.AutomaticEnv()

	// Obter e validar o offset com valor padrão
	offsetValue := viper.GetString("KAFKA_AUTO_OFFSET_RESET")
	offset := enums.AutoOffsetReset(offsetValue)
	if offsetValue == "" || (offset != enums.OFFSET_RESET_EARLIEST &&
		offset != enums.OFFSET_RESET_LATEST &&
		offset != enums.OFFSET_RESET_ERROR &&
		offset != enums.OFFSET_RESET_SMALLEST &&
		offset != enums.OFFSET_RESET_BEGINNING &&
		offset != enums.OFFSET_RESET_LARGEST &&
		offset != enums.OFFSET_RESET_END) {
		offset = enums.OFFSET_RESET_LATEST // Valor padrão
	}

	// Obtem e valida o mecanismo SASL com valor padrão
	saslValue := viper.GetString("KAFKA_SASL_MECHANISM")
	saslMechanism := enums.SaslMechanims(saslValue)
	if saslValue == "" || (saslMechanism != enums.SASL_MECHANISM_NONE &&
		saslMechanism != enums.SASL_MECHANISM_GSSAPI &&
		saslMechanism != enums.SASL_MECHANISM_PLAIN &&
		saslMechanism != enums.SASL_MECHANISM_SCRAM_SHA256 &&
		saslMechanism != enums.SASL_MECHANISM_SCRAM_SHA512 &&
		saslMechanism != enums.SASL_MECHANISM_OAUTHBEARER) {
		saslMechanism = enums.SASL_MECHANISM_PLAIN // Valor padrão
	}

	// Obtem e valida o protocolo de segurança com valor padrão
	secProtocolValue := viper.GetString("KAFKA_SECURITY_PROTOCOL")
	securityProtocol := enums.SecurityProtocol(secProtocolValue)
	if secProtocolValue == "" || (securityProtocol != enums.SECURITY_PROTOCOL_PLAINTEXT &&
		securityProtocol != enums.SECURITY_PROTOCOL_SASL_PLAINTEXT &&
		securityProtocol != enums.SECURITY_PROTOCOL_SSL &&
		securityProtocol != enums.SECURITY_PROTOCOL_SASL_SSL) {
		securityProtocol = enums.SECURITY_PROTOCOL_SASL_PLAINTEXT // Valor padrão
	}

	// Obtem e valida a fonte de credenciais com valor padrão
	authSourceValue := viper.GetString("KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE")
	authSource := enums.BasicAuthCredentialsSource(authSourceValue)

	// Verificar username e password do schema registry
	schemaRegistryUsername := viper.GetString("KAFKA_SCHEMA_REGISTRY_USERNAME")
	schemaRegistryPassword := viper.GetString("KAFKA_SCHEMA_REGISTRY_PASSWORD")

	// Se não foi especificado um authSource mas foram informadas credenciais,
	// assume USER_INFO como padrão
	if authSourceValue == "" && schemaRegistryUsername != "" && schemaRegistryPassword != "" {
		authSource = enums.BASIC_AUTH_CREDENTIALS_SOURCE_USER_INFO
	} else if authSourceValue == "" || (authSource != enums.BASIC_AUTH_CREDENTIALS_SOURCE_NONE &&
		authSource != enums.BASIC_AUTH_CREDENTIALS_SOURCE_USER_INFO &&
		authSource != enums.BASIC_AUTH_CREDENTIALS_SOURCE_SASL_INHERIT) {
		authSource = enums.BASIC_AUTH_CREDENTIALS_SOURCE_USER_INFO // Valor padrão
	}

	// Obtem valores para outros campos com defaults quando necessário
	brokers := viper.GetString("KAFKA_BROKERS")
	if brokers == "" {
		return errors.New("KAFKA_BROKERS is not set")
	}

	groupId := viper.GetString("KAFKA_GROUPID")
	if groupId == "" {
		return errors.New("KAFKA_GROUPID is not set")
	}

	requestTimeout := viper.GetInt("KAFKA_TIMEOUT")
	if requestTimeout <= 0 {
		requestTimeout = 5000 // 5 segundos padrão
	}

	schemaRegistryUrl := viper.GetString("KAFKA_SCHEMA_REGISTRY_URL")
	if schemaRegistryUrl == "" {
		return errors.New("KAFKA_SCHEMA_REGISTRY_URL is not set")
	}

	// Verificar prioridades
	producerPriority := enums.ProducerOrderPriority(viper.GetString("KAFKA_PRODUCER_PRIORITY"))
	if producerPriority == "" || (producerPriority != enums.PRODUCER_ORDER_PRIORITY_ORDER &&
		producerPriority != enums.PRODUCER_ORDER_PRIORITY_BALANCED &&
		producerPriority != enums.PRODUCER_ORDER_PRIORITY_HIGH_PERFORMANCE) {
		producerPriority = enums.PRODUCER_ORDER_PRIORITY_BALANCED // Valor padrão
	}

	consumerPriority := enums.ConsumerOrderPriority(viper.GetString("KAFKA_CONSUMER_PRIORITY"))
	if consumerPriority == "" || (consumerPriority != enums.CONSUMER_ORDER_PRIORITY_ORDER &&
		consumerPriority != enums.CONSUMER_ORDER_PRIORITY_BALANCED &&
		consumerPriority != enums.CONSUMER_ORDER_PRIORITY_HIGH_PERFORMANCE &&
		consumerPriority != enums.CONSUMER_ORDER_PRIORITY_RISKY) {
		consumerPriority = enums.CONSUMER_ORDER_PRIORITY_BALANCED // Valor padrão
	}

	schemaRegistryOptions := config.NewSchemaRegistryOptions()
	schemaRegistryOptions.SetUrl(schemaRegistryUrl)

	// Configurar as credenciais e fonte de autenticação com base nos valores fornecidos
	if authSource == enums.BASIC_AUTH_CREDENTIALS_SOURCE_USER_INFO {
		// Se for USER_INFO, usamos as credenciais específicas do Schema Registry
		schemaRegistryOptions.SetBasicAuthUser(schemaRegistryUsername)
		schemaRegistryOptions.SetBasicAuthSecret(schemaRegistryPassword)
	} else if authSource == enums.BASIC_AUTH_CREDENTIALS_SOURCE_SASL_INHERIT {
		// Se for SASL_INHERIT, usamos as credenciais SASL do Kafka
		schemaRegistryOptions.SetBasicAuthUser(viper.GetString("KAFKA_USERNAME"))
		schemaRegistryOptions.SetBasicAuthSecret(viper.GetString("KAFKA_PASSWORD"))
	}
	// Se for NONE, não definimos credenciais

	schemaRegistryOptions.SetRequestTimeout(requestTimeout)
	schemaRegistryOptions.SetBasicAuthCredentialsSource(authSource)
	schemaRegistryOptions.Validate()

	kafkaOptions := config.NewKafkaOptions()
	kafkaOptions.SetBrokers(brokers)
	kafkaOptions.SetGroupId(groupId)
	kafkaOptions.SetOffset(offset)
	kafkaOptions.SetSaslMechanism(saslMechanism)
	kafkaOptions.SetSecurityProtocol(securityProtocol)
	kafkaOptions.SetUserName(viper.GetString("KAFKA_USERNAME"))
	kafkaOptions.SetPassword(viper.GetString("KAFKA_PASSWORD"))
	kafkaOptions.SetRequestTimeout(requestTimeout)
	kafkaOptions.SetProducerPriority(producerPriority)
	kafkaOptions.SetConsumerPriority(consumerPriority)
	kafkaOptions.SetSchemaRegistry(schemaRegistryOptions)
	kafkaOptions.Validate()

	// Usa as funções factory para obter as interfaces
	schemaRegistry, err := registryModule.NewSchemaRegistrySetup(kafkaOptions)
	if err != nil {
		fmt.Println("Error on initialize Schema Registry")
		return err
	}

	ioc.schemaRegistrySetup = schemaRegistry

	producerSetup, err := producerModule.NewKafkaProducerSetup(kafkaOptions)
	if err != nil {
		fmt.Println("Error on Setup Producer")
		return err
	}

	ioc.producerSetup = producerSetup

	consumerSetup, err := consumerModule.NewKafkaConsumerSetup(kafkaOptions)
	if err != nil {
		fmt.Println("Error on Setup Consumer")
		return err
	}

	ioc.consumerSetup = consumerSetup

	return nil
}

// ==========================================================================
// Métodos Públicos
// ==========================================================================

func (ioc *kafkaIoC) GetConsumer() (setup.IKafkaConsumerSetup, setup.ISchemaRegistrySetup) {
	return ioc.consumerSetup, ioc.schemaRegistrySetup
}

func (ioc *kafkaIoC) GetProducer() (setup.IKafkaProducerSetup, setup.ISchemaRegistrySetup) {
	return ioc.producerSetup, ioc.schemaRegistrySetup
}

func (ioc *kafkaIoC) GetSchemaRegistry() setup.ISchemaRegistrySetup {
	return ioc.schemaRegistrySetup
}
