package consumer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/IBM/sarama"
	"github.com/kanhaiyagupta9045/pratilipi/order/internal/models"
	"github.com/kanhaiyagupta9045/pratilipi/order/internal/service"
)

const (
	USER_TOPIC    = "user_topic"
	PRODUCT_TOPIC = "product_topic"
)

var srv service.OrderService

func NewOrderConsumer(service *service.OrderService) {
	srv = *service
}

type Consumer struct{}

func (c *Consumer) Setup(session sarama.ConsumerGroupSession) error   { return nil }
func (c *Consumer) Cleanup(session sarama.ConsumerGroupSession) error { return nil }
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		switch msg.Topic {

		case USER_TOPIC:
			var event models.UserEvent
			err := json.Unmarshal(msg.Value, &event)
			if err != nil {
				log.Printf("Error unmarshalling message from USER_TOPIC: %v", err)
				continue
			}

			switch event.EventType {
			case "User Registered":
				if err := srv.CreateUser(event.Data); err != nil {
					log.Printf("Error While Inserting user in the order db: %s", err.Error())
				}
			case "User Profile Updated":
				if err := srv.UpdateUser(event.Data); err != nil {
					log.Printf("Error While Updating user in the order db: %s", err.Error())
				}
			default:
				log.Printf("Unknown user event type: %s", event.EventType)
			}
		case PRODUCT_TOPIC:

			var productEvent models.ProductCreatedEvent
			err := json.Unmarshal(msg.Value, &productEvent)
			if err == nil && productEvent.EventType == "Product Created" {
				if err := srv.CreateOrder(productEvent.Data); err != nil {
					log.Printf("Error While Creating product: %s", err.Error())
				}
				session.MarkMessage(msg, "")
				continue
			}

			var inventoryEvent models.InventoryUpdate
			err = json.Unmarshal(msg.Value, &inventoryEvent)
			if err != nil {
				log.Printf("Error unmarshalling message from PRODUCT_TOPIC: %v", err)
				continue
			}

			if inventoryEvent.EventType == "Inventory Updated" {
				if err := srv.UpdateInventory(inventoryEvent.ProductID, inventoryEvent.StockLevel); err != nil {
					log.Printf("Error While Updating Inventory: %s", err.Error())
				}
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
	consumerGroup, err := sarama.NewConsumerGroup(brokers, "order_service", config)
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
		err := consumerGroup.Consume(ctx, []string{USER_TOPIC, PRODUCT_TOPIC}, &Consumer{})
		if err != nil {
			log.Fatalf("Error from consumer: %v", err)
		}

		if ctx.Err() != nil {
			break
		}
	}
}
