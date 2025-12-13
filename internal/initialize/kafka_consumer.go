package initialize

import (
	"app/global"
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

// DeliveryHandler defines the interface for processing Kafka messages
type DeliveryHandler interface {
	Handle(ctx context.Context, msg kafka.Message) error
}

// StartKafkaConsumer initializes and starts the Kafka consumer
func StartKafkaConsumer(handler DeliveryHandler) {
	global.Logger.Info("Starting Kafka Consumer...")

	// Configure Reader
	broker := fmt.Sprintf("%s:%d", global.Config.Kafka.Host, global.Config.Kafka.Port)

	// We subscribe to the configured topics
	topics := global.Config.Kafka.Topics
	global.Logger.Info(fmt.Sprintf("Subscribing to topics: %v", topics))

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{broker},
		GroupTopics:    topics,
		GroupID:        global.Config.Kafka.GroupID,
		MinBytes:       1,                // Fetch immediately
		MaxBytes:       10e6,             // 10MB
		MaxWait:        10 * time.Second, // Wait up to 10s if NO data. Realtime if data exists.
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
		ErrorLogger: kafka.LoggerFunc(func(msg string, a ...interface{}) {
			global.Logger.Error(fmt.Sprintf("Kafka Consumer Error: "+msg, a...))
		}),
	})

	// Run consumer in a separate goroutine
	go func() {
		defer func() {
			if err := r.Close(); err != nil {
				global.Logger.Error(fmt.Sprintf("Failed to close Kafka reader: %v", err))
			}
		}()

		for {
			ctx := context.Background()
			m, err := r.FetchMessage(ctx)
			if err != nil {
				global.Logger.Error(fmt.Sprintf("Failed to fetch message: %v", err))
				time.Sleep(time.Second) // Wait before retrying
				continue
			}

			// Delegate handling to the third_party layer
			if err := handler.Handle(ctx, m); err != nil {
				global.Logger.Error(fmt.Sprintf("Error handling message: %v", err))
			}

			if err := r.CommitMessages(ctx, m); err != nil {
				global.Logger.Error(fmt.Sprintf("Failed to commit message: %v", err))
			}
		}
	}()

	global.Logger.Info("Kafka Consumer started")
}
