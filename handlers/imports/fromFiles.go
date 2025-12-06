package imports

import (
	"log"
	"net/http"
	"os"

	"bober.app/internal/db"
	"bober.app/internal/jsonRespond"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func ImportFilesAI(c *gin.Context) {
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
	jsonRespond.SendJSON(c, gin.H{
		"status": "OK",
	})
}
