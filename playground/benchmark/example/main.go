package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dieg657/kafka-toolkit-lib/internal/common/constants"
	"github.com/Dieg657/kafka-toolkit-lib/playground/benchmark"
)

func main() {
	fmt.Println("=== Kafka Toolkit Lib - Benchmark de Produção ===")
	fmt.Println("Iniciando benchmark usando mensagens JSON...")

	// Cria contexto com cancelamento por sinal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configura tratamento de sinais
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("\nRecebido sinal %v, terminando benchmark graciosamente...\n", sig)
		cancel()
	}()

	// Inicializa o contexto e dependências seguindo o padrão da biblioteca
	benchCtx, err := benchmark.InitializeBenchmark()
	if err != nil {
		log.Fatalf("Erro ao inicializar benchmark: %v", err)
	}

	// Combina o contexto de benchmark com o contexto de cancelamento
	ctx = context.WithValue(ctx, constants.IocKey, benchCtx.Value(constants.IocKey))

	// Configura o benchmark de produção
	config := benchmark.NewDefaultConfig()
	config.ProducerWorkers = 8                    // 8 workers produtores
	config.Duration = 30 * time.Second            // 30 segundos de teste
	config.ReportInterval = 5 * time.Second       // Relatório a cada 5s
	config.TopicName = "benchmark-producer-topic" // Tópico específico

	fmt.Printf("Configuração do benchmark:\n")
	fmt.Printf("- Producers: %d workers\n", config.ProducerWorkers)
	fmt.Printf("- Duração: %s\n", config.Duration)
	fmt.Printf("- Tópico: %s\n", config.TopicName)
	fmt.Printf("- Pressione Ctrl+C para interromper a qualquer momento\n")

	// Executa o benchmark de produção
	result := benchmark.RunProducerOnlyBenchmark(ctx, config)

	// Verifica se foi cancelado por sinal
	if ctx.Err() == context.Canceled {
		fmt.Println("\n⚠️  Benchmark interrompido pelo usuário")
		fmt.Println("Resultados parciais até o momento da interrupção:")
	}

	// Análise adicional dos resultados
	fmt.Println("\n🔍 ANÁLISE DETALHADA:")

	// Analisa a taxa de sucesso
	successRate := float64(result.SuccessMessages) / float64(result.TotalMessages) * 100
	if successRate >= 99.5 {
		fmt.Println("✅ Taxa de sucesso excelente - sistema funcionando perfeitamente")
	} else if successRate >= 95.0 {
		fmt.Printf("⚠️  Taxa de sucesso boa (%.2f%%) - algumas falhas ocasionais\n", successRate)
	} else {
		fmt.Printf("❌ Taxa de sucesso baixa (%.2f%%) - verifique configurações e conectividade\n", successRate)
	}

	// Analisa o throughput
	if result.MessagesPerSecond >= 1000 {
		fmt.Printf("🚀 Throughput alto (%.0f msg/s) - excelente performance\n", result.MessagesPerSecond)
	} else if result.MessagesPerSecond >= 500 {
		fmt.Printf("⚠️  Throughput médio (%.0f msg/s) - há espaço para otimização\n", result.MessagesPerSecond)
	} else {
		fmt.Printf("❌ Throughput baixo (%.0f msg/s) - necessita otimização urgente\n", result.MessagesPerSecond)
	}

	// Analisa latência
	if result.AverageLatencyMs <= 5.0 {
		fmt.Printf("✅ Latência baixa (%.2f ms) - resposta rápida\n", result.AverageLatencyMs)
	} else if result.AverageLatencyMs <= 20.0 {
		fmt.Printf("⚠️  Latência média (%.2f ms) - aceitável para maioria dos casos\n", result.AverageLatencyMs)
	} else {
		fmt.Printf("❌ Latência alta (%.2f ms) - pode impactar performance\n", result.AverageLatencyMs)
	}

	// Recomendações específicas baseadas nos resultados
	fmt.Println("\n💡 RECOMENDAÇÕES ESPECÍFICAS:")

	if result.MessagesPerSecond < 1000 {
		fmt.Println("📈 Para melhorar throughput:")
		fmt.Println("   - Ajuste batch.size (ex: 16384 ou 32768)")
		fmt.Println("   - Configure linger.ms (ex: 5-10ms)")
		fmt.Println("   - Considere aumentar buffer.memory")
		fmt.Println("   - Verifique se o número de workers está adequado")
	}

	if result.P99LatencyMs > 50.0 {
		fmt.Println("⏱️  Para reduzir latência P99:")
		fmt.Printf("   - P99 atual: %.2f ms (alto)\n", result.P99LatencyMs)
		fmt.Println("   - Ajuste request.timeout.ms")
		fmt.Println("   - Configure delivery.timeout.ms")
		fmt.Println("   - Verifique configurações de acks")
	}

	if result.FailedMessages > 0 {
		fmt.Println("🔧 Para reduzir falhas:")
		fmt.Printf("   - %d mensagens falharam\n", result.FailedMessages)
		fmt.Println("   - Verifique logs detalhados do Kafka")
		fmt.Println("   - Ajuste configurações de retry")
		fmt.Println("   - Verifique conectividade de rede")
	}

	// Estimativa de capacidade
	estimatedHourlyMessages := int64(result.MessagesPerSecond * 3600)
	estimatedDailyMessages := estimatedHourlyMessages * 24

	fmt.Println("\n📊 ESTIMATIVA DE CAPACIDADE:")
	fmt.Printf("   - Por hora: ~%d mensagens\n", estimatedHourlyMessages)
	fmt.Printf("   - Por dia: ~%d mensagens\n", estimatedDailyMessages)

	if ctx.Err() == context.Canceled {
		fmt.Printf("\n⚠️  Benchmark interrompido pelo usuário após %s\n",
			result.ElapsedTime.Round(time.Second))
	} else {
		fmt.Printf("\n🎯 Benchmark de produção concluído com sucesso em %s!\n",
			result.ElapsedTime.Round(time.Second))
	}
}
