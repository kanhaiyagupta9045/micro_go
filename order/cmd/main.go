package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/kanhaiyagupta9045/pratilipi/order/internal/consumer"
	"github.com/kanhaiyagupta9045/pratilipi/order/internal/handlers"
	"github.com/kanhaiyagupta9045/pratilipi/order/internal/producer"
	"github.com/kanhaiyagupta9045/pratilipi/order/internal/repository"
	"github.com/kanhaiyagupta9045/pratilipi/order/internal/routes"
	"github.com/kanhaiyagupta9045/pratilipi/order/internal/service"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}
}
func main() {
	repo, err := repository.NewOrderRepositry(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Panic(err.Error())
	}
	producer, err := producer.NewProducer([]string{"localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
		return
	}

	srv := service.NewOrderService(repo, producer)
	consumer.NewOrderConsumer(srv)

	handler := handlers.NewOrderHandler(srv)

	router := routes.OrderRoutes(handler)

	server := http.Server{
		Addr:    fmt.Sprintf(":" + os.Getenv("PORT")),
		Handler: router,
	}

	go consumer.StartConsumer()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()
	log.Println("Server started on port", os.Getenv("ADDR"))
	<-stop

	log.Println("Shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")

}
