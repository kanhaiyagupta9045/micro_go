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
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/handler"
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/kafka"
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/repository"
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/routes"
	"github.com/kanhaiyagupta9045/pratilipi/user/internal/service"
)

func init() {
	if err := godotenv.Load("/app/.env"); err != nil {
		panic(err)
	}
}

func main() {
	repo, err := repository.NewUserRepository(os.Getenv("DOCKER_POSTRES_URL"))

	if err != nil {
		panic(err.Error())
	}

	kafkaproducer, err := kafka.NewProducer([]string{"kafka:29092"})
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

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()
	log.Printf("Server started on port %s", os.Getenv("ADDR"))
	<-stop

	log.Println("Shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")

}
