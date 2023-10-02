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
	r.POST("/login/callback", controllers.LoginCallback)
	r.GET("/votes", controllers.GetVotes)
	protected := r.Group("/post")
	protected.Use(middlewares.JwtAuthMiddleware())
	protected.POST("/vote", controllers.PostVote)
	r.Run()
}
