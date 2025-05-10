package product

import (
	"context"
	"fmt"
	"inventory-service/internal/config"
	"inventory-service/internal/infrastructure/database"
	"inventory-service/internal/infrastructure/persistence"
	"log"
	"time"
)

type ProductServiceClient interface {
	DecreaseStock(ctx context.Context, productID string, quantity int) error
	Close()
}

type ProductClient struct {
	db          *database.MongoDBConnector
	productRepo persistence.ProductRepository
}

func NewProductServiceClient(cfg *config.Config, existingDB *database.MongoDBConnector) (ProductServiceClient, error) {
	db := existingDB
	productRepo := persistence.NewMongoProductRepository(db)

	log.Println("Using product client with MongoDB")
	return &ProductClient{
		db:          db,
		productRepo: productRepo,
	}, nil
}

func (c *ProductClient) DecreaseStock(ctx context.Context, productID string, quantity int) error {
	startTime := time.Now()
	log.Printf("[%s] Attempting to decrease stock for product %s by %d",
		startTime.Format(time.RFC3339Nano), productID, quantity)

	product, err := c.productRepo.GetByID(ctx, productID)
	if err != nil {
		log.Printf("[%s] Error getting product %s: %v [latency: %v]",
			time.Now().Format(time.RFC3339Nano), productID, err, time.Since(startTime))
		return err
	}

	if product == nil {
		err := fmt.Errorf("product not found: %s", productID)
		log.Printf("[%s] %v [latency: %v]",
			time.Now().Format(time.RFC3339Nano), err, time.Since(startTime))
		return err
	}

	if product.Stock < quantity {
		err := fmt.Errorf("insufficient stock for product %s: requested %d, available %d",
			productID, quantity, product.Stock)
		log.Printf("[%s] %v [latency: %v]",
			time.Now().Format(time.RFC3339Nano), err, time.Since(startTime))
		return err
	}

	product.Stock -= quantity
	updatedProduct, err := c.productRepo.Update(ctx, product)
	if err != nil {
		log.Printf("[%s] Error updating stock for product %s: %v [latency: %v]",
			time.Now().Format(time.RFC3339Nano), productID, err, time.Since(startTime))
		return err
	}

	endTime := time.Now()
	log.Printf("[%s] Successfully decreased stock for product %s from %d to %d [latency: %v]",
		endTime.Format(time.RFC3339Nano), updatedProduct.ID, updatedProduct.Stock+quantity,
		updatedProduct.Stock, endTime.Sub(startTime))

	return nil
}

func (c *ProductClient) Close() {
	log.Println("Product client closed")
}
