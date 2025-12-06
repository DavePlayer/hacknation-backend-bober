package middleware

import (
	"net/http"
	"os"
	"time"

	"bober.app/internal/jsonRespond"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// get cookie
		tokenString, err := c.Cookie("token")
		if err != nil {
			jsonRespond.Fail(c, http.StatusUnauthorized, "Unauthorized access", nil)
			c.Abort()
			return
		}

		// parse token
		secret := []byte(os.Getenv("SECRET"))

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// verify alg
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenInvalidClaims
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			jsonRespond.Fail(c, http.StatusUnauthorized, "Invalid token", nil)
			c.Abort()
			return
		}

		// get claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			jsonRespond.Fail(c, http.StatusUnauthorized, "Invalid token claims", nil)
			c.Abort()
			return
		}

		// sprawdź exp
		if exp, ok := claims["exp"].(float64); ok {
			if int64(exp) < time.Now().Unix() {
				jsonRespond.Fail(c, http.StatusUnauthorized, "Token expired", nil)
				c.Abort()
				return
			}
		}

		// attach user id to context
		c.Set("userID", claims["sub"])

		c.Next()
	}
}
