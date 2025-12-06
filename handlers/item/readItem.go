package item

import (
	"net/http"

	"bober.app/internal/db"
	"bober.app/internal/jsonRespond"
	"bober.app/models"
	"github.com/gin-gonic/gin"
)

func ReadItem(c *gin.Context) {

	id := c.Param("id")

	dbConn, err := db.OpenDB()
	if err != nil {
		jsonRespond.Error(c, http.StatusInternalServerError, "Failed to open DB!", err)
		return
	}

	var item models.Item
	if err := dbConn.First(&item, "id = ?", id).Error; err != nil {
		jsonRespond.Error(c, http.StatusNotFound, "Item not found", err)
		return
	}

	jsonRespond.SendJSON(c, item)
}
