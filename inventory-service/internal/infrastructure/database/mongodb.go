package database

import (
	"context"
	"inventory-service/internal/config"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConnector struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func (m *MongoDBConnector) ProductCollection() *mongo.Collection {
	return m.Database.Collection("products")
}

func (m *MongoDBConnector) CategoryCollection() *mongo.Collection {
	return m.Database.Collection("categories")
}

func (m *MongoDBConnector) initIndexes(ctx context.Context) error {
	productNameIndex := mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	}

	categoryProductIndex := mongo.IndexModel{
		Keys: bson.M{"category_id": 1},
	}

	_, err := m.ProductCollection().Indexes().CreateMany(ctx, []mongo.IndexModel{
		productNameIndex,
		categoryProductIndex,
	})
	if err != nil {
		return err
	}

	categoryNameIndex := mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err = m.CategoryCollection().Indexes().CreateMany(ctx, []mongo.IndexModel{
		categoryNameIndex,
	})

	return err
}

func NewMongoDB(cfg *config.Config) (*MongoDBConnector, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.MongoDB.Timeout)*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoDB.URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to MongoDB")

	db := client.Database(cfg.MongoDB.Database)

	mongodb := &MongoDBConnector{
		Client:   client,
		Database: db,
	}

	if err := mongodb.initIndexes(ctx); err != nil {
		log.Printf("Warning: failed to create indexes: %v", err)
	} else {
		log.Println("MongoDB indexes created successfully")
	}

	return mongodb, nil
}

func (m *MongoDBConnector) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}
