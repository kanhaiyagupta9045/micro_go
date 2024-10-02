package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kanhaiyagupta9045/product_service/internal/handlers"
	"github.com/kanhaiyagupta9045/product_service/internal/kafka"
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
	repo := repository.NewProductRepository(os.Getenv("MONGODB_URL"))

	kafkaproducer, err := kafka.NewProducer([]string{"localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
		return
	}

	service := service.NewProductService(repo, kafkaproducer)

	handler := handlers.NewProductHandler(service)

	router := routes.ProductRoutes(handler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":" + os.Getenv("PORT")),
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
