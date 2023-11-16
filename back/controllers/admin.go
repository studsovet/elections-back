package controllers

import (
	"elections-back/db"
	"fmt"
	"net/http"
	"strconv"

	token "elections-back/utils"

	"github.com/gin-gonic/gin"
)

func CreateElection(c *gin.Context) {
	var election db.Election
	if err := c.ShouldBindJSON(&election); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	election.Status = db.Statuses[0]
	election.Save()
	c.JSON(http.StatusOK, gin.H{"message": "success", "election": election})
}

func SetPublicKey(c *gin.Context) {
	// TODO: make it available only in the waiting state

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

	_, err := token.ParsePublicKey(public_key.Key)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprint(err)})
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
	election, err := db.GetElection(id.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	status := election.Status
	status_num := -1
	for i, s := range db.Statuses {
		if s == status {
			status_num = i
			break
		}
	}
	if status_num == -1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not such status: " + status})
		return
	}
	if status_num+1 == len(db.Statuses) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot move next last status: " + status})
		return
	}
	new_status := db.Statuses[status_num+1]

	if new_status == db.Started {
		if _, err := db.GetPublicKey(id.ID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Before moving to status " + new_status + ", add public key"})
			return
		}
	}

	if new_status == db.Decrypted {
		if _, err := db.GetPrivateKey(id.ID); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Before moving to status " + new_status + ", add private key"})
			return
		}
	}

	err = db.ElectionUpdateStatus(id.ID, new_status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "approved query parse error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "success"})
}
