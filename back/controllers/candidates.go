package controllers

import (
	"elections-back/db"
	"net/http"

	token "elections-back/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func BecomeCandidate(c *gin.Context) {
	var err error
	var candidate db.Candidate

	if err = c.ShouldBindJSON(&candidate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id db.ElectionId
	if err = c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	candidate.ID = uuid.NewString()

	if candidate.UserId, err = token.ExtractTokenID(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	candidate.ElectionId = id.ID
	election, err := db.GetElection(candidate.ElectionId)
	if (election.Status != db.Waiting) || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not found"})
		return
	}

	candidate.Approved = false
	candidate.WaitingForApprove = true

	_, err = candidate.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func MyCandidateStatus(c *gin.Context) {
	var err error
	var election_id db.ElectionId
	var candidate_id string
	var candidate db.Candidate

	if err = c.ShouldBindUri(&election_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if candidate_id, err = token.ExtractTokenID(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if candidate, err = db.GetCandidate(election_id.ID, candidate_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, candidate)
}
