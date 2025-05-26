package routes

import (
	"context"
	"inventory-service/internal/application"
	"inventory-service/internal/infrastructure/cache"
	"inventory-service/internal/infrastructure/database"
	"inventory-service/internal/infrastructure/messaging"
	"inventory-service/internal/infrastructure/persistence"
	"inventory-service/internal/infrastructure/product"
	"inventory-service/internal/interfaces/handlers"
	"log"
	"proto/inventory"

	gogrpc "google.golang.org/grpc"
)

func RegisterGRPCServices(
	grpcServer *gogrpc.Server,
	db *database.MongoDBConnector,
	consumer messaging.EventConsumer,
	productClient product.ProductServiceClient,
	redisClient *cache.RedisClient,
) {
	mongoProductRepo := persistence.NewMongoProductRepository(db)
	categoryRepo := persistence.NewMongoCategoryRepository(db)

	var productRepo persistence.ProductRepository
	if redisClient != nil {
		productRepo = persistence.NewRedisProductRepository(mongoProductRepo, redisClient)
		log.Println("Using Redis-backed product repository with caching")
	} else {
		productRepo = mongoProductRepo
		log.Println("Using MongoDB product repository without caching")
	}

	productUseCase := application.NewProductUseCase(productRepo)
	categoryUseCase := application.NewCategoryUseCase(categoryRepo)
	metrics := application.NewMetrics()

	orderEventHandler := application.NewOrderEventHandler(productClient, metrics)

	inventoryHandler := handlers.NewInventoryHandler(productUseCase, categoryUseCase)

	inventory.RegisterInventoryServiceServer(grpcServer, inventoryHandler)

	err := consumer.SubscribeToOrderCreated(func(ctx context.Context, event *messaging.OrderCreatedEvent) error {
		metrics.IncEventsProcessed()
		return orderEventHandler.HandleOrderCreated(ctx, event)
	})

	if err != nil {
		log.Fatalf("Failed to subscribe to order.created events: %v", err)
	}

	log.Println("Successfully subscribed to order.created events")
}
