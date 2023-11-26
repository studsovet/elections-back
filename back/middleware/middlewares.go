package middlewares

import (
	"elections-back/db"
	"elections-back/utils"
	token "elections-back/utils"
	"fmt"
	"net/http"

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
	if (c.Request.Method == "GET" && (c.Request.URL.Path == "/ping" || c.Request.URL.Path == "/auth/elk")) ||
		(c.Request.Method == "POST" && c.Request.URL.Path == "/auth/redirect") {
		c.Next()
		return
	}
	bearerToken := token.ExtractToken(c)
	//fmt.Println(bearerToken)
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
	err = utils.GetSaveStudentData(c, bearerToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprint(err)})
		return
	}
	c.Next()
}

func AdminAuthMiddleware(c *gin.Context) {
	id, err := token.ExtractTokenEmail(c)
	if err != nil {
		c.AbortWithStatusJSON(403, gin.H{"error": err.Error()})
		return
	}
	is_admin, err := db.IsAdmin(id)
	if err != nil {
		c.AbortWithStatusJSON(403, gin.H{"error": err.Error()})
		return
	}
	if !is_admin {
		c.AbortWithStatusJSON(403, gin.H{"error": "Not admin"})
	}
	c.Next()
}

func ObserverAuthMiddleware(c *gin.Context) {
	id, err := token.ExtractTokenEmail(c)
	if err != nil {
		c.AbortWithStatusJSON(403, gin.H{"error": err.Error()})
		return
	}
	is_observer, err := db.IsObserver(id)
	if err != nil {
		c.AbortWithStatusJSON(403, gin.H{"error": err.Error()})
		return
	}
	if !is_observer {
		c.AbortWithStatusJSON(403, gin.H{"error": "Not observer"})
	}
	c.Next()
}
