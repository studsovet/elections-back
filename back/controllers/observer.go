package controllers

import (
	db "elections-back/db"
	token "elections-back/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetPrivateKey(c *gin.Context) {
	var input db.Key

	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := db.GetUserByID(userId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !user.IsObserver {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only observers can upload tokens!"})
		return
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Type == "public" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can't change public key"})
		return
	}
	if input.PartID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can't change composed private key"})
		return
	}

	f, err := input.IsTokenExist()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !f {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This token doesn't exist!"})
		return
	}

	input.SaveKey()

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
