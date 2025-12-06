package imports

import (
	"net/http"

	"bober.app/internal/db"
	"bober.app/internal/jsonRespond"
	"github.com/gin-gonic/gin"

	"github.com/xuri/excelize/v2"
)

func ImportXLSX(c *gin.Context) {
	_, err := db.OpenDB()
	if err != nil {
		jsonRespond.Error(c, http.StatusInternalServerError, "could not connect to database", err)
		return
	}

	// 2. Pobranie pliku z form-data
	fileHeader, err := c.FormFile("file")
	if err != nil {
		jsonRespond.Error(c, http.StatusBadRequest, "file is required", err)
		return
	}

	// 3. Otwieramy io.Reader z przesłanego pliku
	file, err := fileHeader.Open()
	if err != nil {
		jsonRespond.Error(c, http.StatusBadRequest, "could not open uploaded file", err)
		return
	}
	defer file.Close()

	// 4. Wczytujemy Excela z readera
	xls, err := excelize.OpenReader(file)
	if err != nil {
		jsonRespond.Error(c, http.StatusBadRequest, "invalid excel file", err)
		return
	}
	defer func() {
		_ = xls.Close()
	}()

	// 5. Przykład: odczyt wierszy z pierwszego arkusza
	sheetName := xls.GetSheetName(0)
	if sheetName == "" {
		jsonRespond.Error(c, http.StatusBadRequest, "excel has no sheets", nil)
		return
	}

	rows, err := xls.GetRows(sheetName)
	if err != nil {
		jsonRespond.Error(c, http.StatusInternalServerError, "could not read rows", err)
		return
	}

	// Tu robisz swoją logikę: insert do bazy, mapowanie, etc.
	// Przykładowo tylko zwrócę ilość wierszy:
	jsonRespond.SendJSON(c, gin.H{
		"status":     "OK",
		"sheet":      sheetName,
		"rows_count": len(rows),
	})
}
