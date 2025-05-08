package bootstrap

import (
	config "TaskManagmentApis/configs"
	"TaskManagmentApis/internal/database"
	"TaskManagmentApis/internal/handlers"
	"TaskManagmentApis/internal/repositories"
	service "TaskManagmentApis/internal/services"
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type Handlers struct {
	Auth *handlers.AuthHandler
}

type AppContainer struct {
	DB           *gorm.DB
	RedisService database.RedisService
	Handler      Handlers
}

func InitalizeApp() (*AppContainer, error) {
	// Load Configuration
	log.Println("🔧 Loading configuration...")
	config.LoadEnv()

	// Connect to the PostgreSQL db
	log.Println("💾 Connecting to the database...")
	dbService := database.NewDBService()
	db, err := dbService.Connect()
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to connect to database: %w", err)
	}

	// Connect to Redis db
	log.Println("🔗 Connecting to Redis...")
	ctx := context.Background()
	redisService, err := database.NewRedisService(ctx)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to connect to Redis: %w", err)

	}

	// repo->service->handler

	// Initialize repo
	log.Println("📦 Initializing repositories...")
	authRepo := repositories.NewAuthRepository(db)

	// initialize service
	log.Println("🧠 Initializing services...")
	authService := service.NewAuthService(authRepo)

	// Initialize handler
	log.Println("🧠 Initializing services...")
	authHandler := handlers.NewAuthHandler(authService)

	return &AppContainer{
		DB:           db,
		RedisService: redisService,
		Handler: Handlers{
			Auth: authHandler,
		},
	}, nil

}
