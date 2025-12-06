package item

import (
	"net/http"
	"strconv"

	"bober.app/internal/db"
	"bober.app/internal/jsonRespond"
	"bober.app/models"
	"github.com/gin-gonic/gin"
)

func ReadItems(c *gin.Context) {
	pageParam := c.Param("page")

	dbConn, err := db.OpenDB()
	if err != nil {
		jsonRespond.Error(c, http.StatusInternalServerError, "Failed to open DB", err)
		return
	}

	const pageSize = 10
	page := 1
	if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
		page = p
	}

	var total int64
	if err := dbConn.Model(&models.Item{}).Count(&total).Error; err != nil {
		jsonRespond.Error(c, http.StatusInternalServerError, "Failed to count items", err)
		return
	}

	var items []models.Item
	offset := (page - 1) * pageSize

	if err := dbConn.Limit(pageSize).Offset(offset).Find(&items).Error; err != nil {
		jsonRespond.Error(c, http.StatusInternalServerError, "Failed to fetch items", err)
		return
	}
	response := gin.H{
		"total": total,
		"items": items,
	}

	jsonRespond.SendJSON(c, response)
}
