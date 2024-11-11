package main

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"github.com/Youssef-Shehata/yapocalypse/pkg/types"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	ConsumerGroup      = "yaps-group"
	ConsumerTopic      = "yaps"
	KafkaServerAddress = "localhost:9092"
)

type (
	Yap      = types.Yap
	Consumer struct{}

)
func (*Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (consumer *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		yapID := string(msg.Key)
		var yap Yap
		err := json.Unmarshal(msg.Value, &yap)
		if err != nil {
			log.Printf("failed to unmarshal yap: %v", err)
			continue
		}

		log.Printf("consuming : id(%v) \n \t YAP : %v \n", yapID, yap)

		//TODO: go addToFeeds

		sess.MarkMessage(msg, "")
	}
	return nil
}

func setupConsumerGroup(ctx context.Context) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	consumerGroup, err := sarama.NewConsumerGroup([]string{KafkaServerAddress}, ConsumerGroup, config)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{ConsumerTopic}, &Consumer{}); err != nil {
				log.Printf("error from consumer: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	return consumerGroup, nil
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumerGroup, err := setupConsumerGroup(ctx)
	if err != nil {
		log.Fatalf("failed to set up consumer group: %v", err)
	}
	defer consumerGroup.Close()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	log.Println("Kafka CONSUMER 📨 started")
	<-ctx.Done()
	log.Println("Kafka CONSUMER 📨 stopped")
}

