package database

import (
	config "TaskManagmentApis/configs"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBService interface {
	Connect() (*gorm.DB, error)
}

type PostgresDB struct{}

func NewDBService() DBService {
	return &PostgresDB{}
}

// connect
func (p *PostgresDB) Connect() (*gorm.DB, error) {
	cfg := config.Config

	// construct the dsn
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	// Log safe message
	log.Println("📡 Attempting to connect to PostgreSQL...")

	// Try to connect
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("❌ Failed to connect to PostgreSQL:", err)
		return nil, err
	}

	// optional: Pint the databases
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("❌ Failed to get DB instance:", err)
		return nil, err
	}

	if err := sqlDB.Ping(); err != nil {
		log.Println("❌ PostgreSQL ping failed:", err)
		return nil, err
	}

	log.Println("✅ Successfully connected to PostgreSQL database.")
	return db, nil
}
