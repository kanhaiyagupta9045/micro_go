package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/handler"
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/kafka"
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/repository"
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/routes"
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/service"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}
}

func main() {
	repo, err := repository.NewUserRepository(os.Getenv("DATABASE_URL"))

	if err != nil {
		panic(err.Error())
	}

	kafkaproducer, err := kafka.NewProducer([]string{"localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
		return
	}
	svc := service.NewUserService(repo, kafkaproducer)
	handler := handler.NewUserHandler(svc)

	router := routes.Router(handler)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":" + os.Getenv("ADDR")),
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}

}
