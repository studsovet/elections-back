package controllers

import (
	"fmt"
	"net/http"

	"elections-back/db"

	"github.com/gin-gonic/gin"
)

func GetPublicKey(c *gin.Context) {
	var election_id db.ElectionId
	if err := c.ShouldBindUri(&election_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	key, err := db.GetPublicKey(election_id.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/x-pem-file", []byte(key.Key))

}

func GetCandidates(c *gin.Context) {
	var election_id db.ElectionId
	if err := c.ShouldBindUri(&election_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	candidates, err := db.GetCandidates(election_id.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, candidates)
}

func GetElection(c *gin.Context) {
	var election_id db.ElectionId
	if err := c.ShouldBindUri(&election_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	election, err := db.GetElection(election_id.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint(err)})
		return
	}

	c.JSON(http.StatusOK, election)
}

func GetFilteredElections(c *gin.Context) {
	elections, err := db.GetElections()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var elections_id []int64
	for _, e := range elections {
		if e.Status != db.Draft { // TODO move to mongo
			elections_id = append(elections_id, e.ID)
		}
	}

	c.JSON(http.StatusOK, elections_id)
}

func PostVote(c *gin.Context) {
	var election_id db.ElectionId
	if err := c.ShouldBindUri(&election_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	election, err := db.GetElection(election_id.ID)
	if (election.Status != db.Started) || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not found"})
		return
	}

	var vote db.EncryptedVote
	if err := c.ShouldBindJSON(&vote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	vote.VoterID = "voter_id" // TODO!

	voted, err := db.IsVoted(election_id.ID, vote.VoterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if voted {
		c.JSON(http.StatusBadRequest, gin.H{"error": "already voted"})
		return
	}

	// TODO: add basic checks (eligible for voting in this election...)

	_, err = vote.Save(election_id.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func GetEncryptedVotes(c *gin.Context) {
	var election_id db.ElectionId
	if err := c.ShouldBindUri(&election_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	votes, err := db.GetEncryptedVotes(election_id.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, votes)

}

func ElectionNotImplemented(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "Not implemented!"})
}
