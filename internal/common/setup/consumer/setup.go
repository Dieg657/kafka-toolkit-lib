package setup

import (
	"os"

	"github.com/Dieg657/kafka-toolkit-lib/internal/common/config"
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/enums"
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/setup"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/spf13/viper"
)

// ==========================================================================
// Tipos e Propriedades
// ==========================================================================

// kafkaConsumerSetup implementação concreta privada
type kafkaConsumerSetup struct {
	consumerKafka *kafka.Consumer
}

// ==========================================================================
// Construtores
// ==========================================================================

// NewKafkaConsumerSetup cria uma nova instância da interface IKafkaConsumerSetup
func NewKafkaConsumerSetup(options config.IKafkaOptions) (setup.IKafkaConsumerSetup, error) {
	consumer := &kafkaConsumerSetup{}
	err := consumer.New(options)
	if err != nil {
		return nil, err
	}
	return consumer, nil
}

// ==========================================================================
// Configurações de Prioridade do Consumidor
// ==========================================================================

// Mapa de configurações do consumidor pelo tipo de prioridade escolhida pelo usuário
var consumerPriorityConfigs = map[enums.ConsumerOrderPriority]func(*kafka.ConfigMap){
	enums.CONSUMER_ORDER_PRIORITY_ORDER: func(configMap *kafka.ConfigMap) {
		// Configurações para garantia máxima de ordem e consistência
		configMap.SetKey("enable.auto.commit", "false")        // Commit manual para garantir processamento correto
		configMap.SetKey("auto.offset.reset", "earliest")      // Começa do início para não perder mensagens
		configMap.SetKey("max.partition.fetch.bytes", 1048576) // 1 MB por partição
		configMap.SetKey("fetch.max.bytes", 5242880)           // 5 MB total
		configMap.SetKey("max.poll.interval.ms", 300000)       // 5 min para processar cada lote
		configMap.SetKey("session.timeout.ms", 10000)          // 10 s timeout da sessão
		configMap.SetKey("heartbeat.interval.ms", 3000)        // 3 s para heartbeat
		configMap.SetKey("isolation.level", "read_committed")  // Lê apenas mensagens confirmadas
		configMap.SetKey("fetch.min.bytes", 1)                 // Não espera acumular dados
		configMap.SetKey("fetch.wait.max.ms", 500)             // Espera no máximo 500ms
		configMap.SetKey("retry.backoff.ms", 100)              // 100ms entre retentativas
		configMap.SetKey("fetch.error.backoff.ms", 500)        // 500ms de espera após erro de fetch
	},
	enums.CONSUMER_ORDER_PRIORITY_BALANCED: func(configMap *kafka.ConfigMap) {
		// Balanceamento entre consistência e performance
		configMap.SetKey("enable.auto.commit", "false")        // Commit manual para segurança
		configMap.SetKey("auto.offset.reset", "earliest")      // Começa do início para não perder mensagens
		configMap.SetKey("max.partition.fetch.bytes", 1048576) // 1 MB por partição
		configMap.SetKey("fetch.max.bytes", 10485760)          // 10 MB total
		configMap.SetKey("fetch.message.max.bytes", 262144)    // 256KB por mensagem
		configMap.SetKey("max.poll.interval.ms", 300000)       // 5 min para processar
		configMap.SetKey("session.timeout.ms", 30000)          // 30 s timeout
		configMap.SetKey("heartbeat.interval.ms", 10000)       // 10 s para heartbeat
		configMap.SetKey("isolation.level", "read_committed")  // Lê apenas mensagens confirmadas
		configMap.SetKey("fetch.min.bytes", 1024)              // Espera pelo menos 1KB
		configMap.SetKey("fetch.wait.max.ms", 1000)            // Espera até 1s
		configMap.SetKey("retry.backoff.ms", 200)              // 200ms entre retentativas
		configMap.SetKey("fetch.error.backoff.ms", 500)        // 500ms de espera após erro de fetch
	},
	enums.CONSUMER_ORDER_PRIORITY_HIGH_PERFORMANCE: func(configMap *kafka.ConfigMap) {
		// Configurações para alto throughput
		configMap.SetKey("enable.auto.commit", "true")          // Commit automático para velocidade
		configMap.SetKey("auto.commit.interval.ms", 5000)       // 5 s intervalo de commit
		configMap.SetKey("auto.offset.reset", "latest")         // Começa de mensagens mais recentes
		configMap.SetKey("max.partition.fetch.bytes", 10485760) // 10 MB por partição
		configMap.SetKey("fetch.max.bytes", 52428800)           // 50 MB total
		configMap.SetKey("fetch.message.max.bytes", 1048576)    // 1MB por mensagem
		configMap.SetKey("max.poll.interval.ms", 600000)        // 10 min para processar
		configMap.SetKey("session.timeout.ms", 60000)           // 60 s timeout
		configMap.SetKey("heartbeat.interval.ms", 20000)        // 20 s para heartbeat
		configMap.SetKey("isolation.level", "read_uncommitted") // Lê todas as mensagens
		configMap.SetKey("fetch.min.bytes", 65536)              // Espera 64KB para buscar
		configMap.SetKey("fetch.wait.max.ms", 100)              // Espera apenas 100ms
		configMap.SetKey("retry.backoff.ms", 50)                // 50ms entre retentativas
		configMap.SetKey("fetch.error.backoff.ms", 200)         // 200ms de espera após erro de fetch
	},
	enums.CONSUMER_ORDER_PRIORITY_RISKY: func(configMap *kafka.ConfigMap) {
		// Configurações para máximo throughput com menos garantias
		configMap.SetKey("enable.auto.commit", "true")          // Commit automático
		configMap.SetKey("auto.commit.interval.ms", 1000)       // 1 s intervalo de commit
		configMap.SetKey("auto.offset.reset", "latest")         // Começa de mensagens mais recentes
		configMap.SetKey("max.partition.fetch.bytes", 52428800) // 50 MB por partição
		configMap.SetKey("fetch.max.bytes", 104857600)          // 100 MB total
		configMap.SetKey("fetch.message.max.bytes", 4194304)    // 4MB por mensagem
		configMap.SetKey("max.poll.interval.ms", 900000)        // 15 min para processar
		configMap.SetKey("session.timeout.ms", 120000)          // 120 s timeout
		configMap.SetKey("heartbeat.interval.ms", 40000)        // 40 s para heartbeat
		configMap.SetKey("isolation.level", "read_uncommitted") // Lê todas as mensagens
		configMap.SetKey("fetch.min.bytes", 131072)             // Espera 128KB antes de buscar
		configMap.SetKey("fetch.wait.max.ms", 50)               // Espera no máximo 50ms
		configMap.SetKey("reconnect.backoff.ms", 10)            // Apenas 10ms para reconectar
		configMap.SetKey("retry.backoff.ms", 10)                // 10ms entre retentativas
		configMap.SetKey("fetch.error.backoff.ms", 100)         // 100ms de espera após erro de fetch
	},
}

