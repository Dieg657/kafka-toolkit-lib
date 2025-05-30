package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Dieg657/kafka-toolkit-lib/internal/common/constants"
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/enums"
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/ioc"
	"github.com/Dieg657/kafka-toolkit-lib/internal/common/message"
	"github.com/Dieg657/kafka-toolkit-lib/internal/consumer"
	"github.com/Dieg657/kafka-toolkit-lib/internal/publisher"
	avro "github.com/Dieg657/kafka-toolkit-lib/playground/avro_example/avroschema"
	"github.com/google/uuid"
)

func main() {
	// Define variáveis de ambiente para o teste
	_ = os.Setenv("KAFKA_GROUPID", "avro-example-group")

	// Cria contexto com cancelamento por sinal
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Configura tratamento de sinais
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("\nRecebido sinal %v, terminando graciosamente...\n", sig)
		cancel()
	}()

	iocContainer, err := ioc.GetKafkaIoC()
	if err != nil {
		panic(err)
	}
	ctx = context.WithValue(ctx, constants.IocKey, iocContainer)

	topic := "avro-topic"

	// Publica uma mensagem
	go func() {
		time.Sleep(1 * time.Second) // Aguarda um pouco para o consumer estar pronto
		payload := avro.AvroPayload{Field1: "avro", Field2: 789}
		correlationId := uuid.New()
		msg, _ := message.NewForData(correlationId, payload, nil)
		err = publisher.PublishMessage(ctx, topic, msg, enums.AvroSerialization)
		if err != nil {
			fmt.Println("Erro ao publicar:", err)
			cancel()
			return
		}
		fmt.Println("Mensagem publicada com Avro!")
	}()

	// Contador de mensagens recebidas para demonstração
	messageCount := 0

	// Consome a mensagem
	handler := func(msg message.Message[avro.AvroPayload]) error {
		messageCount++
		fmt.Printf("Mensagem recebida (Avro) #%d: %+v\n", messageCount, msg)

		// Para demonstração, termina após receber uma mensagem
		// Remova esta linha se quiser consumir continuamente
		if messageCount >= 1 {
			fmt.Println("Exemplo concluído. Use Ctrl+C para terminar se executando continuamente.")
			cancel()
		}
		return nil
	}

	fmt.Println("Iniciando consumo... (Pressione Ctrl+C para terminar)")

	// Loop principal com verificação de contexto
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Contexto cancelado, terminando consumer...")
			return
		default:
			err := consumer.ConsumeMessage(ctx, topic, enums.AvroDeserialization, enums.OnDeserializationIgnoreMessage, handler)
			if err != nil {
				fmt.Println("Erro ao consumir:", err)
				// Se o contexto foi cancelado, não tenta novamente
				if ctx.Err() != nil {
					return
				}
				// Pausa breve em caso de erro, mas verifica contexto
				select {
				case <-ctx.Done():
					return
				case <-time.After(2 * time.Second):
					continue
				}
			}
		}
	}
}
