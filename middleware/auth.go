package middleware

import (
	"net/http"
	"os"

	"bober.app/internal/jsonRespond"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// if c.Request.Method == "OPTIONS" {
		// 	c.Next()
		// 	return
		// }
		// --- Pobieranie ciasteczka ------------------------------------------------
		tokenString, err := c.Cookie("token")
		if err != nil || tokenString == "" {
			jsonRespond.Fail(c, http.StatusUnauthorized, "Unauthorized access", err)
			c.Abort()
			return
		}

		// --- Parse token ----------------------------------------------------------
		secret := []byte(os.Getenv("SECRET"))

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {

			// Sprawdź, czy algorytm jest poprawny (HMAC)
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenInvalidClaims
			}

			return secret, nil
		})

		if err != nil || !token.Valid {
			jsonRespond.Fail(c, http.StatusUnauthorized, "Invalid token", nil)
			c.Abort()
			return
		}

		// --- Pobranie claims ------------------------------------------------------
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			jsonRespond.Fail(c, http.StatusUnauthorized, "Invalid claims", nil)
			c.Abort()
			return
		}

		// --- Walidacja exp --------------------------------------------------------
		sub, ok := claims["sub"].(float64)
		if !ok {
			jsonRespond.Fail(c, http.StatusUnauthorized, "Invalid sub claim", nil)
			c.Abort()
			return
		}
		c.Set("userID", int64(sub))

		// --- Przekazanie userID ---------------------------------------------------
		// if sub, ok := claims["sub"]; ok {
		// c.Set("userID", sub)
		// }

		c.Next()
	}
}
