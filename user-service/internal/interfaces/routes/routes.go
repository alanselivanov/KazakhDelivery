package routes

import (
	"log"
	"user-service/internal/application"
	"user-service/internal/config"
	"user-service/internal/infrastructure/database"
	"user-service/internal/infrastructure/persistence"
	"user-service/internal/interfaces/handlers"

	"proto/user"

	"google.golang.org/grpc"
)

type Services struct {
	RedisCache *database.RedisCache
}

func RegisterGRPCServices(grpcServer *grpc.Server, db *database.MongoDB, cfg *config.Config) *Services {
	redisCache, err := database.NewRedisCache(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v. Continuing without caching.", err)
		redisCache = nil
	} else {
		log.Println("Successfully connected to Redis")
	}

	userRepo := persistence.NewMongoUserRepository(db)

	userUseCase := application.NewUserUseCase(userRepo, redisCache)

	userHandler := handlers.NewUserHandler(userUseCase)

	user.RegisterUserServiceServer(grpcServer, userHandler)

	return &Services{
		RedisCache: redisCache,
	}
}

func RegisterGRPCServicesWithInMemoryDB(grpcServer *grpc.Server, db *database.InMemoryDB) *Services {
	userRepo := persistence.NewUserRepository(db)

	userUseCase := application.NewUserUseCase(userRepo, nil)

	userHandler := handlers.NewUserHandler(userUseCase)

	user.RegisterUserServiceServer(grpcServer, userHandler)

	return &Services{
		RedisCache: nil,
	}
}
