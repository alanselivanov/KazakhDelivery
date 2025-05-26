package routes

import (
	"log"
	"user-service/internal/application"
	"user-service/internal/config"
	"user-service/internal/infrastructure/database"
	"user-service/internal/infrastructure/mail"
	"user-service/internal/infrastructure/persistence"
	"user-service/internal/interfaces/handlers"

	"proto/user"

	"google.golang.org/grpc"
)

type Services struct {
	RedisCache  *database.RedisCache
	MailService *mail.MailService
}

func RegisterGRPCServices(grpcServer *grpc.Server, db *database.MongoDB, cfg *config.Config) *Services {
	redisCache, err := database.NewRedisCache(cfg)
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v. Continuing without caching.", err)
		redisCache = nil
	} else {
		log.Println("Successfully connected to Redis")
	}

	mailService := mail.NewMailService(cfg)
	if cfg.SMTP.Username == "" || cfg.SMTP.Password == "" {
		log.Println("Warning: SMTP credentials not provided. Email functionality will be disabled.")
		mailService = nil
	} else {
		log.Println("Mail service configured successfully")
	}

	userRepo := persistence.NewMongoUserRepository(db)

	userUseCase := application.NewUserUseCase(userRepo, redisCache, mailService)

	userHandler := handlers.NewUserHandler(userUseCase)

	user.RegisterUserServiceServer(grpcServer, userHandler)

	return &Services{
		RedisCache:  redisCache,
		MailService: mailService,
	}
}

func RegisterGRPCServicesWithInMemoryDB(grpcServer *grpc.Server, db *database.InMemoryDB) *Services {
	userRepo := persistence.NewUserRepository(db)

	userUseCase := application.NewUserUseCase(userRepo, nil, nil)

	userHandler := handlers.NewUserHandler(userUseCase)

	user.RegisterUserServiceServer(grpcServer, userHandler)

	return &Services{
		RedisCache:  nil,
		MailService: nil,
	}
}
