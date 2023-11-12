package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCouncilOrganizations(c *gin.Context) {
	var organizations []string = []string{"Консульская организация"}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": organizations})
}

func GetFaculty(c *gin.Context) {
	var faculties []string = []string{"Факультет"}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": faculties})
}

func GetDormitory(c *gin.Context) {
	var dormitories []string = []string{"Общежитие"}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": dormitories})
}
