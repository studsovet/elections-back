package controllers

import (
	"elections-back/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func BecomeCandidate(c *gin.Context) {
	var candidate db.Candidate
	if err := c.ShouldBindJSON(&candidate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var id db.ElectionId
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	candidate.ID = "candidate_id" // TODO
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
	var election_id db.ElectionId
	if err := c.ShouldBindUri(&election_id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var candidate_id string = "candidate_id" // TODO
	candidate, err := db.GetCandidate(election_id.ID, candidate_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, candidate)
}
