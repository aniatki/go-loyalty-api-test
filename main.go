package main

import (
	"fmt"
	"net/http"
	"os"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Item struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Tags        []Tag   `json:"tags" gorm:"many2many:item_tags"`
}

type Tag struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type CreateItemInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Tags        []Tag   `json:"tags"`
}

func InitDB() {
	dsn := ftm.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", 
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL_MODE"),
		)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	DB = db
	fmt.Println("Connected to database")

	// AutoMigrate will create tables and add missing columns
	err = DB.AutoMigrate(&Item{}, &Tag{})
	if err != nil {
		panic("Failed to migrate database")
	}
}

func main() {
	InitDB()

	r := gin.Default()

	// Items routes
	r.POST("/items", createItem)
	r.GET("/items", getItems)

	// Tags routes
	r.GET("/tags", getTags)
	r.POST("/tags", createTag)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = 8080
	}

	r.Run(":" + port)
}

// Handler functions

func createItem(c *gin.Context) {
	var input CreateItemInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item := Item{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
		Tags:        input.Tags,
	}

	if err := DB.Create(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func getItems(c *gin.Context) {
	var items []Item
	if err := DB.Preload("Tags").Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch items"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func getTags(c *gin.Context) {
	var tags []Tag
	if err := DB.Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}
	c.JSON(http.StatusOK, tags)
}

func createTag(c *gin.Context) {
	var tag Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(http.StatusOK, tag)
}
