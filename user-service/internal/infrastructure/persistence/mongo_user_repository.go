package persistence

import (
	"context"
	"errors"

	"user-service/internal/domain"
	"user-service/internal/infrastructure/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoUserRepository struct {
	db *database.MongoDB
}

func NewMongoUserRepository(db *database.MongoDB) *mongoUserRepository {
	return &mongoUserRepository{db: db}
}

func (r *mongoUserRepository) Create(ctx context.Context, user *domain.User) (*domain.User, error) {
	userDTO := &database.UserDTO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
	}

	_, err := r.db.UserCollection().InsertOne(ctx, userDTO)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors.New("user already exists")
		}
		return nil, err
	}

	return user, nil
}

func (r *mongoUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var userDTO database.UserDTO

	filter := bson.M{"_id": id}
	err := r.db.UserCollection().FindOne(ctx, filter).Decode(&userDTO)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &domain.User{
		ID:        userDTO.ID,
		Username:  userDTO.Username,
		Email:     userDTO.Email,
		Password:  userDTO.Password,
		CreatedAt: userDTO.CreatedAt,
	}, nil
}

func (r *mongoUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	var userDTO database.UserDTO

	filter := bson.M{"username": username}
	err := r.db.UserCollection().FindOne(ctx, filter).Decode(&userDTO)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &domain.User{
		ID:        userDTO.ID,
		Username:  userDTO.Username,
		Email:     userDTO.Email,
		Password:  userDTO.Password,
		CreatedAt: userDTO.CreatedAt,
	}, nil
}

func (r *mongoUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var userDTO database.UserDTO

	filter := bson.M{"email": email}
	err := r.db.UserCollection().FindOne(ctx, filter).Decode(&userDTO)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}

	return &domain.User{
		ID:        userDTO.ID,
		Username:  userDTO.Username,
		Email:     userDTO.Email,
		Password:  userDTO.Password,
		CreatedAt: userDTO.CreatedAt,
	}, nil
}
