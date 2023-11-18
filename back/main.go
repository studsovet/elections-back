package main

import (
	middlewares "elections-back/middleware"
	"net/http"

	"github.com/gin-gonic/gin"

	"elections-back/controllers"
	db "elections-back/db"
)

func main() {
	db.ConnectDB()

	r := gin.Default()
	r.Use(middlewares.TokenAuthMiddleware)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Old handlers
	/*
		r.POST("/register", controllers.Register)
		r.POST("/login", controllers.Login)

		protected := r.Group("/election")
		protected.Use(middlewares.JwtAuthMiddleware())

		adminGroup := protected.Group("/admin")
		adminGroup.Use(middlewares.AdminAuthMiddleware())
		adminGroup.POST("/start", controllers.ElectionStart)
		adminGroup.POST("/stop", controllers.ElectionStop)
		adminGroup.POST("/recovery", controllers.PrivateKeyRecovery)
		adminGroup.POST("/addobserver", controllers.AddObserver)
		adminGroup.POST("/addcandidate", controllers.AddCandidate)

		observerGroup := protected.Group("/observer")
		observerGroup.Use(middlewares.ObserverAuthMiddleware())
		observerGroup.POST("/setkey", controllers.SetPrivateKey)

		voteGroup := protected.Group("/vote")
		voteGroup.GET("/public.pem", controllers.GetPublicKey)
		voteGroup.POST("/vote", controllers.PostVote)
		voteGroup.POST("/voteencrypted", controllers.PostEncryptedVote)
		voteGroup.GET("/result", controllers.ElectionResult)
		voteGroup.GET("/getcandidates", controllers.GetCandidates)
	*/

	// New handlers

	authGroup := r.Group("/auth")

	// auth
	authGroup.GET("/elk", controllers.RedirectToELK)
	authGroup.POST("/redirect", controllers.Login)
	authGroup.GET("/me", controllers.GetMe)

	electionsGroup := r.Group("/elections")

	// elector
	electionsGroup.POST("/becomeCandidate/:electionId", controllers.BecomeCandidate)
	electionsGroup.GET("/myCandidateStatus/:electionId", controllers.MyCandidateStatus)
	electionsGroup.GET("/all", controllers.GetFilteredElections)
	electionsGroup.GET("/get/:electionId", controllers.GetElection)
	electionsGroup.GET("/getCandidates/:electionId", controllers.GetCandidates)
	electionsGroup.GET("/getVoices/:electionId", controllers.GetEncryptedVotes)
	electionsGroup.GET("/getResults/:electionId", controllers.ElectionNotImplemented)
	electionsGroup.POST("/vote/:electionId", controllers.PostVote)
	electionsGroup.GET("/publicKey/:electionId", controllers.GetPublicKey)

	// observer
	electionsGroup.POST("/setPrivateKey/:electionId", controllers.PostSavePrivateKey)

	// admin
	electionsGroup.POST("/setPublicKey/:electionId", controllers.SetPublicKey)
	electionsGroup.POST("/create", controllers.CreateElection)
	electionsGroup.POST("/approveCandidate/:electionId/:candidateId", controllers.ApproveCandidate)
	electionsGroup.POST("/next/:electionId", controllers.Next)
	electionsGroup.GET("/getAllCandidates", controllers.GetAllCandidates)

	dictionariesGroup := r.Group("/dictionaries")

	// dictionary
	dictionariesGroup.GET("/councilOrganizations", controllers.GetCouncilOrganizations)
	dictionariesGroup.GET("/faculty", controllers.GetFaculty)
	dictionariesGroup.GET("/dormitory", controllers.GetDormitory)

	r.Run()
}
