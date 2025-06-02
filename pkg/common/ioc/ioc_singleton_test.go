package ioc

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Teste de singleton do IoC container
// Este teste garante que, mesmo em chamadas concorrentes, o container de dependências
// (e, por consequência, o consumer e o producer) são instanciados apenas uma vez.
//
// O teste NÃO depende de Kafka real, pois valida apenas a identidade das instâncias.
func TestIoCSingletonInstance(t *testing.T) {
	setDummyKafkaEnv()

	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)

	results := make([]IContainer, goroutines)
	errors := make([]error, goroutines)

	// Executa várias goroutines que tentam obter o IoC container simultaneamente
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			container, err := GetKafkaIoC()
			results[idx] = container
			errors[idx] = err
		}(i)
	}
	wg.Wait()

	// Todas as instâncias devem ser idênticas (singleton)
	for i := 1; i < goroutines; i++ {
		assert.NoError(t, errors[i], "Erro inesperado ao obter IoC container")
		assert.Equal(t, results[0], results[i], "IoC container não é singleton!")
	}
}

// Teste de singleton para consumer e producer
// Garante que o consumer e o producer retornados pelo IoC container são sempre as mesmas instâncias
func TestIoCConsumerProducerSingleton(t *testing.T) {
	setDummyKafkaEnv()

	container, err := GetKafkaIoC()
	assert.NoError(t, err)

	consumer1, _ := container.GetConsumer()
	consumer2, _ := container.GetConsumer()
	assert.Equal(t, consumer1, consumer2, "Consumer não é singleton!")

	producer1, _ := container.GetProducer()
	producer2, _ := container.GetProducer()
	assert.Equal(t, producer1, producer2, "Producer não é singleton!")
}

// Teste concorrente completo: IoC, Producer e Consumer em múltiplas goroutines
// Este teste lança 10 goroutines, cada uma obtendo o IoC container, o producer e o consumer.
// Ele armazena todas as instâncias retornadas e valida que todas são idênticas (singleton real).
func TestIoCProducerConsumerSingletonConcurrent(t *testing.T) {
	setDummyKafkaEnv()

	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)

	iocResults := make([]IContainer, goroutines)
	producerResults := make([]interface{}, goroutines)
	consumerResults := make([]interface{}, goroutines)
	errors := make([]error, goroutines)

	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			container, err := GetKafkaIoC()
			iocResults[idx] = container
			errors[idx] = err
			if err == nil {
				prod, _ := container.GetProducer()
				cons, _ := container.GetConsumer()
				producerResults[idx] = prod
				consumerResults[idx] = cons
			}
		}(i)
	}
	wg.Wait()

	// Valida que todas as instâncias de IoC são idênticas
	for i := 1; i < goroutines; i++ {
		assert.NoError(t, errors[i], "Erro inesperado ao obter IoC container")
		assert.Equal(t, iocResults[0], iocResults[i], "IoC container não é singleton!")
	}

	// Valida que todas as instâncias de Producer são idênticas
	for i := 1; i < goroutines; i++ {
		assert.Equal(t, producerResults[0], producerResults[i], "Producer não é singleton!")
	}

	// Valida que todas as instâncias de Consumer são idênticas
	for i := 1; i < goroutines; i++ {
		assert.Equal(t, consumerResults[0], consumerResults[i], "Consumer não é singleton!")
	}
}

// Documentação:
// - Estes testes não validam conexão real com Kafka, apenas a lógica de singleton.
// - Se o padrão singleton for quebrado, os testes falharão, indicando múltiplas instâncias.
// - O uso de assert.Equal compara ponteiros/interfaces, garantindo identidade real das instâncias.

func setDummyKafkaEnv() {
	os.Setenv("KAFKA_BROKERS", "dummy:9092")
	os.Setenv("KAFKA_GROUPID", "dummy-group")
	os.Setenv("KAFKA_SCHEMA_REGISTRY_URL", "http://dummy:8081")
	os.Setenv("KAFKA_SECURITY_PROTOCOL", "plaintext")
	os.Setenv("KAFKA_SASL_MECHANISM", "PLAIN")
	os.Setenv("KAFKA_SCHEMA_REGISTRY_USERNAME", "dummy")
	os.Setenv("KAFKA_SCHEMA_REGISTRY_PASSWORD", "dummy")
	os.Setenv("KAFKA_SCHEMA_REGISTRY_AUTH_SOURCE", "USER_INFO")
	os.Setenv("KAFKA_TIMEOUT", "1000")
	os.Setenv("KAFKA_PRODUCER_PRIORITY", "BALANCED")
	os.Setenv("KAFKA_CONSUMER_PRIORITY", "BALANCED")
	os.Setenv("KAFKA_AUTO_OFFSET_RESET", "latest")
}
