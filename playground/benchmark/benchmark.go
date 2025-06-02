package benchmark

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Dieg657/kafka-toolkit-lib/pkg/common/constants"
	"github.com/Dieg657/kafka-toolkit-lib/pkg/common/enums"
	"github.com/Dieg657/kafka-toolkit-lib/pkg/common/ioc"
	"github.com/Dieg657/kafka-toolkit-lib/pkg/common/message"
	"github.com/Dieg657/kafka-toolkit-lib/pkg/publisher"
	"github.com/google/uuid"
)

// ==========================================================================
// Tipos para Mensagens JSON
// ==========================================================================

// JsonMessage define a estrutura da mensagem JSON simples para benchmark
type JsonMessage struct {
	HoraMensagem string `json:"hora_mensagem"`
	Producer     string `json:"producer"`
	MessageID    string `json:"message_id"`
	Counter      int64  `json:"counter"`
	Timestamp    int64  `json:"timestamp"` // Para futuras medições de latência end-to-end
}

// ==========================================================================
// Configurações e Resultados do Benchmark
// ==========================================================================

// BenchmarkConfig define as configurações para o benchmark de produção
type BenchmarkConfig struct {
	ProducerWorkers int           // Número de workers produtores
	Duration        time.Duration // Duração do teste
	ReportInterval  time.Duration // Intervalo para reportar estatísticas
	TopicName       string        // Nome do tópico Kafka
	PreWarmup       bool          // Indica se deve fazer warm-up
}

// ProducerResult contém os resultados do benchmark de produção
type ProducerResult struct {
	TotalMessages     int64         // Total de mensagens enviadas
	SuccessMessages   int64         // Mensagens enviadas com sucesso
	FailedMessages    int64         // Mensagens que falharam
	ElapsedTime       time.Duration // Tempo total de produção
	MessagesPerSecond float64       // Taxa de mensagens por segundo
	AverageLatencyMs  float64       // Latência média de produção em ms
	P95LatencyMs      float64       // 95º percentil de latência
	P99LatencyMs      float64       // 99º percentil de latência
}

// ProducerWorker gerencia métricas de produção
type ProducerWorker struct {
	successCount int64
	failureCount int64
	totalLatency int64
	latencies    []int64
	latencyMutex sync.Mutex
	ctx          context.Context
}

// NewDefaultConfig cria uma configuração padrão para benchmark de produção
func NewDefaultConfig() BenchmarkConfig {
	return BenchmarkConfig{
		ProducerWorkers: runtime.NumCPU(),
		Duration:        30 * time.Second,
		ReportInterval:  5 * time.Second,
		TopicName:       "benchmark-producer-topic",
		PreWarmup:       true,
	}
}

// InitializeBenchmark configura o contexto e dependências para o benchmark
func InitializeBenchmark() (context.Context, error) {
	// Define variável de ambiente necessária
	_ = os.Setenv("KAFKA_GROUPID", "benchmark-producer-group")

	ctx := context.Background()
	iocContainer, err := ioc.GetKafkaIoC()
	if err != nil {
		return nil, fmt.Errorf("falha ao inicializar IoC: %w", err)
	}

	// Adiciona o container IoC ao contexto
	ctx = context.WithValue(ctx, constants.IocKey, iocContainer)

	return ctx, nil
}

// ==========================================================================
// Métodos do Producer Worker
// ==========================================================================

// RunProducerBenchmark executa o benchmark de produção
func (pw *ProducerWorker) RunProducerBenchmark(ctx context.Context, config BenchmarkConfig) ProducerResult {
	fmt.Printf("🚀 Iniciando benchmark de PRODUÇÃO com %d workers por %s\n",
		config.ProducerWorkers, config.Duration)

	ctx, cancel := context.WithTimeout(ctx, config.Duration)
	defer cancel()

	var wg sync.WaitGroup
	startTime := time.Now()

	// Inicia os workers de produção
	for i := 0; i < config.ProducerWorkers; i++ {
		wg.Add(1)
		go pw.producerWorker(ctx, &wg, config.TopicName, i)
	}

	// Aguarda todos os workers terminarem
	wg.Wait()
	endTime := time.Now()
	elapsedTime := endTime.Sub(startTime)

	return pw.calculateProducerResult(elapsedTime)
}

