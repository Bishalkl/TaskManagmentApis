package bootstrap

import (
	config "TaskManagmentApis/configs"
	"TaskManagmentApis/internal/database"
	"context"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type AppContainer struct {
	DB           *gorm.DB
	RedisService database.RedisService
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

	return &AppContainer{
		DB:           db,
		RedisService: redisService,
	}, nil

}
