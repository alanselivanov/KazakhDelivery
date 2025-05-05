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
	"inventory-service/internal/infrastructure/database"
	"inventory-service/internal/interfaces/routes"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	// Initialize MongoDB connection
	mongoDB, err := database.NewMongoDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
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

		cancel()
	}()

	grpcServer := grpc.NewServer()

	routes.RegisterGRPCServices(grpcServer, mongoDB)

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
