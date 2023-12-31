package controllers

import (
	"elections-back/db"
	"fmt"
	"io/ioutil"
	"net/http"
	"slices"
	"strings"

	token "elections-back/utils"

	"github.com/gin-gonic/gin"
)

/*
func SetPrivateKey(c *gin.Context) {
	var input db.Key

	status := db.GetLastStatus()

	if status.Code != 1 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Vote is not stopped!"})
		return
	}

	userId, err := token.ExtractTokenID(c)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := db.GetUserByID(userId)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if !user.IsObserver {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Only observers can upload tokens!"})
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

	key, err := db.GetPrivateKeyPart(input.PartID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if key.Owner != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can't change THIS public key"})
		return
	}

	input.SaveKey()

	c.JSON(http.StatusOK, gin.H{"message": "OK"})
}
*/

func PostSavePrivateKey(c *gin.Context) {
	var electionId db.ElectionId
	if err := c.ShouldBindUri(&electionId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	election, err := db.GetElection(electionId.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint(err)})
		return
	}

	allowed_statuses := []string{"finished"}
	if !slices.Contains(allowed_statuses, election.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't save private key in status `" + election.Status +
			"`, allowed are `" + strings.Join(allowed_statuses, ", ") + "`"})
		return
	}

	var private_key db.PrivateKey
	key_bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	private_key.Key = string(key_bytes[:])
	private_key.ID = electionId.ID

	privateKey, err := token.ParsePrivateKey(private_key.Key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint(err)})
		return
	}

	dbPublicKey, err := db.GetPublicKey(private_key.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint(err)})
		return
	}

	publicKey, err := token.ParsePublicKey(dbPublicKey.Key)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint(err)})
		return
	}

	if !token.IsKeyMatched(publicKey, privateKey) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "keys not matched"})
		return
	}

	private_key.Save()
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
