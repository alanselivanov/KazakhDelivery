package main

import (
	"log"
	"net"

	"user-service/internal/config"
	"user-service/internal/infrastructure/database"
	"user-service/internal/interfaces/routes"

	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()

	db := database.NewInMemoryDB()

	grpcServer := grpc.NewServer()

	routes.RegisterGRPCServices(grpcServer, db)

	lis, err := net.Listen("tcp", ":"+cfg.Server.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("User Service running on port %s", cfg.Server.Port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
