package consumer

import (
	"sync"

	"github.com/IBM/sarama"
	"github.com/Youssef-Shehata/yapocalypse/cmd/types"
)


const (
    ConsumerGroup      = "yaps-group"
    ConsumerTopic      = "yaps"
    ConsumerPort       = ":8083"
    KafkaServerAddress = "localhost:9092"
)

type Yap = types.Yap


type UserYaps map[string][]Yap

type YapStore struct {
    data UserYaps
    mu   sync.RWMutex
}

func (ns *YapStore) Add(userID string, yap Yap) {
    ns.mu.Lock()
    defer ns.mu.Unlock()

    ns.data[userID] = append(ns.data[userID], yap)
}

func (ns *YapStore) Get(userID string) []Yap {
    ns.mu.RLock()
    defer ns.mu.RUnlock()

    return ns.data[userID]
}

type Consumer struct {
    store *YapStore
}

func (*Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }
