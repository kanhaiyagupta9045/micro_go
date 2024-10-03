package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"github.com/kanhaiyagupta9045/product_service/internal/models"
	"github.com/kanhaiyagupta9045/product_service/internal/service"
)

const (
	USER_TOPIC    = "user_topic"
	PRODUCT_TOPIC = "product_topic"
	Order_Topic   = "user_topic"
)

var srv *service.ProductService

func NewOrderConsumer(service *service.ProductService) {
	srv = service
}

type Consumer struct{}

func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error   { return nil }
func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error { return nil }
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		switch msg.Topic {
		case Order_Topic:
			var event models.OrderEvent
			err := json.Unmarshal(msg.Value, &event)
			if err != nil {
				log.Printf("Error unmarshalling message from USER_TOPIC: %v", err)
				continue
			}
			if err := srv.UpdateInventoryEvent(event.ProductID); err != nil {
				log.Printf("Error While Updating inventory in the order db: %s", err.Error())
			}

		default:
			log.Printf("Unknown topic: %s", msg.Topic)
		}
		session.MarkMessage(msg, "")
	}
	return nil
}

func StartConsumer() {
	brokers := []string{"localhost:9092"}
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumerGroup, err := sarama.NewConsumerGroup(brokers, "product_service", config)
	if err != nil {
		log.Fatalf("Error creating consumer group client: %v", err)
	}
	defer func() {
		if err := consumerGroup.Close(); err != nil {
			log.Fatalf("Error closing consumer group: %v", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		<-signals
		cancel()
	}()
	for {
		err := consumerGroup.Consume(ctx, []string{Order_Topic}, &Consumer{})
		if err != nil {
			log.Fatalf("Error from consumer: %v", err)
		}

		if ctx.Err() != nil {
			break
		}
	}
}
