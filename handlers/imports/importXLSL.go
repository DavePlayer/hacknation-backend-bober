package imports

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"bober.app/internal/db"
	"bober.app/internal/jsonRespond"
	"bober.app/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/xuri/excelize/v2"
)

func parseDate(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("empty date string")
	}

	// 1. spróbuj normalne formaty tekstowe
	layouts := []string{
		"02.01.2006",
		"2006-01-02",
		"02-01-2006",
		"02/01/2006",

		"02-01-06",
		"02/01/06",
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
		"02.01.2006 15:04",
		"2006-01-02 15:04",
		"2006-01-02 15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, nil
		}
	}

	// 2. spróbuj, czy to serial Excela (np. "45628")
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		// false = system 1900 (Windows) – prawie na pewno to masz
		if t, err2 := excelize.ExcelDateToTime(f, false); err2 == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("cannot parse date %q", s)
}

func ImportXLSX(c *gin.Context) {
	_, err := db.OpenDB()
	if err != nil {
		jsonRespond.Error(c, http.StatusInternalServerError, "could not connect to database", err)
		return
	}

	hmacSecret := []byte(os.Getenv("SECRET"))
	if len(hmacSecret) == 0 {
		jsonRespond.Fail(c, http.StatusInternalServerError, "SECRET env is not defined", err)
		return
	}

	tokenCookie, err := c.Cookie("token")

	if err != nil {
		jsonRespond.Error(c, http.StatusInternalServerError, "Error when decoding token", err)
		return
	}

	token, err := jwt.Parse(tokenCookie, func(token *jwt.Token) (any, error) {
		return hmacSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		jsonRespond.Error(c, http.StatusInternalServerError, "Error when decoding token", err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	subVal, ok := claims["sub"].(float64)
	log.Printf("%v", claims)
	if !ok {
		jsonRespond.Error(c, http.StatusInternalServerError, "sub claim is not a string lol", nil)
		return
	}

	issuerId := uint(subVal)

	if err != nil {
		jsonRespond.Fail(c, http.StatusUnauthorized, "token cookie is not set", err)
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

	header := rows[0]
	col := map[string]int{}
	for i, h := range header {
		col[strings.TrimSpace(h)] = i
	}

	// helper do pobierania wartości z konkretnej kolumny po nagłówku
	get := func(row []string, name string) string {
		idx, ok := col[name]
		if !ok || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}

	var imported []models.ImportedItem

	for _, row := range rows[1:] { // od drugiego wiersza (po nagłówku)
		if len(row) == 0 {
			continue
		}

		item := models.ImportedItem{
			Name:        get(row, "Nazwa"),
			Type:        get(row, "Typ"),
			Description: get(row, "Opis"),

			Issue_number: get(row, "Numer sprawy"),
			Where_stored: get(row, "Gdzie przechowywany"),
			Where_found:  get(row, "Gdzie znaleziony"),
			Voivodeship:  get(row, "Województwo"),
		}

		if t, err := parseDate(get(row, "Data przekazania dokumentu")); err == nil {
			item.Document_transfer_date = t
		}
		if t, err := parseDate(get(row, "Data wpisu")); err == nil {
			item.Entry_date = t
		}
		if t, err := parseDate(get(row, "Data znalezienia")); err == nil {
			item.Found_date = t
		}

		// status – w Excelu może być np. "nowy", "zarchiwizowany", itp.
		item.Status = get(row, "Status")

		log.Printf("RAW date values: %q\n", get(row, "Data wpisu"))

		imported = append(imported, item)
	}

	dbConn, err := db.OpenDB()
	if err != nil {
		jsonRespond.Error(c, http.StatusBadRequest, "Failed to open db!", err)
		return
	}
	for _, item := range imported {
		dbConn.Create(&models.Item{
			// TODO finish this shit
			Name:                   item.Name,
			Type:                   item.Type,
			Description:            item.Description,
			Document_transfer_date: item.Document_transfer_date,
			Found_date:             item.Found_date,
			Issue_number:           item.Issue_number,
			Where_stored:           item.Where_stored,
			Where_found:            item.Where_found,
			Voivodeship:            item.Voivodeship,
			Status:                 item.Status,
			Issuer_id:              issuerId,
		})
	}

	// Tu robisz swoją logikę: insert do bazy, mapowanie, etc.
	// Przykładowo tylko zwrócę ilość wierszy:
	jsonRespond.SendJSON(c, gin.H{
		"status":     "OK",
		"sheet":      sheetName,
		"parsed":     imported,
		"rows_count": len(rows),
	})
}
