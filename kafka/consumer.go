package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

var usertopic = "user_topic"
var producttopic = "product_topic"
var shippingtopic = "shipping_topic"

type consumer struct{}

func (consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		switch message.Topic {
		case usertopic:

		case producttopic:

		case shippingtopic:

		default:
			fmt.Println("Unknown event type received")
		}

		sess.MarkMessage(message, "")

	}

	return nil
}
func ConsumeMessages(ctx context.Context, brokers []string, groupID string, topics []string) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %s", err)
	}

	defer consumerGroup.Close()

	for {
		if err := consumerGroup.Consume(ctx, topics, &consumer{}); err != nil {
			// After this the ConsumeClaim() function is called automatically
			log.Fatalf("Error from consumer: %s", err)
		}

		// Check if context is done
		if ctx.Err() != nil {
			log.Println("Context error: ", ctx.Err())
			return
		}
	}
}
