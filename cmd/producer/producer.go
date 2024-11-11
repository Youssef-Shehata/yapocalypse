package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/IBM/sarama"
	"github.com/Youssef-Shehata/yapocalypse/pkg/logger"
	"github.com/Youssef-Shehata/yapocalypse/pkg/types"
	"github.com/pkg/errors"
)

const (
	ProducerPort       = ":8081"
	KafkaServerAddress = "localhost:9092"
	KafkaTopic         = "yaps"
    ERROR = logger.ERROR
    INFO = logger.INFO
)

type (
	User = types.User
	Yap  = types.Yap
)
type config struct {
	p      sarama.SyncProducer
	logger *logger.Logger
}

func Init() (config, *http.ServeMux) {

	logger, err := logger.NewLogger("./producer_server.log")
	if err != nil {
		log.Printf("Could not initialize logger: %v\n", err)
	}

	producer, err := setupProducer()
	if err != nil {
		log.Fatalf("failed to initialize producer: %v", err)
	}

	cfg := config{p: producer, logger: logger}
	mux := http.NewServeMux()
	return cfg, mux

}
func setupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()

	config.Producer.Return.Successes = true

	//single message per yap
	//config.Producer.Idempotent = true

	producer, err := sarama.NewSyncProducer([]string{KafkaServerAddress}, config)

	if err != nil {
		return nil, fmt.Errorf("failed to setup producer: %w", err)
	}

	return producer, nil
}

func (cfg *config) sendKafkaMessage(yap Yap) error {

	yapJson, err := json.Marshal(yap)

	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: KafkaTopic,
		Key:   sarama.StringEncoder(yap.ID.String()),
		Value: sarama.StringEncoder(yapJson),
	}

	p, o, err := cfg.p.SendMessage(msg)
	if err == nil {
		log.Printf("Yap is sent to kafka and stored in :\n \t Partition %v \n \t Offset %v \n", p, o)
	}
	return err
}





func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	res, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "failed to marshal json", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(res)
}
func (cfg config) sendMessageHandler(w http.ResponseWriter, r *http.Request) {

	var yap Yap
	json.NewDecoder(r.Body).Decode(&yap)
	if err := cfg.sendKafkaMessage(yap); err != nil {

        cfg.logger.Log(ERROR ,errors.Wrap(err ,"sending message"))
		respondWithJSON(w, http.StatusInternalServerError, err)

		return
	}
	respondWithJSON(w, http.StatusOK, yap)
}

func main() {
	cfg, mux := Init()
	defer cfg.p.Close()
	defer cfg.logger.Close()

	mux.HandleFunc("POST /yap", cfg.sendMessageHandler)

	server := http.Server{Handler: mux, WriteTimeout: 10 * time.Second, ReadTimeout: 10 * time.Second, Addr: "localhost" + ProducerPort}

	cfg.logger.Log(INFO, fmt.Errorf("Kafka PRODUCER 📨 started at http://localhost%s\n", ProducerPort))

	if err := server.ListenAndServe(); err != nil {
		cfg.logger.Log(ERROR, fmt.Errorf("starting server:  %v", err))
	}
}
