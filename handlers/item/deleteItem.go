package item

import (
	"net/http"

	"bober.app/internal/db"
	"bober.app/internal/jsonRespond"
	"bober.app/models"
	"github.com/gin-gonic/gin"
)

func DeleteItem(c *gin.Context) {

	id := c.Param("id")
	if id == "" {
		jsonRespond.Error(c, http.StatusBadRequest, "Missing ID", nil)
		return
	}

	dbConn, err := db.OpenDB()
	if err != nil {
		jsonRespond.Error(c, http.StatusInternalServerError, "Failed to open DB!", err)
		return
	}

	//check if item exists
	var item models.Item
	if err := dbConn.First(&item, "id = ?", id).Error; err != nil {
		jsonRespond.Error(c, http.StatusNotFound, "Item not found", err)
		return
	}

	if err := dbConn.Delete(&item).Error; err != nil {
		jsonRespond.Error(c, http.StatusBadRequest, "Failed to delete item!", err)
		return
	}

	jsonRespond.SendJSON(c, "Item Deleted")
}
