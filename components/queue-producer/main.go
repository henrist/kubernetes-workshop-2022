package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

const message = "hello world"

func getenvOrPanic(name string) string {
	value, found := os.LookupEnv(name)
	if !found {
		panic(fmt.Sprintf("Missing %s", name))
	}

	return value
}

func sendManyMessages(ctx context.Context, writer *kafka.Writer) {
	messages := []kafka.Message{}

	for i := 0; i < 1000; i++ {
		messages = append(messages, kafka.Message{
			Value: []byte(message),
		})
	}

	err := writer.WriteMessages(ctx, messages...)
	if err != nil {
		log.Fatal("failed to write messages:", err)
		return
	}

	log.Printf("sent messages\n")
}

func keepSendingMessages(ctx context.Context, writer *kafka.Writer) {
	for {
		err := writer.WriteMessages(ctx, kafka.Message{
			Value: []byte(message),
		})

		if err != nil {
			log.Fatal("failed to write messages:", err)
			break
		}

		log.Printf("sent message\n")

		time.Sleep(2 * time.Second)
	}
}

func main() {
	log.Println("Running application")

	kafkaUrl := getenvOrPanic("KAFKA_URL")
	kafkaTopic := getenvOrPanic("KAFKA_TOPIC")

	burstThenQuit := os.Getenv("BURST_THEN_QUIT") == "true"

	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaUrl),
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}

	ctx := context.Background()

	if burstThenQuit {
		log.Printf("Burst then quit\n")
		sendManyMessages(ctx, writer)
	} else {
		keepSendingMessages(ctx, writer)
	}

	if err := writer.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}
