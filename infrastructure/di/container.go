package di

import (
	"mindbridge/application/usecases"
	"mindbridge/config"
	domainRepo "mindbridge/domain/repositories"
	"mindbridge/infrastructure/repositories"
)

type Container struct {
	UserRepo    domainRepo.IUserRepository
	AuthService domainRepo.IAuthService
	AuthUseCase *usecases.AuthUseCase
	RedisClient *repositories.RedisClient
	ChatRepo    domainRepo.IChatRepository
	FileRepo    domainRepo.IFileRepository
}

func NewContainer() (*Container, error) {
	// MongoDB
	mongoClient, err := repositories.NewMongoClient(config.GetMongoURI())
	if err != nil {
		return nil, err
	}
	collection := mongoClient.Database(config.MongoDBName).Collection("users")
	userRepo := repositories.NewUserRepository(collection)

	// Redis
	redisClient, err := repositories.NewRedisClient()
	if err != nil {
		// Redis is optional; continue without it
		redisClient = nil
	}

	// JWT service
	authService := repositories.NewJWTService(config.JWTSecret, redisClient)

	// Auth use-case
	authUseCase := usecases.NewAuthUseCase(userRepo, authService, redisClient)

	// Chat & File repos
	chatRepo := repositories.NewChatRepository()
	fileRepo := repositories.NewFileRepository()

	return &Container{
		UserRepo:    userRepo,
		AuthService: authService,
		AuthUseCase: authUseCase,
		RedisClient: redisClient,
		ChatRepo:    chatRepo,
		FileRepo:    fileRepo,
	}, nil
}
