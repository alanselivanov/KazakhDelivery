package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"inventory-service/internal/config"
	"inventory-service/internal/infrastructure/cache"
	"inventory-service/internal/infrastructure/database"
	"inventory-service/internal/infrastructure/messaging"
	"inventory-service/internal/infrastructure/product"
	"inventory-service/internal/interfaces/routes"

	gogrpc "google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	mongoDB, err := database.NewMongoDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	redisClient, err := cache.NewRedisClient(cfg)
	if err != nil {
		log.Printf("Failed to connect to Redis, proceeding without caching: %v", err)
		redisClient = nil
	}

	natsConsumer, err := messaging.NewNATSConsumer(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	consumer := natsConsumer
	log.Println("Using NATS consumer")

	productClient, err := product.NewProductServiceClient(cfg, mongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to Product Service: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

		if redisClient != nil {
			if err := redisClient.Close(); err != nil {
				log.Printf("Error during Redis disconnect: %v", err)
			}
		}

		consumer.Close()
		productClient.Close()

		cancel()
	}()

	grpcServer := gogrpc.NewServer()

	routes.RegisterGRPCServices(grpcServer, mongoDB, consumer, productClient, redisClient)

	lis, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Inventory Service running on port %s", cfg.Server.Port)

	go func() {
		<-ctx.Done()
		log.Println("Stopping gRPC server...")
		grpcServer.GracefulStop()
	}()

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
