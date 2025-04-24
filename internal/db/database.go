package db

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL_MODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// AutoMigrate with retry logic
	if err := migrateWithRetry(db); err != nil {
		return nil, err
	}

	return db, nil
}

func migrateWithRetry(db *gorm.DB) error {
	err := db.AutoMigrate(&models.Item{}, &models.Tag{})
	if err != nil && strings.Contains(err.Error(), "23505") {
		db.Exec("DELETE FROM tags WHERE id NOT IN (SELECT MIN(id) FROM tags GROUP BY name)")
		err = db.AutoMigrate(&models.Tag{})
	}
	return err
}
