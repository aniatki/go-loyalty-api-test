package main

import (
	"log"
	"os"

	"github.com/aniatki/loyalty-api/internal/config"
	"github.com/aniatki/loyalty-api/internal/handlers"
	"github.com/aniatki/loyalty-api/internal/repositories"
	"github.com/aniatki/loyalty-api/internal/routes"
	"github.com/aniatki/loyalty-api/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
