package main

import (
	"gin_jwt/controllers"
	"gin_jwt/middlewares"
	"gin_jwt/models"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(gin.Logger())

	r.Static("/public", "./public")

	r.LoadHTMLGlob("templates/**/*")

	models.ConnectDatabase()
	models.DBMigrate()

	r.GET("/blogs", middlewares.AuthMiddleware, controllers.BlogsIndex)
	r.GET("/blogs/:id", middlewares.AuthMiddleware, controllers.BlogsShow)

	r.GET("/signup", controllers.SignupPage)
	r.GET("/login", controllers.LoginPage)
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)
	r.DELETE("/logout", controllers.Logout)

	log.Println("Server started!")
	r.Run() // Default Port 8080
}
