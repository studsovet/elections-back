package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
type VoteInput struct {
	Vote        interface{} `bson:"vote" json:"vote" bindings:"required"`
	BallotBoxID int         `bson:"ballotid" json:"ballotid" bindings:"required"`
}

func GetPublicKey(c *gin.Context) {
	key, err := db.GetPublicKey()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bKey, err := hex.DecodeString(key.Data)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pKey := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: bKey,
	}

	pem.Encode(c.Writer, pKey)
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	key, err := db.GetParsedPublicKey()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	vote.BallotBoxID = input.BallotBoxID
	vote.Data = hex.EncodeToString(voteEncoded)

	c.JSON(http.StatusOK, gin.H{"message": "OK"})

	// Save vote
	vote.SaveVote()
	db.AddUserVote(userId, input.BallotBoxID)
}

func PostEncryptedVote(c *gin.Context) {
	var vote db.Vote

	// Try to parse user data
	if err := c.ShouldBindJSON(&vote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check: is user already voted?
	userId, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	f, err := db.IsUserVoted(userId, vote.BallotBoxID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if f {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This user is already voted!"})
		return
	}

	// Save vote
	vote.SaveVote()
	db.AddUserVote(userId, vote.BallotBoxID)
}

func ElectionResult(c *gin.Context) {
	status := db.GetLastStatus()

	if status.Code != 3 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Votes are encoded!"})
		return
	}

	// TODO: count votes
	res, err := db.GetVotes()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"votes": res})
}
*/
func ElectionNotImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented!"})
}
