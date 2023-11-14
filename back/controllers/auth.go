package controllers

import (
	token "elections-back/utils"
	"encoding/base64"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type AuthorizationCallback struct {
	AccessToken string `form:"access_token" binding:"required"`
	ExpiresIn   string `form:"expires_in" binding:"required"`
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	_, err := token.VerifyHSEToken(input.AccessToken) // token
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token invalid"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": input.AccessToken})

	// Then use token.Header to get user data
}

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func GetMe(c *gin.Context) {
	claims := c.MustGet("claims").(jwt.MapClaims)
	asJson, _ := json.Marshal(claims)
	asMap := map[string]interface{}{}
	json.Unmarshal(asJson, &asMap)
	c.JSON(http.StatusOK, gin.H{"first_name": asMap["firstname"], "last_name": asMap["lastname"], "email": asMap["email"]})
}
