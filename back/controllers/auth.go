package controllers

import (
	token "elections-back/utils"
	"encoding/base64"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"strings"

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
	redirect := c.Query("redirect_uri")
	if redirect == "" {
		redirect = os.Getenv("DEFAULT_REDIRECT")
	} else {
		if os.Getenv("ALLOWED_REDIRECTS") != "" {
			allowedRedirects := os.Getenv("ALLOWED_REDIRECTS")
			allowedRedirectsSlice := strings.Split(allowedRedirects, ",")
			allowed := false
			for _, allowedRedirect := range allowedRedirectsSlice {
				if strings.TrimSpace(allowedRedirect) == redirect {
					allowed = true
					break
				}
			}
			if !allowed {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid redirect"})
				return
			}
		}
	}
	state := GenerateState(RouterState{
		ServiceID: os.Getenv("SERVICE_ID"),
		StateData: map[string]interface{}{
			"redirect_uri": redirect,
		},
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
	c.SetCookie("token", input.AccessToken, 60*60, "/", "elections-api.studsovet.me", false, false)
	stateAsJson, err := base64.StdEncoding.DecodeString(input.State)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}
	var state map[string]string
	err = json.Unmarshal(stateAsJson, &state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}
	redirectUri, exists := state["redirect_uri"]
	if !exists {
		redirectUri = os.Getenv("DEFAULT_REDIRECT")
	}
	c.Redirect(302, redirectUri+"?token="+input.AccessToken)
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
