package controllers

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

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

	election, err := db.GetElection(election_id.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint(err)})
		return
	}

	allowed_statuses := []string{"started", "finished", "decrypted", "results"}
	if !slices.Contains(allowed_statuses, election.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't see public key in status `" + election.Status +
			"`, allowed are `" + strings.Join(allowed_statuses, ", ") + "`"})
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

	election, err := db.GetElection(election_id.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint(err)})
		return
	}

	if election.Status == db.Draft || election.Status == db.Created || election.Status == db.Waiting {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't see candidates of not started election"})
		// we don't want candidates to look up on each other in order to steal program etc...
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

	if election.Status == db.Draft {
		id, err := token.ExtractTokenID(c)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		is_admin, err := db.IsAdmin(id)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if !is_admin {
			c.JSON(http.StatusForbidden, gin.H{"error": "not admin"})
			return
		}
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

	if status == db.Draft {
		id, err := token.ExtractTokenID(c)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		is_admin, err := db.IsAdmin(id)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if !is_admin {
			c.JSON(http.StatusForbidden, gin.H{"error": "not admin"})
			return
		}
	}

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

	election, err := db.GetElection(election_id.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint(err)})
		return
	}

	allowed_statuses := []string{"started", "finished", "decrypted", "results"}
	if !slices.Contains(allowed_statuses, election.Status) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can't see votes in status `" + election.Status +
			"`, allowed are `" + strings.Join(allowed_statuses, ", ") + "`"})
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
