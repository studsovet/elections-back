package controllers

import (
	db "elections-back/db"
	token "elections-back/utils"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AuthorizationCallback struct {
	AccessToken string `form:"access_token" binding:"required"`
	ExpiresIn   int    `form:"expires_in" binding:"required"`
	State       string `form:"state" binding:"required"`
	TokenType   string `form:"token_type" binding:"required"`
}

type RouterState struct {
	ServiceID string                 `json:"service_id"`
	StateData map[string]interface{} `json:"state_data"`
}

func GenerateState(state RouterState) string {
	encodedState, _ := json.Marshal(state)
	encoded := base64.StdEncoding.EncodeToString(encodedState)
	return encoded
}

func RedirectToELK(c *gin.Context) {
	state := GenerateState(RouterState{
		ServiceID: os.Getenv("SERVICE_ID"),
		StateData: map[string]interface{}{},
	})

	c.Redirect(302,
		"https://auth.hse.ru/adfs/oauth2/authorize/?"+
			"client_id="+os.Getenv("CLIENT_ID")+"&"+
			"response_type=token&"+
			"redirect_uri=https://dc.studsovet.me/redirect&state="+state+"&"+
			"response_mode=form_post")
}

func Login(c *gin.Context) {
	var input AuthorizationCallback
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := token.VerifyHSEToken(input.AccessToken) // token
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Then use token.Headers to get user data
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func GetCurrentUser(c *gin.Context) (db.User, error) {
	user_id, err := token.ExtractTokenID(c)
	if err != nil {
		return db.User{}, err
	}

	u, err := db.GetUserByID(user_id)

	if err != nil {
		return db.User{}, err
	}
	return u, nil
}
