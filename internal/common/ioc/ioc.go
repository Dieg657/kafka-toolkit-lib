package ioc

import (
	"fmt"
	"strings"
	"sync"

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

	// GetConsumerPriority retorna a prioridade do consumidor
	GetConsumerPriority() enums.ConsumerOrderPriority

	// GetProducerPriority retorna a prioridade do produtor
	GetProducerPriority() enums.ProducerOrderPriority
}

// ==========================================================================
// Tipos e Propriedades
// ==========================================================================

var (
	once         sync.Once
	iocContainer IContainer
)

// kafkaIoC implementação concreta do container de dependências
// Nota: agora é privado (letra minúscula) para esconder a implementação
type kafkaIoC struct {
	consumerSetup       setup.IKafkaConsumerSetup
	producerSetup       setup.IKafkaProducerSetup
	schemaRegistrySetup setup.ISchemaRegistrySetup
	consumerPriority    enums.ConsumerOrderPriority
	producerPriority    enums.ProducerOrderPriority
}

// Mapas para normalização dos parâmetros
var (
	securityProtocolMap = map[string]string{
		"PLAINTEXT":      string(enums.SECURITY_PROTOCOL_PLAINTEXT),
		"SASL_PLAINTEXT": string(enums.SECURITY_PROTOCOL_SASL_PLAINTEXT),
		"SSL":            string(enums.SECURITY_PROTOCOL_SSL),
		"SASL_SSL":       string(enums.SECURITY_PROTOCOL_SASL_SSL),
	}
	saslMechanismMap = map[string]string{
		"PLAIN":         string(enums.SASL_MECHANISM_PLAIN),
		"SCRAM-SHA-256": string(enums.SASL_MECHANISM_SCRAM_SHA256),
		"SCRAM-SHA-512": string(enums.SASL_MECHANISM_SCRAM_SHA512),
		"GSSAPI":        string(enums.SASL_MECHANISM_GSSAPI),
		"OAUTHBEARER":   string(enums.SASL_MECHANISM_OAUTHBEARER),
		"NONE":          string(enums.SASL_MECHANISM_NONE),
	}
	basicAuthCredentialsSourceMap = map[string]enums.BasicAuthCredentialsSource{
		"USER_INFO":    enums.BASIC_AUTH_CREDENTIALS_SOURCE_USER_INFO,
		"SASL_INHERIT": enums.BASIC_AUTH_CREDENTIALS_SOURCE_SASL_INHERIT,
		"NONE":         enums.BASIC_AUTH_CREDENTIALS_SOURCE_NONE,
	}
	autoOffsetResetMap = map[string]string{
		"ERROR":     string(enums.OFFSET_RESET_ERROR),
		"SMALLEST":  string(enums.OFFSET_RESET_SMALLEST),
		"EARLIEST":  string(enums.OFFSET_RESET_EARLIEST),
		"BEGINNING": string(enums.OFFSET_RESET_BEGINNING),
		"LARGEST":   string(enums.OFFSET_RESET_LARGEST),
		"LATEST":    string(enums.OFFSET_RESET_LATEST),
		"END":       string(enums.OFFSET_RESET_END),
	}
	producerPriorityMap = map[string]string{
		"ORDER":            string(enums.PRODUCER_ORDER_PRIORITY_ORDER),
		"BALANCED":         string(enums.PRODUCER_ORDER_PRIORITY_BALANCED),
		"HIGH_PERFORMANCE": string(enums.PRODUCER_ORDER_PRIORITY_HIGH_PERFORMANCE),
	}
	consumerPriorityMap = map[string]string{
		"ORDER":            string(enums.CONSUMER_ORDER_PRIORITY_ORDER),
		"BALANCED":         string(enums.CONSUMER_ORDER_PRIORITY_BALANCED),
		"HIGH_PERFORMANCE": string(enums.CONSUMER_ORDER_PRIORITY_HIGH_PERFORMANCE),
		"RISKY":            string(enums.CONSUMER_ORDER_PRIORITY_RISKY),
	}
)

// Função utilitária para mapear security protocol amigável para valor Kafka
func mapSecurityProtocolToKafka(value string) string {
	mapped, ok := securityProtocolMap[strings.ToUpper(value)]
	if ok {
		return mapped
	}
	return string(enums.SECURITY_PROTOCOL_PLAINTEXT)
}

