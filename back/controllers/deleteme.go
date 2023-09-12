package controllers

import (
	db "elections-back/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetObserversCount(c *gin.Context) {
	n, err := db.CountObservers()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"result":  n,
	})
}
