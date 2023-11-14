package token

import (
	"errors"
	"github.com/MicahParks/keyfunc/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"strings"
)

var keyset *keyfunc.JWKS

func ExtractToken(c *gin.Context) string {
	cookie, err := c.Cookie("token")

	if err == nil {
		return cookie
	}

	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func VerifyHSEToken(t string) (*jwt.Token, error) {
	if keyset == nil {
		jwksUrl := "https://auth.hse.ru/adfs/discovery/keys"
		var err error
		keyset, err = keyfunc.Get(jwksUrl, keyfunc.Options{})
		if err != nil {
			return nil, err
		}
	}
	token, err := jwt.Parse(t, keyset.Keyfunc)
	if err != nil {
		return nil, err
	} else if !token.Valid {
		return nil, errors.New("invalid token")
	}
	audience := token.Claims.(jwt.MapClaims)["aud"]
	if audience == nil || audience.(string) != "microsoft:identityserver:"+os.Getenv("CLIENT_ID") {
		return nil, errors.New("invalid audience")
	}
	return token, nil
}
