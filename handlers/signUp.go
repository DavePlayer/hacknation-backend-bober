package handlers

import (
	"net/http"

	"bober.app/internal/db"
	"bober.app/internal/jsonRespond"
	"bober.app/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	//get email/pass of body
	var body struct {
		Email    string
		Password string
	}
	if c.Bind(&body) != nil {
		jsonRespond.Error(c, http.StatusBadRequest, "Failed to get body", nil)
		return
	}
	//hass pass
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		jsonRespond.Error(c, http.StatusBadRequest, "Failed to generate hash", err)
		return
	}

	// create user
	user := models.User{Email: body.Email, Password: string(hash)}
	db, err := db.OpenDB()
	if err != nil {
		jsonRespond.Error(c, http.StatusBadRequest, "Failed to open db!", err)
		return
	}
	result := db.Create(&user)

	if result.Error != nil {
		jsonRespond.Error(c, http.StatusBadRequest, "Failed to Create user!", result.Error)
		return
	}

	// Respond
	ru := models.ReturnedUser{}.From(user)

	jsonRespond.SendJSON(c, ru)
}
