package item

import (
	"net/http"
	"time"

	"bober.app/internal/db"
	"bober.app/internal/jsonRespond"
	"bober.app/models"
	"github.com/gin-gonic/gin"
)

func CreateItem(c *gin.Context) {

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
		jsonRespond.Error(c, http.StatusBadRequest, "Failed to get body", nil)
		return
	}

	dbConn, err := db.OpenDB()
	if err != nil {
		jsonRespond.Error(c, http.StatusBadRequest, "Failed to open DB!", err)
		return
	}

	item := models.Item{
		Issuer_id:              body.IssuerID,
		Name:                   body.Name,
		Type:                   body.Type,
		Description:            body.Description,
		Document_transfer_date: body.DocumentTransferDate,
		Entry_date:             body.EntryDate,
		Found_date:             body.FoundDate,
		Issue_number:           body.IssueNumber,
		Where_stored:           body.WhereStored,
		Where_found:            body.WhereFound,
		Voivodeship:            body.Voivodeship,
	}

	// Save to DB
	if err := dbConn.Create(&item).Error; err != nil {
		jsonRespond.Error(c, http.StatusBadRequest, "Failed to create item!", err)
		return
	}

	jsonRespond.SendJSON(c, item)
}
