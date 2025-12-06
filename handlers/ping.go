package handlers

import (
	"bober.app/internal/jsonRespond"
	"github.com/gin-gonic/gin"
)

func HandlePing(c *gin.Context) {
	jsonRespond.SendJSON(c, gin.H{"OK": "OK"})
}