// producerWorker é um worker individual de produção
func (pw *ProducerWorker) producerWorker(ctx context.Context, wg *sync.WaitGroup, topic string, workerID int) {
	defer wg.Done()

	var messageCounter int64

	for {
		select {
		case <-ctx.Done():
			return
		default:
			currentCount := atomic.AddInt64(&messageCounter, 1)

			// Cria mensagem com timestamp para futuras medições de latência end-to-end
			payload := JsonMessage{
				HoraMensagem: time.Now().Format("15:04:05.000 02/01/2006"),
				Producer:     fmt.Sprintf("producer_%d", workerID),
				MessageID:    uuid.New().String(),
				Counter:      currentCount,
				Timestamp:    time.Now().UnixNano(), // Timestamp para futuras medições
			}

			correlationID := uuid.New()
			msg, err := message.NewForData(correlationID, payload, nil)
			if err != nil {
				atomic.AddInt64(&pw.failureCount, 1)
				continue
			}

			startTime := time.Now()
			err = publisher.PublishMessage(ctx, topic, msg, enums.JsonSerialization)
			latency := time.Since(startTime)

			if err != nil {
				atomic.AddInt64(&pw.failureCount, 1)
			} else {
				atomic.AddInt64(&pw.successCount, 1)
				atomic.AddInt64(&pw.totalLatency, int64(latency))

				pw.latencyMutex.Lock()
				pw.latencies = append(pw.latencies, int64(latency))
				pw.latencyMutex.Unlock()
			}
		}
	}
}

// calculateProducerResult calcula os resultados do producer
func (pw *ProducerWorker) calculateProducerResult(elapsedTime time.Duration) ProducerResult {
	totalMsgs := atomic.LoadInt64(&pw.successCount) + atomic.LoadInt64(&pw.failureCount)
	successMsgs := atomic.LoadInt64(&pw.successCount)
	failedMsgs := atomic.LoadInt64(&pw.failureCount)

	var avgLatency, p95Latency, p99Latency float64
	if totalMsgs > 0 {
		avgLatency = float64(atomic.LoadInt64(&pw.totalLatency)) / float64(totalMsgs) / float64(time.Millisecond)
	}

	// Calcula percentis
	pw.latencyMutex.Lock()
	if len(pw.latencies) > 0 {
		// Para simplificar, usamos aproximação dos percentis
		p95Index := int(float64(len(pw.latencies)) * 0.95)
		p99Index := int(float64(len(pw.latencies)) * 0.99)
		if p95Index < len(pw.latencies) {
			p95Latency = float64(pw.latencies[p95Index]) / float64(time.Millisecond)
		}
		if p99Index < len(pw.latencies) {
			p99Latency = float64(pw.latencies[p99Index]) / float64(time.Millisecond)
		}
	}
	pw.latencyMutex.Unlock()

	return ProducerResult{
		TotalMessages:     totalMsgs,
		SuccessMessages:   successMsgs,
		FailedMessages:    failedMsgs,
		ElapsedTime:       elapsedTime,
		MessagesPerSecond: float64(totalMsgs) / elapsedTime.Seconds(),
		AverageLatencyMs:  avgLatency,
		P95LatencyMs:      p95Latency,
		P99LatencyMs:      p99Latency,
	}
}

// ==========================================================================
// Benchmark de Produção
// ==========================================================================

// RunProducerOnlyBenchmark executa benchmark focado apenas na produção
func RunProducerOnlyBenchmark(ctx context.Context, config BenchmarkConfig) ProducerResult {
	fmt.Println("=== BENCHMARK DE PRODUÇÃO ===")
	fmt.Printf("Configuração:\n")
	fmt.Printf("- Producers: %d workers\n", config.ProducerWorkers)
	fmt.Printf("- Duração: %s\n", config.Duration)
	fmt.Printf("- Tópico: %s\n\n", config.TopicName)

	// Cria worker de produção
	producerWorker := &ProducerWorker{ctx: ctx}

	// Executa benchmark de produção
	result := producerWorker.RunProducerBenchmark(ctx, config)

	printProducerBenchmarkResult(result)
	return result
}

