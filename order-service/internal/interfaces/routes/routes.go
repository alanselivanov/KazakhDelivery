package routes

import (
	"log"
	"order-service/internal/application"
	"order-service/internal/config"
	"order-service/internal/infrastructure/database"
	"order-service/internal/infrastructure/messaging"
	"order-service/internal/infrastructure/persistence"
	"order-service/internal/interfaces/handlers"

	"proto/order"

	"google.golang.org/grpc"
)

type Services struct {
	RedisCache *database.RedisCache
}

func RegisterGRPCServices(grpcServer *grpc.Server, db *database.MongoDBConnector, publisher messaging.EventPublisher) *Services {
	cfg := config.LoadConfig()
	redisCache, err := database.NewRedisCache(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v. Continuing without caching.", err)
		redisCache = nil
	}

	orderRepo := persistence.NewMongoOrderRepository(db)

	orderUseCase := application.NewOrderUseCase(orderRepo, publisher, redisCache)

	orderHandler := handlers.NewOrderHandler(orderUseCase)

	order.RegisterOrderServiceServer(grpcServer, orderHandler)

	return &Services{
		RedisCache: redisCache,
	}
}
