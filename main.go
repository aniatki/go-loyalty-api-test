package main

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strings"
)

var DB *gorm.DB

type Item struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Tags        []Tag   `json:"tags" gorm:"many2many:item_tags"`
}

type CreateItemInput struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Tags        []Tag   `json:"tags"`
}

type UpdateItemTagsInput struct {
	TagIDs []uint `json:"tag_ids"`
}

type Tag struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"unique; not null"`
}

// Main entrypoint
func InitDB() *gorm.DB {
	requiredVars := []string{"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT", "DB_SSL_MODE"}
	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			panic(fmt.Sprintf("Missing required environment variable: %s", v))
		}
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL_MODE"),
	)

	fmt.Printf("Connecting with DSN: host=%s user=%s password=**** dbname=%s port=%s sslmode=%s\n",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSL_MODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Sprintf("Failed to get database instance: %v", err))
	}

	if err := sqlDB.Ping(); err != nil {
		panic(fmt.Sprintf("Failed to ping database: %v", err))
	}

	fmt.Println("Successfully connected to database")

	if err := db.AutoMigrate(&Item{}, &Tag{}); err != nil {
		panic(fmt.Sprintf("Migration failed: %v", err))
	}

	err = db.AutoMigrate(&Item{}, &Tag{})
	if err != nil {
		if strings.Contains(err.Error(), "23505") {
			fmt.Println("Warning: Duplicate tags detected. Cleaning up...")

			db.Exec("DELETE FROM tags WHERE id NOT IN (SELECT MIN(id) FROM tags GROUP BY name)")

			if err := db.AutoMigrate(&Tag{}); err != nil {
				panic(fmt.Sprintf("Retry migration failed: %v", err))
			}
		} else {
			panic(fmt.Sprintf("Migration failed: %v", err))
		}
	}

	DB = db
	return db
}

//func ResetDB() {
//	DB := InitDB()
//	DB.Migrator().DropTable(&Item{}, &Tag{})
//	DB.AutoMigrate(&Item{}, &Tag{})
//}

func main() {
	err := godotenv.Load("C:/Users/User/OneDrive/Desktop/Repositories/loyalty-api/.env")
	if err != nil {
		log.Fatal("Error loading .env file from absolute path:", err)
	}
	InitDB()

	r := gin.Default()

	r.POST("/items", createItem)
	r.GET("/items", getItems)
	r.PATCH("/items/:id", updateItemTags)

	r.GET("/tags", getTags)
	r.POST("/tags", createTag)
	r.DELETE("/tags/:id", deleteTag)
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = string(8080)
	}

	r.Run(":" + port)
}

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
		Tags:        nil,
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

func updateItemTags(c *gin.Context) {
	var input UpdateItemTagsInput
	if err := c.ShouldBindBodyWithJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var item Item
	if err := DB.Preload("Tags").First(&item, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	var tags []Tag
	if len(input.TagIDs) > 0 {
		if err := DB.Where("id IN ?", input.TagIDs).Find(&tags).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
			return
		}
	}

	if err := DB.Model(&item).Association("Tags").Replace(tags); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tags"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tags updated successfully", "item": item})
}

func getTags(c *gin.Context) {
	var tags []Tag
	if err := DB.Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tags"})
		return
	}

	c.JSON(http.StatusOK, tags)
}

func getFormattedTag(tag string) string {
	tag = strings.ToLower(tag)
	words := strings.Fields(tag)
	return strings.Join(words, "")
}

func createTag(c *gin.Context) {
	var tag Tag

	if err := c.ShouldBindJSON(&tag); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag.Name = getFormattedTag(tag.Name)

	var existing Tag
	err := DB.Where("name = ?", tag.Name).First(&existing).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Tag already exists"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err := DB.Create(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tag"})
		return
	}

	c.JSON(http.StatusOK, tag)

}

func deleteTag(c *gin.Context) {
	id := c.Param("id")
	result := DB.Where("id = ?", id).Delete(&Tag{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}
	c.Status(http.StatusNoContent)
}
