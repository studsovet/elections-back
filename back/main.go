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

	protected := r.Group("/election")
	protected.Use(middlewares.JwtAuthMiddleware())

	adminGroup := protected.Group("/admin")
	adminGroup.Use(middlewares.AdminAuthMiddleware())
	adminGroup.POST("/start", controllers.ElectionStart)
	adminGroup.POST("/stop", controllers.ElectionStop)
	adminGroup.POST("/recovery", controllers.PrivateKeyRecovery)
	adminGroup.POST("/addobserver", controllers.AddObserver)

	observerGroup := protected.Group("/observer")
	observerGroup.Use(middlewares.ObserverAuthMiddleware())
	observerGroup.POST("/setkey", controllers.SetPrivateKey)

	voteGroup := protected.Group("/vote")
	voteGroup.POST("/public.pem", controllers.GetPublicKey)
	voteGroup.POST("/vote", controllers.PostVote)
	voteGroup.POST("/voteencrypted", controllers.PostEncryptedVote)
	voteGroup.POST("/result", controllers.ElectionResult)

	r.Run()
}
