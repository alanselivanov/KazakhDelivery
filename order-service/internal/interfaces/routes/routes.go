package routes

import (
	"order-service/internal/application"
	"order-service/internal/infrastructure/database"
	"order-service/internal/infrastructure/messaging"
	"order-service/internal/infrastructure/persistence"
	"order-service/internal/interfaces/handlers"

	"proto/order"

	"google.golang.org/grpc"
)

func RegisterGRPCServices(grpcServer *grpc.Server, db *database.MongoDBConnector, publisher messaging.EventPublisher) {
	orderRepo := persistence.NewMongoOrderRepository(db)

	orderUseCase := application.NewOrderUseCase(orderRepo, publisher)

	orderHandler := handlers.NewOrderHandler(orderUseCase)

	order.RegisterOrderServiceServer(grpcServer, orderHandler)
}
