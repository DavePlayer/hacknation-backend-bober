package handlers

import (
	"bober.app/internal/db"
	"bober.app/internal/jsonRespond"
	"bober.app/models"
	"github.com/gin-gonic/gin"
)

func MeHandler(c *gin.Context) {

	userID := c.GetInt64("userID")

	var user models.User
	dbConn, _ := db.OpenDB()
	dbConn.First(&user, userID)

	returned := models.ReturnedUser{}.From(user)

	jsonRespond.SendJSON(c, returned)
}
