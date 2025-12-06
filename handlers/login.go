package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var jwtKey = []byte("supersecretkey") // tajny klucz JWT

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func loginHandler(c *gin.Context) {
	var creds Credentials

	if err := c.BindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalida credentials"})
	}

}
