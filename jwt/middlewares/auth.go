package middlewares

import (
	"fmt"
	"gin_jwt/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *gin.Context) {
	// Retrieve the cookie from the request
	tokenStr, err := c.Cookie("Auth")
	if err != nil {
		c.HTML(http.StatusUnauthorized,
			"errors/error.tpl",
			gin.H{"error": "No auth token"})
		c.Abort()
	}

	// Extract the JWT token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// TODO: Move this to env variable.
		hmacSampleSecret := "e1bed9f5-81d7-4810-9f9b-307d2761c4d4"

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(hmacSampleSecret), nil
	})
	if err != nil {
		c.HTML(http.StatusUnauthorized,
			"errors/error.tpl",
			gin.H{"error": "Failed to parse JWT token"})
		c.Abort()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.HTML(http.StatusUnauthorized,
			"errors/error.tpl",
			gin.H{"error": "JWT Claims failed"})
		c.Abort()
	}

	// Check expiry of the token
	if claims["ttl"].(float64) < float64(time.Now().Unix()) {
		c.HTML(http.StatusUnauthorized,
			"errors/error.tpl",
			gin.H{"error": "JWT token expired!"})
		c.Abort()
	}

	// Extract the user from the token
	var user models.User
	models.DB.Where("id = ?", claims["userID"]).First(&user)

	if user.ID == 0 {
		c.HTML(http.StatusUnauthorized,
			"errors/error.tpl",
			gin.H{"error": "Could not find the user!"})
		c.Abort()
	}

	// Set the current user in the context
	c.Set("user", user)

	// Go to the next in chain
	c.Next()
}
