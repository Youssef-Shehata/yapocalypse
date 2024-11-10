package producer 

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/Youssef-Shehata/yapocalypse/cmd/types"
	"github.com/IBM/sarama"
)

const (
	ProducerPort       = ":8081"
	KafkaServerAddress = "localhost:9092"
	KafkaTopic         = "yaps"
)

type (
	User = types.User
	Yap  = types.Yap
)

type Producer struct{
    Sync_producer sarama.SyncProducer

}

func (producer *Producer)SendKafkaMessage( yap Yap) error {

	yapJson, err := json.Marshal(yap)

	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: KafkaTopic,
		Key:   sarama.StringEncoder(yap.ID.String()),
		Value: sarama.StringEncoder(yapJson),
	}

	p, o, err := producer.Sync_producer.SendMessage(msg)
	log.Printf("Yap is sent to kafka and stored in :\n \t Partition %v \n \t Offset %v \n", p, o)
	return err
}

func SetupProducer() (Producer, error) {
	config := sarama.NewConfig()

	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{KafkaServerAddress}, config)

	if err != nil {
		return Producer{}, fmt.Errorf("failed to setup producer: %w", err)
	}

	return Producer{Sync_producer: producer}, nil
}

