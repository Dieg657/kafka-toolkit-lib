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

// kafkaProducerSetup implementação concreta privada
type kafkaProducerSetup struct {
	producerKafka *kafka.Producer
}

// NewKafkaProducerSetup cria uma nova instância da interface IKafkaProducerSetup
func NewKafkaProducerSetup(options config.IKafkaOptions) (setup.IKafkaProducerSetup, error) {
	producer := &kafkaProducerSetup{}
	err := producer.New(options)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

// Mapa de configurações do produtor pelo tipo de prioridade escolhida pelo usuário
var producerPriorityConfigs = map[enums.ProducerOrderPriority]func(*kafka.ConfigMap){
	enums.PRODUCER_ORDER_PRIORITY_ORDER: func(configMap *kafka.ConfigMap) {
		// Configurações para garantia máxima de ordem
		configMap.SetKey("acks", "all")                              // Aguardar confirmação de todos as replicas
		configMap.SetKey("retries", 5)                               // Número de retentativas
		configMap.SetKey("enable.idempotence", true)                 // Garante que mensagens não sejam duplicadas
		configMap.SetKey("max.in.flight.requests.per.connection", 1) // Apenas uma requisição por vez para garantir ordem
		configMap.SetKey("batch.num.messages", 1000)                 // Máximo de 1000 mensagens por lote
		configMap.SetKey("linger.ms", 0)                             // Não aguardar para formar lotes
		configMap.SetKey("queue.buffering.max.messages", 100000)     // Limite de mensagens em buffer
		configMap.SetKey("queue.buffering.max.kbytes", 1048576)      // 1GB em buffer
		configMap.SetKey("request.timeout.ms", 30000)                // 30 segundos timeout
		configMap.SetKey("message.timeout.ms", 120000)               // 2 minutos para entregar
		configMap.SetKey("retry.backoff.ms", 100)                    // 100ms entre retentativas
		configMap.SetKey("compression.type", "snappy")               // Compressão eficiente
		configMap.SetKey("queue.buffering.max.ms", 0)                // Não atrasa entregas para tentar otimizar lotes
	},
	enums.PRODUCER_ORDER_PRIORITY_BALANCED: func(configMap *kafka.ConfigMap) {
		// Balanceamento entre ordem e performance
		configMap.SetKey("acks", "all")                              // Aguardar confirmação de todos as replicas
		configMap.SetKey("retries", 3)                               // Número de retentativas
		configMap.SetKey("enable.idempotence", true)                 // Garante que mensagens não sejam duplicadas
		configMap.SetKey("max.in.flight.requests.per.connection", 5) // Permitir 5 requisições em paralelo
		configMap.SetKey("batch.num.messages", 10000)                // Máximo de 10000 mensagens por lote
		configMap.SetKey("linger.ms", 5)                             // Aguardar 5ms para formar lotes maiores
		configMap.SetKey("queue.buffering.max.messages", 200000)     // Limite de mensagens em buffer
		configMap.SetKey("queue.buffering.max.kbytes", 2097152)      // 2GB em buffer
		configMap.SetKey("request.timeout.ms", 30000)                // 30 segundos timeout
		configMap.SetKey("message.timeout.ms", 120000)               // 2 minutos para entregar
		configMap.SetKey("retry.backoff.ms", 100)                    // 100ms entre retentativas
		configMap.SetKey("compression.type", "lz4")                  // Compressão mais eficiente e mais rápida
		configMap.SetKey("queue.buffering.max.ms", 5)                // Atrasa até 5ms para otimizar lotes
	},
	enums.PRODUCER_ORDER_PRIORITY_HIGH_PERFORMANCE: func(configMap *kafka.ConfigMap) {
		// Configurações para máxima performance
		configMap.SetKey("acks", 1)                                   // Aguardar apenas confirmação do líder
		configMap.SetKey("retries", 1)                                // Apenas uma retentativa
		configMap.SetKey("enable.idempotence", false)                 // Sem idempotência para maior throughput
		configMap.SetKey("max.in.flight.requests.per.connection", 10) // Muitas requisições em paralelo
		configMap.SetKey("batch.num.messages", 50000)                 // Máximo de 50000 mensagens por lote
		configMap.SetKey("linger.ms", 10)                             // Aguardar 10ms para formar lotes grandes
		configMap.SetKey("queue.buffering.max.messages", 500000)      // Limite alto de mensagens em buffer
		configMap.SetKey("queue.buffering.max.kbytes", 4194304)       // 4GB em buffer
		configMap.SetKey("request.timeout.ms", 20000)                 // 20 segundos timeout
		configMap.SetKey("message.timeout.ms", 60000)                 // 1 minuto para entregar
		configMap.SetKey("retry.backoff.ms", 100)                     // 100ms entre retentativas
		configMap.SetKey("compression.type", "lz4")                   // Compressão mais eficiente e mais rápida
		configMap.SetKey("socket.keepalive.enable", true)             // Manter conexões ativas
		configMap.SetKey("queue.buffering.max.ms", 10)                // Atrasa até 10ms para otimizar lotes
	},
}

// ==========================================================================
// Construtores
// ==========================================================================

func (producerSetup *kafkaProducerSetup) New(options config.IKafkaOptions) error {
	viper.AutomaticEnv()

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	configMap := &kafka.ConfigMap{
		"bootstrap.servers":  options.GetBrokers(),
		"client.id":          hostname,
		"request.timeout.ms": options.GetRequestTimeout(),
		"security.protocol":  options.GetSecurityProtocol(),
		"sasl.mechanism":     options.GetSaslMechanisms(),
		"sasl.username":      options.GetUserName(),
		"sasl.password":      options.GetPassword(),
	}

	setProducerOrderPriority(enums.ProducerOrderPriority(options.GetProducerPriority()), configMap)

	producer, err := kafka.NewProducer(configMap)

	if err != nil {
		return err
	}

	producerSetup.producerKafka = producer
	return nil
}

// ==========================================================================
// Métodos Públicos
// ==========================================================================

func (producerSetup *kafkaProducerSetup) GetKafkaProducer() *kafka.Producer {
	return producerSetup.producerKafka
}

// ==========================================================================
// Métodos Privados
// ==========================================================================

// setProducerOrderPriority configura o produtor com base na prioridade escolhida
func setProducerOrderPriority(priority enums.ProducerOrderPriority, configMap *kafka.ConfigMap) {
	configFunc := producerPriorityConfigs[priority]
	configFunc(configMap)
}
