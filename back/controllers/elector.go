package controllers

import (
	"fmt"
	"net/http"

	"elections-back/db"

	token "elections-back/utils"

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

	voter_id, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	voted, err := db.IsVoted(election.ID, voter_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	election.IsVoted = voted

	c.JSON(http.StatusOK, election)
}

func GetFilteredElections(c *gin.Context) {
	return_voted := false
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provide `status` param to query"})
		return
	} else if status == db.Voted {
		return_voted = true
		status = db.Started
	}
	elections, err := db.GetElections(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	voter_id, err := token.ExtractTokenID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	elections_id := []string{}
	for _, e := range elections {
		if status != db.Started {
			elections_id = append(elections_id, e.ID)
		} else {
			voted, err := db.IsVoted(e.ID, voter_id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if (return_voted && voted) || (!return_voted && !voted) {
				elections_id = append(elections_id, e.ID)
			}
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

	if vote.VoterID, err = token.ExtractTokenID(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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
