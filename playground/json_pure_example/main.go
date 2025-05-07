package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
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
	_ = os.Setenv("KAFKA_GROUPID", "json-pure-example-group")

	ctx := context.Background()
	iocContainer, err := ioc.NewKafkaIoC(ctx)
	if err != nil {
		panic(err)
	}
	ctx = context.WithValue(ctx, constants.IocKey, iocContainer)

	topic := "json-pure-topic"

	// Publica uma mensagem
	go func() {
		payload := MyPayload{Field1: "pure-json", Field2: 456}
		correlationId := uuid.New()
		msg, _ := message.NewForData(correlationId, payload, nil)
		err := publisher.PublishMessage(ctx, topic, msg, enums.JsonSerialization)
		if err != nil {
			fmt.Println("Erro ao publicar:", err)
			os.Exit(1)
		}
		fmt.Println("Mensagem publicada com PureJson!")
	}()

	handler := func(msg message.Message[MyPayload]) error {
		b, _ := json.MarshalIndent(msg, "", "  ")
		fmt.Printf("Mensagem recebida (PureJson): %s\n", b)
		os.Exit(0)
		return nil
	}
	for {
		err := consumer.ConsumeMessage(ctx, topic, enums.JsonDeserialization, handler)
		if err != nil {
			fmt.Println("Erro ao consumir:", err)
		}
		time.Sleep(2 * time.Second)
	}
}
