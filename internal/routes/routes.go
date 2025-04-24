package routes

import (
	"github.com/gin-gonic/gin"
	"loyalty-api/internal/handlers"
)

func SetupRoutes(r *gin.Engine, itemHandler *handlers.ItemHandler, tagHandler *handlers.TagHandler) {
	r.POST("/items", itemHandler.CreateItem)
	r.GET("/items", itemHandler.GetItems)
	r.PATCH("/items/:id/tags", itemHandler.UpdateItemTags)

	r.GET("/tags", tagHandler.GetTags)
	r.POST("/tags", tagHandler.CreateTag)
	r.DELETE("/tags/:id", tagHandler.DeleteTag)
}
