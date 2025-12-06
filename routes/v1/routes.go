package v1

import (
	"bober.app/handlers"
	"bober.app/handlers/item"
	"bober.app/middleware"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup) {
	r.GET("/", handlers.HandlePing)
	r.POST("/register", handlers.SignUp)
	r.POST("/login", middleware.Auth(), handlers.LoginHandler)

	r.POST("/item", item.CreateItem)
	r.DELETE("/item/:id", item.DeleteItem)
	r.PUT("/item/:id", item.DeleteItem)
}
