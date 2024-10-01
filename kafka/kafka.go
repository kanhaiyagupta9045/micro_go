package kafka

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func StartKafkaConsumer() {
	brokers := []string{"kafka:29092"}
	if err := InitProducer(brokers); err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go ConsumeMessages(ctx, brokers, "user_event_group", []string{usertopic, producttopic, shippingtopic})

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down gracefully...")
	os.Exit(0)

}
