package router

import (
	"fmt"

	v1 "bober.app/routes/v1"
	//"github.com/gin-contrib/cors"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	r := gin.Default()

	config := cors.Config{
		AllowOrigins: []string{
			"http://127.0.0.1:3000",
			"http://localhost:3000",
			"http://10.250.192.156:3000",
			"http://0.0.0.0:3000",
		},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}
	r.Use(cors.New(config))
	r.Use(func(c *gin.Context) {
		fmt.Println("REQ", c.Request.Method, c.Request.URL.Path, "Origin:", c.Request.Header.Get("Origin"))
		c.Next()
	})

	api := r.Group("/api")
	v1Group := api.Group("/v1")

	v1.Register(v1Group)

	return r
}
