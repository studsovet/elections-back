package middlewares

import (
	token "elections-back/utils"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

/*
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := token.ExtractTokenID(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		user, err := db.GetUserByID(userId)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		if !user.IsAdmin {
			c.String(http.StatusForbidden, "Forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}

func ObserverAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := token.ExtractTokenID(c)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		user, err := db.GetUserByID(userId)
		if err != nil {
			c.String(http.StatusUnauthorized, "Unauthorized")
			c.Abort()
			return
		}

		if !user.IsObserver {
			c.String(http.StatusForbidden, "Forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}
*/

func TokenAuthMiddleware(c *gin.Context) {

	c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
	c.Writer.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	c.Writer.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	if (c.Request.Method == "GET" && (c.Request.URL.Path == "/ping" || c.Request.URL.Path == "/auth/elk")) ||
		(c.Request.Method == "POST" && c.Request.URL.Path == "/auth/redirect") {
		c.Next()
		return
	}
	bearerToken := token.ExtractToken(c)
	fmt.Println(bearerToken)
	if bearerToken == "" {
		c.AbortWithStatusJSON(401, gin.H{"error": "token not found"})
		return
	}
	token, err := token.VerifyHSEToken(bearerToken)
	if err != nil {
		c.AbortWithStatusJSON(403, gin.H{"error": fmt.Sprint(err)})
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	c.Set("claims", claims)
	c.Next()
}
