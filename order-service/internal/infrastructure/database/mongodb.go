package database

import (
	"context"
	"log"
	"order-service/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBConnector struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func (m *MongoDBConnector) OrderCollection() *mongo.Collection {
	return m.Database.Collection("orders")
}

func (m *MongoDBConnector) initIndexes(ctx context.Context) error {

	userIDIndex := mongo.IndexModel{
		Keys: bson.M{"user_id": 1},
	}

	_, err := m.OrderCollection().Indexes().CreateOne(ctx, userIDIndex)
	if err != nil {
		return err
	}

	return nil
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
