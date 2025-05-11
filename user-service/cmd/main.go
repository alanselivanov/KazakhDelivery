package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"user-service/internal/config"
	"user-service/internal/infrastructure/database"
	"user-service/internal/interfaces/routes"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	mongoDB, err := database.NewMongoDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	grpcServer := grpc.NewServer()

	services := routes.RegisterGRPCServices(grpcServer, mongoDB, cfg)

	lis, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("User service is running on port %s", cfg.Server.Port)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutdown signal received, closing connections...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	grpcServer.GracefulStop()
	if err := mongoDB.Close(ctx); err != nil {
		log.Fatalf("Error while closing MongoDB connection: %v", err)
	}

	if services.RedisCache != nil {
		if err := services.RedisCache.Close(); err != nil {
			log.Fatalf("Error while closing Redis connection: %v", err)
		}
	}

	log.Println("Server stopped successfully")
}
