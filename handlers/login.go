package handlers

import (
	"net/http"
	"os"
	"time"

	"bober.app/internal/db"
	"bober.app/internal/jsonRespond"
	"bober.app/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if c.Bind(&body) != nil {
		jsonRespond.Error(c, http.StatusBadRequest, "Failed to get body", nil)
		return
	}

	dbConn, err := db.OpenDB()
	if err != nil {
		jsonRespond.Error(c, http.StatusBadRequest, "Failed to open db!", err)
		return
	}

	var user models.User
	dbConn.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		jsonRespond.Error(c, http.StatusBadRequest, "Invalid email", nil)
		return
	}

	// compare password
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)) != nil {
		jsonRespond.Error(c, http.StatusBadRequest, "Password invalid", nil)
		return
	}

	// create JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	secret := []byte(os.Getenv("SECRET"))

	tokenString, err := token.SignedString(secret)
	if err != nil {
		jsonRespond.Error(c, http.StatusInternalServerError, "Failed to sign token", err)
		return
	}

	// Save JWT in cookie
	c.SetCookie(
		"token",
		tokenString,
		60*60*24*30, // 30 dni
		"/",
		"",
		false,
		true, // httpOnly
	)

	// return success message
	jsonRespond.SendJSON(c, user.ID)
}
