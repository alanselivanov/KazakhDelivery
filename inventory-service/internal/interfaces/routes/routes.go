package routes

import (
	"inventory-service/internal/application"
	"inventory-service/internal/infrastructure/database"
	"inventory-service/internal/infrastructure/persistence"
	"inventory-service/internal/interfaces/handlers"
	"proto/inventory"

	"google.golang.org/grpc"
)

// RegisterGRPCServices now uses MongoDB repositories
func RegisterGRPCServices(grpcServer *grpc.Server, db *database.MongoDB) {
	productRepo := persistence.NewMongoProductRepository(db)
	categoryRepo := persistence.NewMongoCategoryRepository(db)

	productUseCase := application.NewProductUseCase(productRepo)
	categoryUseCase := application.NewCategoryUseCase(categoryRepo)

	inventoryHandler := handlers.NewInventoryHandler(productUseCase, categoryUseCase)

	inventory.RegisterInventoryServiceServer(grpcServer, inventoryHandler)
}

// For backward compatibility with existing code
func RegisterMongoGRPCServices(grpcServer *grpc.Server, db *database.MongoDB) {
	RegisterGRPCServices(grpcServer, db)
}
