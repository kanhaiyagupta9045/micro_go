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
	"github.com/kanhaiyagupta9045/product_service/internal/consumer"
	"github.com/kanhaiyagupta9045/product_service/internal/handlers"
	"github.com/kanhaiyagupta9045/product_service/internal/producer"
	"github.com/kanhaiyagupta9045/product_service/internal/repository"
	"github.com/kanhaiyagupta9045/product_service/internal/routes"
	"github.com/kanhaiyagupta9045/product_service/internal/service"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		panic(err.Error())
	}
}
func main() {
	reps, err := repository.NewProductRepository(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err.Error())
	}

	kafkaproducer, err := producer.NewProducer([]string{"localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
		return
	}

	service := service.NewProductService(reps, kafkaproducer)
	consumer.NewOrderConsumer(service)
	handler := handlers.NewProductHandler(service)
	router := routes.ProductRoutes(handler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":" + os.Getenv("PORT")),
		Handler: router,
	}

	go consumer.StartConsumer()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()
	log.Println("Server started on port", os.Getenv("ADDR"))
	<-stop

	log.Println("Shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")

}
