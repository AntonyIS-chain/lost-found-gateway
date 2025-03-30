package postgresDB

import (
	"fmt"
	"time"

	"github.com/AntonyIS-chain/lost-found-gateway/config"
	"github.com/AntonyIS-chain/lost-found-gateway/internal/core/domain"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDBClient struct {
	DB *gorm.DB
}

// NewPostgresDBClient creates a single reusable database connection.
func NewPostgresDBClient(appConfig *config.Config) (*PostgresDBClient, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=require",
		appConfig.POSTGRES_HOST, appConfig.POSTGRES_PORT,
		appConfig.POSTGRES_USER, appConfig.POSTGRES_DB, appConfig.POSTGRES_PASSWORD,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable connection pooling
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection pool: %w", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	// Run migrations once
	err = db.AutoMigrate(&domain.UserToken{})
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return &PostgresDBClient{DB: db}, nil
}