// Função utilitária para mapear sasl.mechanism amigável para valor Kafka
func mapSaslMechanismToKafka(value string) string {
	mapped, ok := saslMechanismMap[strings.ToUpper(value)]
	if ok {
		return mapped
	}
	return string(enums.SASL_MECHANISM_PLAIN)
}

// Função utilitária para mapear basicAuthCredentialsSource amigável para valor Kafka
func mapBasicAuthCredentialsSourceToKafka(value string) enums.BasicAuthCredentialsSource {
	mapped, ok := basicAuthCredentialsSourceMap[strings.ToUpper(value)]
	if ok {
		return mapped
	}
	return enums.BASIC_AUTH_CREDENTIALS_SOURCE_USER_INFO
}

// Função utilitária para mapear autoOffsetReset amigável para valor Kafka
func mapAutoOffsetResetToKafka(value string) string {
	mapped, ok := autoOffsetResetMap[strings.ToUpper(value)]
	if ok {
		return mapped
	}
	return string(enums.OFFSET_RESET_LATEST)
}

// Função utilitária para mapear producerPriority amigável para valor Kafka
func mapProducerPriorityToKafka(value string) string {
	mapped, ok := producerPriorityMap[strings.ToUpper(value)]
	if ok {
		return mapped
	}
	return string(enums.PRODUCER_ORDER_PRIORITY_ORDER)
}

// Função utilitária para mapear consumerPriority amigável para valor Kafka
func mapConsumerPriorityToKafka(value string) string {
	mapped, ok := consumerPriorityMap[strings.ToUpper(value)]
	if ok {
		return mapped
	}
	return string(enums.CONSUMER_ORDER_PRIORITY_ORDER)
}

// ==========================================================================
// Factory
// ==========================================================================

// GetKafkaIoC retorna uma instância única do container de dependências
func GetKafkaIoC() (IContainer, error) {
	once.Do(func() {
		var err error
		iocContainer, err = newKafkaIoC()
		if err != nil {
			panic(err)
		}
	})
	return iocContainer, nil
}

