package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"loyalty-api/internal/config"
	"loyalty-api/internal/handlers"
	"loyalty-api/internal/repositories"
	"loyalty-api/internal/routes"
	"loyalty-api/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	db := config.InitDB()

	itemRepo := repositories.NewItemRepository(db)
	tagRepo := repositories.NewTagRepository(db)

	itemService := services.NewItemService(itemRepo, tagRepo)
	tagService := services.NewTagService(tagRepo)

	itemHandler := handlers.NewItemHandler(itemService)
	tagHandler := handlers.NewTagHandler(tagService)

	r := gin.Default()
	routes.SetupRoutes(r, itemHandler, tagHandler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	r.Run(":" + port)
}
