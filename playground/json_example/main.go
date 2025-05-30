package main

import (
	"context"
	"encoding/json"
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
	"github.com/google/uuid"
)

type MyPayload struct {
	Field1 string `json:"field1"`
	Field2 int    `json:"field2"`
}

func main() {
	// Define variáveis de ambiente para o teste
	_ = os.Setenv("KAFKA_GROUPID", "json-example-group")

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

	topic := "json-topic"

	// Publica uma mensagem
	go func() {
		time.Sleep(1 * time.Second) // Aguarda um pouco para o consumer estar pronto
		payload := MyPayload{Field1: "json-example", Field2: 456}
		correlationId := uuid.New()
		msg, _ := message.NewForData(correlationId, payload, nil)
		err := publisher.PublishMessage(ctx, topic, msg, enums.JsonSerialization)
		if err != nil {
			fmt.Println("Erro ao publicar:", err)
			cancel()
			return
		}
		fmt.Println("Mensagem publicada com Json!")
	}()

	// Contador de mensagens recebidas para demonstração
	messageCount := 0

	handler := func(msg message.Message[MyPayload]) error {
		messageCount++
		b, _ := json.MarshalIndent(msg, "", "  ")
		fmt.Printf("Mensagem recebida (Json) #%d: %s\n", messageCount, b)

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
			err := consumer.ConsumeMessage(ctx, topic, enums.JsonDeserialization, enums.OnDeserializationIgnoreMessage, handler)
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