// newKafkaIoC cria e retorna uma nova instância de IContainer
func newKafkaIoC() (IContainer, error) {
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

	// Obtem e valida o offset com valor default: LATEST em caso de não ser informado
	offsetValue := viper.GetString("KAFKA_AUTO_OFFSET_RESET")
	offset := enums.AutoOffsetReset(mapAutoOffsetResetToKafka(offsetValue))

	// Obtem o mecanismo SASL com valor default: PLAIN em caso de não ser informado
	saslValue := viper.GetString("KAFKA_SASL_MECHANISM")
	saslMechanisms := enums.SaslMechanisms(mapSaslMechanismToKafka(saslValue))

	// Obtem o protocolo de segurança com valor default: SASL_PLAINTEXT em caso de não ser informado
	secProtocolValue := viper.GetString("KAFKA_SECURITY_PROTOCOL")
	securityProtocol := enums.SecurityProtocol(mapSecurityProtocolToKafka(secProtocolValue))

	// Configura a prioridade do produtor - Valor default: ORDER em caso de não ser informado
	producerPriority := enums.ProducerOrderPriority(mapProducerPriorityToKafka(viper.GetString("KAFKA_PRODUCER_PRIORITY")))

	// Configura a prioridade do consumidor - Valor default: ORDER em caso de não ser informado
	consumerPriority := enums.ConsumerOrderPriority(mapConsumerPriorityToKafka(viper.GetString("KAFKA_CONSUMER_PRIORITY")))

	schemaRegistryOptions := config.NewSchemaRegistryOptions()
	schemaRegistryOptions.SetUrl(viper.GetString("KAFKA_SCHEMA_REGISTRY_URL"))

	// Obtem e define a fonte de credenciais com valor default: USER_INFO em caso de não ser informado
	authSource := mapBasicAuthCredentialsSourceToKafka(viper.GetString("KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE"))
	schemaRegistryOptions.SetBasicAuthCredentialsSource(authSource)

	// Configura as credenciais e fonte de autenticação com base nos valores fornecidos
	if enums.BasicAuthCredentialsSource(schemaRegistryOptions.GetBasicAuthCredentialsSource()) == enums.BASIC_AUTH_CREDENTIALS_SOURCE_USER_INFO {
		// Se for USER_INFO, usamos as credenciais específicas do Schema Registry
		schemaRegistryOptions.SetBasicAuthUser(viper.GetString("KAFKA_SCHEMA_REGISTRY_USERNAME"))
		schemaRegistryOptions.SetBasicAuthSecret(viper.GetString("KAFKA_SCHEMA_REGISTRY_PASSWORD"))
	} else if enums.BasicAuthCredentialsSource(schemaRegistryOptions.GetBasicAuthCredentialsSource()) == enums.BASIC_AUTH_CREDENTIALS_SOURCE_SASL_INHERIT {
		// Se for SASL_INHERIT, usamos as credenciais SASL do Kafka
		schemaRegistryOptions.SetBasicAuthUser(viper.GetString("KAFKA_USERNAME"))
		schemaRegistryOptions.SetBasicAuthSecret(viper.GetString("KAFKA_PASSWORD"))
	}

	// Configura o timeout para requisições com valor default: 5000 em caso de não ser informado
	schemaRegistryOptions.SetRequestTimeout(viper.GetInt("KAFKA_TIMEOUT"))

	// Valida as configurações
	schemaRegistryOptions.Validate()

	kafkaOptions := config.NewKafkaOptions()

	// Configura os brokers
	kafkaOptions.SetBrokers(viper.GetString("KAFKA_BROKERS"))

	// Configura o grupo de consumidores
	kafkaOptions.SetGroupId(viper.GetString("KAFKA_GROUPID"))

	// Configura o offset
	kafkaOptions.SetOffset(offset)

	// Configura o mecanismo SASL
	kafkaOptions.SetSaslMechanisms(saslMechanisms)

	// Configura o protocolo de segurança
	kafkaOptions.SetSecurityProtocol(securityProtocol)

	// Configura o nome do usuário
	kafkaOptions.SetUserName(viper.GetString("KAFKA_USERNAME"))

	// Configura a senha do usuário
	kafkaOptions.SetPassword(viper.GetString("KAFKA_PASSWORD"))

	// Configura a prioridade do produtor
	kafkaOptions.SetProducerPriority(producerPriority)

	// Configura a prioridade do consumidor
	kafkaOptions.SetConsumerPriority(consumerPriority)

	// Configura o schema registry
	kafkaOptions.SetSchemaRegistry(schemaRegistryOptions)

	// Valida as configurações
	kafkaOptions.Validate()

	// Cria o schema registry
	schemaRegistry, err := registryModule.NewSchemaRegistrySetup(kafkaOptions)
	if err != nil {
		return fmt.Errorf("falha ao inicializar Schema Registry: %w", err)
	}

	ioc.schemaRegistrySetup = schemaRegistry

	// Cria o produtor
	producerSetup, err := producerModule.NewKafkaProducerSetup(kafkaOptions)
	if err != nil {
		return fmt.Errorf("falha ao configurar Producer: %w", err)
	}

	ioc.producerSetup = producerSetup

	// Cria o consumidor
	consumerSetup, err := consumerModule.NewKafkaConsumerSetup(kafkaOptions)
	if err != nil {
		return fmt.Errorf("falha ao configurar Consumer: %w", err)
	}

	// Atribui as interfaces e prioridades
	ioc.consumerSetup = consumerSetup
	ioc.consumerPriority = consumerPriority
	ioc.producerPriority = producerPriority

	return nil
}

// ==========================================================================
// Métodos Públicos
// ==========================================================================

// GetConsumer retorna as interfaces para consumer e schema registry
func (ioc *kafkaIoC) GetConsumer() (setup.IKafkaConsumerSetup, setup.ISchemaRegistrySetup) {
	return ioc.consumerSetup, ioc.schemaRegistrySetup
}

// GetProducer retorna as interfaces para producer e schema registry
func (ioc *kafkaIoC) GetProducer() (setup.IKafkaProducerSetup, setup.ISchemaRegistrySetup) {
	return ioc.producerSetup, ioc.schemaRegistrySetup
}

// GetSchemaRegistry retorna a interface para o schema registry
func (ioc *kafkaIoC) GetSchemaRegistry() setup.ISchemaRegistrySetup {
	return ioc.schemaRegistrySetup
}

// GetConsumerPriority retorna a prioridade do consumidor
func (ioc *kafkaIoC) GetConsumerPriority() enums.ConsumerOrderPriority {
	return ioc.consumerPriority
}

// GetProducerPriority retorna a prioridade do produtor
func (ioc *kafkaIoC) GetProducerPriority() enums.ProducerOrderPriority {
	return ioc.producerPriority
}