// ==========================================================================
// Métodos Públicos
// ==========================================================================

// GetKafkaConsumer retorna a instância do consumidor Kafka
func (cs *kafkaConsumerSetup) GetKafkaConsumer() *kafka.Consumer {
	return cs.consumerKafka
}

// ==========================================================================
// Métodos Privados
// ==========================================================================

// New inicializa um novo consumidor Kafka com as configurações especificadas
func (cs *kafkaConsumerSetup) New(options config.IKafkaOptions) error {
	// Obter nome do host para identificação do cliente
	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	// Determinar a prioridade do consumidor
	consumerPriority := determineConsumerPriority(options)

	// Configuração base do consumidor
	configMap := &kafka.ConfigMap{
		"bootstrap.servers": options.GetBrokers(),
		"group.id":          options.GetGroupId(),
		"client.id":         hostname,
		"security.protocol": options.GetSecurityProtocol(),
		"sasl.mechanism":    options.GetSaslMechanism(),
		"sasl.username":     options.GetUserName(),
		"sasl.password":     options.GetPassword(),
		"auto.offset.reset": options.GetOffset(),
	}

	// Aplicar configurações específicas da prioridade escolhida
	applyConsumerPriorityConfig(consumerPriority, configMap)

	// Criar o consumidor Kafka
	consumer, err := kafka.NewConsumer(configMap)
	if err != nil {
		return err
	}

	cs.consumerKafka = consumer
	return nil
}

// determineConsumerPriority determina a prioridade do consumidor com base nas opções ou variáveis de ambiente
func determineConsumerPriority(options config.IKafkaOptions) enums.ConsumerOrderPriority {
	viper.AutomaticEnv()
	priorityStr := options.GetConsumerPriority()
	var consumerPriority enums.ConsumerOrderPriority

	if priorityStr != "" {
		consumerPriority = enums.ConsumerOrderPriority(priorityStr)
	} else {
		consumerPriority = enums.ConsumerOrderPriority(viper.GetString("KAFKA_CONSUMER_PRIORITY"))
	}

	// Validation
	if consumerPriority == "" || (consumerPriority != enums.CONSUMER_ORDER_PRIORITY_ORDER &&
		consumerPriority != enums.CONSUMER_ORDER_PRIORITY_BALANCED &&
		consumerPriority != enums.CONSUMER_ORDER_PRIORITY_HIGH_PERFORMANCE &&
		consumerPriority != enums.CONSUMER_ORDER_PRIORITY_RISKY) {
		consumerPriority = enums.CONSUMER_ORDER_PRIORITY_BALANCED
	}

	return consumerPriority
}

// applyConsumerPriorityConfig aplica as configurações de prioridade ao ConfigMap
func applyConsumerPriorityConfig(priority enums.ConsumerOrderPriority, configMap *kafka.ConfigMap) {
	if configFunc, exists := consumerPriorityConfigs[priority]; exists {
		configFunc(configMap)
	}
}
