package controllers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	db "elections-back/db"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddingObserver struct {
	ID string `bson:"id" json:"id" bindings:"required"`
}

func ElectionStart(c *gin.Context) {
	// Create RSE private and public key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	publickey := &privatekey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privatekey)
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Encode keys for future use
	privateKeyHex := hex.EncodeToString(privateKeyBytes)
	publicKeyHex := hex.EncodeToString(publicKeyBytes)

	// Private key separation
	n, err := db.CountObservers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DropKeys()

	var privateParts []string = make([]string, n)
	var partSize = len(privateKeyHex) / n

	for i := 1; i < n; i++ {
		privateParts[i-1] = privateKeyHex[(i-1)*partSize : i*partSize]
		(&db.Key{
			Data:   "",
			Type:   "private",
			PartID: i,
		}).SaveKey()
	}
	privateParts[n-1] = privateKeyHex[(n-1)*partSize:]
	(&db.Key{
		Data:   "",
		Type:   "private",
		PartID: n,
	}).SaveKey()

	(&db.Key{
		Data:   publicKeyHex,
		Type:   "public",
		PartID: 0,
	}).SaveKey()

	// TODO: send parts of private key to observers
	log.Println(privateParts)

	c.JSON(http.StatusOK, gin.H{"message": "success"})

	db.ClearVotes()
	db.ClearUserVotes()
	db.SetStatus(0, "Election started")
}

func ElectionStop(c *gin.Context) {
	err := db.PrivateKeyRecovery()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success"})
	db.SetStatus(1, "Election stoped.")
	go db.DecodeVotes()
}

func AddObserver(c *gin.Context) {
	var input AddingObserver

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := db.GetUserByID(input.ID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.IsObserver = true
	user.SaveUser()
}