package routes

import (
	"order-service/internal/application"
	"order-service/internal/infrastructure/database"
	"order-service/internal/infrastructure/persistence"
	"order-service/internal/interfaces/handlers"

	"proto/order"

	"google.golang.org/grpc"
)

// RegisterGRPCServices now uses MongoDB repositories
func RegisterGRPCServices(grpcServer *grpc.Server, db *database.MongoDB) {
	orderRepo := persistence.NewMongoOrderRepository(db)

	orderUseCase := application.NewOrderUseCase(orderRepo)

	orderHandler := handlers.NewOrderHandler(orderUseCase)

	order.RegisterOrderServiceServer(grpcServer, orderHandler)
}
