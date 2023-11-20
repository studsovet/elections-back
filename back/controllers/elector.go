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
		return
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
	elections_mask := []bool{}
	for _, e := range elections {
		voted, err := db.IsVoted(e.ID, voter_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		elections_id = append(elections_id, e.ID)
		elections_mask = append(elections_mask, voted)
	}

	c.JSON(http.StatusOK, gin.H{"electionIds": elections_id, "electionMask": elections_mask})
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

	var voterID string
	if voterID, err = token.ExtractTokenID(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	voted, err := db.IsVoted(election_id.ID, voterID)
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

	db.SetVoted(election_id.ID, voterID)

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
