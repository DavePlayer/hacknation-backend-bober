package router

import (
	v1 "bober.app/routes/v1"
	//"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	v1Group := api.Group("/v1")

	v1.Register(v1Group)

	return r
}
