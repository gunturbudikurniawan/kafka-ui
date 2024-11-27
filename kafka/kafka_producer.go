package kafka

import (
	"context"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(brokerAddresses []string, topic string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokerAddresses...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll, // Wait for all replicas
	}

	return &KafkaProducer{writer: writer}
}

func (kp *KafkaProducer) SendMessage(message string) error {
	err := kp.writer.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(message),
		Time:  time.Now(),
	})
	if err != nil {
		log.Printf("Failed to write message to Kafka: %v", err)
		return err
	}
	log.Printf("Message sent successfully to topic %s", kp.writer.Topic)
	return nil
}

func (kp *KafkaProducer) Close() error {
	return kp.writer.Close()
}
