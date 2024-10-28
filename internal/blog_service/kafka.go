package blogservice

import (
	"log/slog"

	"github.com/IBM/sarama"
)

type KafkaService struct {
	log *slog.Logger
}


func NewKafka(log *slog.Logger) KafkaService {
	return KafkaService{
		log: log,
	}
}

func (k *KafkaService)ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(brokers, config)
}


func (k *KafkaService)ConnectConsumer(brokers []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	return sarama.NewConsumer(brokers, config)
}


func (k *KafkaService)PushPostToQueue(message []byte) error {
	brokers := []string {"localhost:9092"}

	// Connect to producer
	producer, err := k.ConnectProducer(brokers)
	if err != nil {
		return err
	}

	defer producer.Close()

	// New message
	msg := &sarama.ProducerMessage{
		Topic: "emails",
		Key: sarama.StringEncoder("Test"),
		Value: sarama.StringEncoder(message),
	}

	// Send message
	slog.Info("Sending message to queue")
	partition, _, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}
	slog.Info("Message sended", partition)
	// k.log.Info("Post stored in ", 
	// 	slog.String("topic", topic), 
	// 	slog.Int("partition", partition), 
	// 	slog.Int("offset", offset),
	// )
	return nil
}