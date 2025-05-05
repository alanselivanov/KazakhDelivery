package database

import (
	"context"
	"log"
	"time"
	"user-service/internal/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

func (m *MongoDB) UserCollection() *mongo.Collection {
	return m.Database.Collection("users")
}

func (m *MongoDB) initUserIndexes(ctx context.Context) error {
	usernameIndex := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}

	emailIndex := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true),
	}

	_, err := m.UserCollection().Indexes().CreateMany(ctx, []mongo.IndexModel{
		usernameIndex,
		emailIndex,
	})

	return err
}

func NewMongoDB(cfg *config.Config) (*MongoDB, error) {
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

	mongodb := &MongoDB{
		Client:   client,
		Database: db,
	}

	if err := mongodb.initUserIndexes(ctx); err != nil {
		log.Printf("Warning: failed to create indexes: %v", err)
	} else {
		log.Println("MongoDB indexes created successfully")
	}

	return mongodb, nil
}

func (m *MongoDB) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}
