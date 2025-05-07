package main

import (
	"context"
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
	_ = os.Setenv("KAFKA_GROUPID", "json-schema-example-group")

	ctx := context.Background()
	iocContainer, err := ioc.NewKafkaIoC(ctx)
	if err != nil {
		panic(err)
	}
	ctx = context.WithValue(ctx, constants.IocKey, iocContainer)

	topic := "json-schema-topic"

	// Publica uma mensagem
	go func() {
		payload := MyPayload{Field1: "json-schema", Field2: 123}
		correlationId := uuid.New()
		msg, _ := message.NewForData(correlationId, payload, nil)
		err := publisher.PublishMessage(ctx, topic, msg, enums.JsonSchemaSerialization)
		if err != nil {
			fmt.Println("Erro ao publicar:", err)
			os.Exit(1)
		}
		fmt.Println("Mensagem publicada com JSON Schema!")
	}()

	// Consome a mensagem
	handler := func(msg message.Message[MyPayload]) error {
		fmt.Printf("Mensagem recebida (JSON Schema): %+v\n", msg)
		os.Exit(0)
		return nil
	}
	for {
		err := consumer.ConsumeMessage(ctx, topic, enums.JsonSchemaDeserialization, handler)
		if err != nil {
			fmt.Println("Erro ao consumir:", err)
		}
		time.Sleep(2 * time.Second)
	}
}
