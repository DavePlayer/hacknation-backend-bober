package v1

import (
	"bober.app/handlers"
	"bober.app/handlers/imports"
	"github.com/gin-gonic/gin"
)

func Register(r *gin.RouterGroup) {
	r.GET("/", handlers.HandlePing)
	r.POST("/register", handlers.SignUp)
	r.POST("/login", handlers.LoginHandler)

	// imports
	r.POST("/import/xlsl", imports.ImportXLSX)
}
