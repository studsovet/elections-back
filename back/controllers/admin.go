package controllers

import (
	"elections-back/db"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateElection(c *gin.Context) {
	var election db.Election
	if err := c.ShouldBindJSON(&election); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	election.ID = uuid.New().String()
	election.Save()
	c.JSON(http.StatusOK, gin.H{"message": "success", "election": election})
}

func SetPublicKey(c *gin.Context) {
	var id db.ElectionId
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var public_key db.PublicKey
	public_key.Key = c.Query("key")
	if public_key.Key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provide `key` in query"})
		return
	}
	public_key.ID = id.ID
	print("key", public_key.ID, public_key.Key)
	public_key.Save()
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func GetAllCandidates(c *gin.Context) {
	candidates, err := db.GetAllCandidates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}
	c.JSON(http.StatusOK, candidates)
}

func ApproveCandidate(c *gin.Context) {
	var id db.CandidateId
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	approved_str := c.Query("approved")
	if approved_str == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provide `approve` param to query"})
		return
	}
	approved, err := strconv.ParseBool(approved_str)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "approved query parse error: " + err.Error()})
		return
	}
	err = db.ApproveCandidate(id.ID, approved)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func Next(c *gin.Context) {
	var id db.ElectionId
	if err := c.ShouldBindUri(&id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := db.ElectionNext(id.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "approved query parse error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
