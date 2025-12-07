package item

import (
	"net/http"
	"time"

	"bober.app/internal/db"
	"bober.app/internal/jsonRespond"
	"bober.app/models"
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2/log"
)

func CreateItem(c *gin.Context) {

	var body struct {
		IssuerID             uint      `json:"issuerId"`
		Name                 string    `json:"itemName"`
		Type                 string    `json:"type"`
		Description          string    `json:"description"`
		DocumentTransferDate time.Time `json:"documentTransferDate"`
		EntryDate            time.Time `json:"entryDate"`
		FoundDate            time.Time `json:"foundDate"`
		IssueNumber          string    `json:"issueNumber"`
		WhereStored          string    `json:"whereStorred"`
		WhereFound           string    `json:"whereFound"`
		Voivodeship          string    `json:"voivodeship"`
	}

	if c.Bind(&body) != nil {
		log.Info("body: %v", body)
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
