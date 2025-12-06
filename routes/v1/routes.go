package v1

import (
	"bober.app/handlers"
	"bober.app/handlers/imports"
	"bober.app/handlers/item"
	"bober.app/middleware"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup) {
	r.GET("/", handlers.HandlePing)
	r.POST("/register", handlers.SignUp)
	r.POST("/login", handlers.LoginHandler)

	r.POST("/item", middleware.Auth(), item.CreateItem)
	r.DELETE("/item/:id", middleware.Auth(), item.DeleteItem)
	r.PUT("/item/:id", middleware.Auth(), item.UpdateItem)
	r.GET("/item/:id", item.ReadItem)
	r.GET("/items/:page", item.ReadItems)
	// imports
	r.POST("/import/xlsl", imports.ImportXLSX)
}
