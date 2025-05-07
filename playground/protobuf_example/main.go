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
	pb "github.com/Dieg657/kafka-toolkit-lib/playground/protobuf_example/pb"
	"github.com/google/uuid"
)

func main() {
	// Define variáveis de ambiente para o teste
	_ = os.Setenv("KAFKA_GROUPID", "protobuf-example-group")

	ctx := context.Background()
	iocContainer, err := ioc.NewKafkaIoC(ctx)
	if err != nil {
		panic(err)
	}
	ctx = context.WithValue(ctx, constants.IocKey, iocContainer)

	topic := "protobuf-topic"

	// Publica uma mensagem
	go func() {
		payload := &pb.ProtoPayload{
			Field1: "protobuf",
			Field2: 321,
		}
		correlationId := uuid.New()
		msg, _ := message.NewForData(correlationId, payload, nil)
		err := publisher.PublishMessage(ctx, topic, msg, enums.ProtobufSerialization)
		if err != nil {
			fmt.Println("Erro ao publicar:", err)
			os.Exit(1)
		}
		fmt.Println("Mensagem publicada com Protobuf!")
	}()

	// Consome a mensagem
	handler := func(msg message.Message[*pb.ProtoPayload]) error {
		fmt.Printf("Mensagem recebida (Protobuf): %+v\n", msg.Data)
		os.Exit(0)
		return nil
	}
	for {
		err := consumer.ConsumeMessage(ctx, topic, enums.ProtobufDeserialization, handler)
		if err != nil {
			fmt.Println("Erro ao consumir:", err)
			os.Exit(1)
		}
		time.Sleep(2 * time.Second)
	}
}
