package kafka

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

type Producer struct {
	producer sarama.SyncProducer
}

var (
	PRODUCT_TOPIC = "product_topic"
)

func NewProducer(brokers []string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Producer{producer: producer}, nil
}
func (p *Producer) ProduceMessage(topic string, message interface{}) error {
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonMessage),
	}

	_, _, err = p.producer.SendMessage(msg)
	return err
}

func (p *Producer) Close() error {
	return p.producer.Close()
}
