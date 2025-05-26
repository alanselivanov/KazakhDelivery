package integration_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"advenced_alan_not_copy/tests/shared/utils"
)

type Product struct {
	ID          string  `bson:"_id" json:"id"`
	Name        string  `bson:"name" json:"name"`
	Description string  `bson:"description" json:"description"`
	Price       float64 `bson:"price" json:"price"`
	Stock       int     `bson:"stock" json:"stock"`
	CategoryID  string  `bson:"category_id" json:"category_id"`
}

func startMongoContainer(ctx context.Context) (testcontainers.Container, string, error) {
	mongoContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo:latest",
			ExposedPorts: []string{"27017/tcp"},
			WaitingFor:   wait.ForLog("Waiting for connections"),
		},
		Started: true,
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to start MongoDB container: %w", err)
	}

	mongoHost, err := mongoContainer.Host(ctx)
	if err != nil {
		return mongoContainer, "", err
	}

	mongoPort, err := mongoContainer.MappedPort(ctx, "27017/tcp")
	if err != nil {
		return mongoContainer, "", err
	}

	connectionString := fmt.Sprintf("mongodb://%s:%s", mongoHost, mongoPort.Port())

	return mongoContainer, connectionString, nil
}

func TestMongoIntegration(t *testing.T) {
	dockerAvailable, dockerMessage := utils.IsDockerAvailable()
	t.Logf("Docker status: %s", dockerMessage)

	if !dockerAvailable {
		t.Skip("Skipping integration test: Docker is not available")
	}

	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	t.Logf("Starting MongoDB container for integration test")
	mongoContainer, connString, err := startMongoContainer(ctx)
	require.NoError(t, err)

	defer func() {
		t.Logf("Terminating MongoDB container")
		if err := mongoContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connString))
	require.NoError(t, err)

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			t.Fatalf("failed to disconnect MongoDB client: %s", err)
		}
	}()

	err = client.Ping(ctx, nil)
	require.NoError(t, err)

	dbName := "inventory_test"
	collection := client.Database(dbName).Collection("products")

	t.Run("Create Product", func(t *testing.T) {
		product := Product{
			ID:          "prod-1",
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       29.99,
			Stock:       100,
			CategoryID:  "category-1",
		}

		_, err := collection.InsertOne(ctx, product)
		require.NoError(t, err)

		var result Product
		err = collection.FindOne(ctx, bson.M{"_id": "prod-1"}).Decode(&result)
		require.NoError(t, err)

		assert.Equal(t, product.ID, result.ID)
		assert.Equal(t, product.Name, result.Name)
		assert.Equal(t, product.Stock, result.Stock)
	})

	t.Run("Update Product Stock", func(t *testing.T) {
		update := bson.M{
			"$inc": bson.M{"stock": -10},
		}
		_, err := collection.UpdateOne(ctx, bson.M{"_id": "prod-1"}, update)
		require.NoError(t, err)

		var result Product
		err = collection.FindOne(ctx, bson.M{"_id": "prod-1"}).Decode(&result)
		require.NoError(t, err)

		assert.Equal(t, 90, result.Stock)
	})

	t.Run("Query Products By Category", func(t *testing.T) {
		products := []interface{}{
			Product{
				ID:          "prod-2",
				Name:        "Test Product 2",
				Description: "Another test product",
				Price:       19.99,
				Stock:       50,
				CategoryID:  "category-1",
			},
			Product{
				ID:          "prod-3",
				Name:        "Test Product 3",
				Description: "Yet another test product",
				Price:       49.99,
				Stock:       25,
				CategoryID:  "category-1",
			},
		}

		_, err := collection.InsertMany(ctx, products)
		require.NoError(t, err)

		cursor, err := collection.Find(ctx, bson.M{"category_id": "category-1"})
		require.NoError(t, err)

		var results []Product
		err = cursor.All(ctx, &results)
		require.NoError(t, err)

		assert.Equal(t, 3, len(results))
	})

	t.Run("Delete Product", func(t *testing.T) {
		_, err := collection.DeleteOne(ctx, bson.M{"_id": "prod-2"})
		require.NoError(t, err)

		count, err := collection.CountDocuments(ctx, bson.M{"_id": "prod-2"})
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)

		count, err = collection.CountDocuments(ctx, bson.M{})
		require.NoError(t, err)
		assert.Equal(t, int64(2), count)
	})

	t.Run("Atomic Stock Update", func(t *testing.T) {
		filter := bson.M{
			"_id":   "prod-1",
			"stock": bson.M{"$gte": 20},
		}
		update := bson.M{
			"$inc": bson.M{"stock": -20},
		}

		result, err := collection.UpdateOne(ctx, filter, update)
		require.NoError(t, err)
		assert.Equal(t, int64(1), result.ModifiedCount)

		var updatedProduct Product
		err = collection.FindOne(ctx, bson.M{"_id": "prod-1"}).Decode(&updatedProduct)
		require.NoError(t, err)
		assert.Equal(t, 70, updatedProduct.Stock)
	})
}
