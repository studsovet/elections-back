package utils

import (
	"elections-back/db"
	"errors"
	"os"
	"strings"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

func ExtractTokenEmail(c *gin.Context) (string, error) {
	tokenString := ExtractToken(c)
	token, err := VerifyHSEToken(tokenString)

	if err != nil {
		return "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		id, ok := claims["email"]
		if !ok {
			return "", errors.New("wrong token")
		}
		switch idType := id.(type) {
		case string:
			uid := string(idType)
			return uid, nil
		}
		return "", errors.New("wrong token")
	}
	return "", errors.New("wrong token")
}

func ExtractTokenID(c *gin.Context) (string, error) {
	email, err := ExtractTokenEmail(c)
	if err != nil {
		return "", err
	}
	ID, err := db.Email2ID(email)
	if err != nil {
		return "", err
	}
	return ID, nil
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
