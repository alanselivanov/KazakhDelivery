package routes

import (
	"user-service/internal/application"
	"user-service/internal/infrastructure/database"
	"user-service/internal/infrastructure/persistence"
	"user-service/internal/interfaces/handlers"

	"proto/user"

	"google.golang.org/grpc"
)

func RegisterGRPCServices(grpcServer *grpc.Server, db *database.InMemoryDB) {
	userRepo := persistence.NewUserRepository(db)

	userUseCase := application.NewUserUseCase(userRepo)

	userHandler := handlers.NewUserHandler(userUseCase)

	user.RegisterUserServiceServer(grpcServer, userHandler)
}
