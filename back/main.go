package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"elections-back/controllers"
	db "elections-back/db"
	middlewares "elections-back/middleware"
)

func main() {
	db.ConnectDB()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	protected_test := r.Group("/deleteme")
	protected_test.Use(middlewares.JwtAuthMiddleware())
	protected_test.GET("/observerCount", controllers.GetObserversCount)

	protected := r.Group("/election")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.POST("/start", controllers.ElectionStart)
	protected.POST("/public.pem", controllers.GetPublicKey)
	protected.POST("/setkey", controllers.SetPrivateKey)
	protected.POST("/vote", controllers.PostVote)
	protected.POST("/stop", controllers.ElectionStop)
	protected.POST("/result", controllers.ElectionResult)

	r.Run()
}
