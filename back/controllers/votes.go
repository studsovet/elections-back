package controllers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"

	db "elections-back/db"
	token "elections-back/utils"

	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VoteInput struct {
	Vote        interface{} `bson:"vote" json:"vote" bindings:"required"`
	BallotBoxID int         `bson:"ballotid" json:"ballotid" bindings:"required"`
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

func GetPublicKey(c *gin.Context) {
	key, err := db.GetPublicKey()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bKey, err := hex.DecodeString(key.Data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pKey := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: bKey,
	}

	pem.Encode(c.Writer, pKey)
}

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

func PostVote(c *gin.Context) {
	var input VoteInput
	var vote db.Vote

	// Try to parse user data
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println(input.Vote)

	// Check: is election running?
	status := db.GetLastStatus()
	log.Println(status)

	if status.Code != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can't vote now!"})
		return
	}

	// Check: is user already voted?
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	f, err := db.IsUserVoted(userId, input.BallotBoxID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if f {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This user is already voted!"})
		return
	}

	// Encode user vote
	voteString, err := json.Marshal(input.Vote)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	key, err := db.GetParsedPublicKey()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	voteEncoded, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		key,
		[]byte(voteString),
		nil,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	vote.BallotBoxID = input.BallotBoxID
	vote.Data = hex.EncodeToString(voteEncoded)

	c.JSON(http.StatusOK, gin.H{"message": "OK"})

	// Save vote
	vote.SaveVote()
	db.AddUserVote(userId, input.BallotBoxID)
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

func ElectionResult(c *gin.Context) {
	// TODO: count votes
	res, err := db.GetVotes()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"votes": res})
}
