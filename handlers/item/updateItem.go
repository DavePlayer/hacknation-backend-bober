package item

import (
	"net/http"
	"time"

	"bober.app/internal/db"
	"bober.app/internal/jsonRespond"
	"bober.app/models"
	"github.com/gin-gonic/gin"
)

func UpdateItem(c *gin.Context) {
	id := c.Param("id")

	var body struct {
		IssuerID             string    `json:"issuer_id"`
		Name                 string    `json:"name"`
		Type                 string    `json:"type"`
		Description          string    `json:"description"`
		DocumentTransferDate time.Time `json:"document_transfer_date"`
		EntryDate            time.Time `json:"entry_date"`
		FoundDate            time.Time `json:"found_date"`
		IssueNumber          string    `json:"issue_number"`
		WhereStored          string    `json:"where_stored"`
		WhereFound           string    `json:"where_found"`
		Voivodeship          string    `json:"voivodeship"`
	}

	if c.Bind(&body) != nil {
		jsonRespond.Error(c, http.StatusBadRequest, "Failed to read body", nil)
		return
	}

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

	// Apply updates
	item.Issuer_id = body.IssuerID
	item.Name = body.Name
	item.Type = body.Type
	item.Description = body.Description
	item.Document_transfer_date = body.DocumentTransferDate
	item.Entry_date = body.EntryDate
	item.Found_date = body.FoundDate
	item.Issue_number = body.IssueNumber
	item.Where_stored = body.WhereStored
	item.Where_found = body.WhereFound
	item.Voivodeship = body.Voivodeship

	if err := dbConn.Save(&item).Error; err != nil {
		jsonRespond.Error(c, http.StatusBadRequest, "Failed to update item!", err)
		return
	}

	jsonRespond.SendJSON(c, "Item Updated!")
}
