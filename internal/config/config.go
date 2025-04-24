package config

import (
	"fmt"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
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
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	if err := db.AutoMigrate(&models.Item{}, &models.Tag{}); err != nil {
		if strings.Contains(err.Error(), "23505") {
			db.Exec("DELETE FROM tags WHERE id NOT IN (SELECT MIN(id) FROM tags GROUP BY name)")
			db.AutoMigrate(&models.Tag{})
		} else {
			panic(fmt.Sprintf("Migration failed: %v", err))
		}
	}

	return db
}