// printProducerBenchmarkResult imprime os resultados do benchmark de produção
func printProducerBenchmarkResult(result ProducerResult) {
	fmt.Println("\n" + strings.Repeat("=", 70))
	fmt.Println("                 RESULTADOS DO BENCHMARK DE PRODUÇÃO")
	fmt.Println(strings.Repeat("=", 70))

	// Resultados de Produção
	fmt.Println("\n🚀 PRODUÇÃO:")
	fmt.Printf("   Total de mensagens: %d\n", result.TotalMessages)
	fmt.Printf("   Mensagens com sucesso: %d (%.2f%%)\n",
		result.SuccessMessages,
		100*float64(result.SuccessMessages)/float64(result.TotalMessages))
	fmt.Printf("   Mensagens falharam: %d (%.2f%%)\n",
		result.FailedMessages,
		100*float64(result.FailedMessages)/float64(result.TotalMessages))
	fmt.Printf("   Throughput: %.2f mensagens/segundo\n", result.MessagesPerSecond)
	fmt.Printf("   Latência média: %.2f ms\n", result.AverageLatencyMs)
	fmt.Printf("   Latência P95: %.2f ms\n", result.P95LatencyMs)
	fmt.Printf("   Latência P99: %.2f ms\n", result.P99LatencyMs)

	// Análise de Performance
	fmt.Println("\n📊 ANÁLISE DE PERFORMANCE:")

	successRate := float64(result.SuccessMessages) / float64(result.TotalMessages) * 100
	if successRate >= 99.5 {
		fmt.Printf("   Taxa de Sucesso: ✅ EXCELENTE (%.2f%%)\n", successRate)
	} else if successRate >= 95.0 {
		fmt.Printf("   Taxa de Sucesso: ⚠️  BOA (%.2f%%)\n", successRate)
	} else {
		fmt.Printf("   Taxa de Sucesso: ❌ BAIXA (%.2f%%) - Verifique configurações\n", successRate)
	}

	if result.MessagesPerSecond >= 1000 {
		fmt.Printf("   Throughput: ✅ ALTO (%.0f msg/s)\n", result.MessagesPerSecond)
	} else if result.MessagesPerSecond >= 500 {
		fmt.Printf("   Throughput: ⚠️  MÉDIO (%.0f msg/s)\n", result.MessagesPerSecond)
	} else {
		fmt.Printf("   Throughput: ❌ BAIXO (%.0f msg/s) - Considere otimizações\n", result.MessagesPerSecond)
	}

	if result.AverageLatencyMs <= 5.0 {
		fmt.Printf("   Latência: ✅ BAIXA (%.2f ms)\n", result.AverageLatencyMs)
	} else if result.AverageLatencyMs <= 20.0 {
		fmt.Printf("   Latência: ⚠️  MÉDIA (%.2f ms)\n", result.AverageLatencyMs)
	} else {
		fmt.Printf("   Latência: ❌ ALTA (%.2f ms) - Verifique rede/broker\n", result.AverageLatencyMs)
	}

	// Recomendações
	fmt.Println("\n💡 RECOMENDAÇÕES:")

	if result.MessagesPerSecond < 1000 {
		fmt.Println("   - Considere ajustar batch.size e linger.ms para melhor throughput")
		fmt.Println("   - Verifique se o número de workers está adequado para sua CPU")
	}

	if result.P99LatencyMs > 50.0 {
		fmt.Println("   - P99 de latência alto - verifique configurações do producer")
		fmt.Println("   - Considere ajustar acks, retries e timeout.ms")
	}

	if result.FailedMessages > 0 {
		fmt.Println("   - Mensagens falharam - verifique logs e conectividade com Kafka")
		fmt.Println("   - Considere ajustar configurações de retry e timeout")
	}

	fmt.Printf("\n   Tempo total: %s\n", result.ElapsedTime.Round(time.Second))
	fmt.Println(strings.Repeat("=", 70))
}
