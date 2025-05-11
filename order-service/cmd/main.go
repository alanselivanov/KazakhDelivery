package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"order-service/internal/config"
	"order-service/internal/infrastructure/database"
	"order-service/internal/infrastructure/messaging"
	"order-service/internal/interfaces/routes"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	mongoDB, err := database.NewMongoDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	natsPublisher, err := messaging.NewNATSPublisher(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	publisher := natsPublisher
	log.Println("Using NATS publisher")

	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	grpcServer := grpc.NewServer()

	services := routes.RegisterGRPCServices(grpcServer, mongoDB, publisher)

	lis, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Order Service running on port %s", cfg.Server.Port)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down gracefully...")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()

		if err := mongoDB.Close(shutdownCtx); err != nil {
			log.Printf("Error during MongoDB disconnect: %v", err)
		}

		publisher.Close()

		if services.RedisCache != nil {
			if err := services.RedisCache.Close(); err != nil {
				log.Printf("Error during Redis disconnect: %v", err)
			}
		}

		grpcServer.GracefulStop()
		cancel()
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
