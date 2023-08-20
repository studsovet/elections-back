package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"elections-back/models"
)

func post(c *gin.Context) {
	t := models.Test{}
	c.BindJSON(&t)
	t.SaveTest()
	c.JSON(http.StatusOK, gin.H{
		"message": "accepted",
	})
}

func get(c *gin.Context) {
	t := models.GetTest()
	c.JSON(http.StatusOK, gin.H{
		"message": t.Test,
	})
}

func main() {
	models.ConnectDB()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.POST("/post", post)
	r.GET("/get", get)
	r.Run()
}
