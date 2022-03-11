package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/segmentio/kafka-go"
)

type config struct {
	KafkaUrl     string
	KafkaTopic   string
	KafkaGroupId string
}

func getenvOrPanic(name string) string {
	value, found := os.LookupEnv(name)
	if !found {
		panic(fmt.Sprintf("Missing %s", name))
	}

	return value
}

func SimulateLoad() {
	done := make(chan int)

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				select {
				case <-done:
					return
				default:
				}
			}
		}()
	}

	time.Sleep(time.Second * 2)
	close(done)
}

func ProcessMessage(ctx context.Context, msg *kafka.Message) {
	log.Printf("message: %s\n", msg.Value)
	SimulateLoad()
	log.Printf("completed simulating load\n")
}

func main() {
	log.Println("Running application")

	config := config{
		KafkaUrl:     getenvOrPanic("KAFKA_URL"),
		KafkaTopic:   getenvOrPanic("KAFKA_TOPIC"),
		KafkaGroupId: getenvOrPanic("KAFKA_CONSUMER_GROUP_ID"),
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{config.KafkaUrl},
		Topic:    config.KafkaTopic,
		GroupID:  config.KafkaGroupId,
		MinBytes: 100,
		MaxBytes: 500,
	})

	ctx := context.Background()

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			break
		}

		log.Printf("message at offset %d: %s = %s\n", msg.Offset, string(msg.Key), string(msg.Value))

		ProcessMessage(ctx, &msg)

		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Printf("failed to commit message: %s\n", err)
		} else {
			log.Printf("commited message\n")
		}
	}

	if err := reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
