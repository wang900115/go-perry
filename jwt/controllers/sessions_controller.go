package controllers

import (
	"gin_jwt/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SignupPage(c *gin.Context) {
	c.HTML(http.StatusOK,
		"sessions/signup.tpl",
		gin.H{})
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK,
		"sessions/login.tpl",
		gin.H{})
}

type formData struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func Signup(c *gin.Context) {
	var data formData
	c.Bind(&data)

	// Check if the user exists already
	if !models.CheckUserAvailability(data.Email) {
		c.HTML(http.StatusBadRequest,
			"errors/error.tpl",
			gin.H{"error": "Email missing"})
		return
	}

	// Create the user
	user := models.UserCreate(data.Email, data.Password)
	if user == nil || user.ID == 0 {
		c.HTML(http.StatusBadRequest,
			"errors/error.tpl",
			gin.H{"error": "user creation failed"})
		return
	}

	// Create JWT token
	tokenString, err := createAndSignJWT(user)
	if err != nil {
		c.HTML(http.StatusBadRequest,
			"errors/error.tpl",
			gin.H{"error": "JWT creation failed"})
		return
	}

	// 2. Send the token in a cookie
	setCookie(c, tokenString)

	c.Redirect(http.StatusFound, "/blogs")
}

func Login(c *gin.Context) {
	var data formData
	c.Bind(&data)

	// Match password
	user := models.UserMatchPassword(data.Email, data.Password)
	if user.ID == 0 {
		c.HTML(http.StatusUnauthorized,
			"errors/error.tpl",
			gin.H{"error": "Unauthorized!"})
		return
	}

	// Create JWT token
	tokenString, err := createAndSignJWT(user)
	if err != nil {
		c.HTML(http.StatusBadRequest,
			"errors/error.tpl",
			gin.H{"error": "JWT creation failed"})
		return
	}

	// 2. Send the token in a cookie
	setCookie(c, tokenString)

	c.Redirect(http.StatusFound, "/blogs")
}

func Logout(c *gin.Context) {
	// Add the JWT token to the block list.
	// or change expiry time of the cookie.

	c.SetCookie("Auth", "deleted", 0, "", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

func createAndSignJWT(user *models.User) (string, error) {
	// 1. Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"ttl":    time.Now().Add(time.Hour * 24 * 100).Unix(),
	})

	// TODO: Move this to env variable.
	hmacSampleSecret := "e1bed9f5-81d7-4810-9f9b-307d2761c4d4"

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(hmacSampleSecret))
}

func setCookie(c *gin.Context, token string) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Auth", token, 3600*24*100, "", "", false, true)
}
