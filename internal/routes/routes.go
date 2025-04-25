package routes

import (
	"github.com/aniatki/loyalty-api/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, itemHandler *handlers.ItemHandler, tagHandler *handlers.TagHandler) {
	r.POST("/items", itemHandler.CreateItem)
	r.GET("/items", itemHandler.GetItems)
	r.PATCH("/items/:id/tags", itemHandler.UpdateItemTags)

	r.GET("/tags", tagHandler.GetTags)
	r.POST("/tags", tagHandler.CreateTag)
	r.DELETE("/tags/:id", tagHandler.DeleteTag)
}
