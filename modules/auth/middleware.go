package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthInfo struct {
	Username string `json:"username"`
}

func AuthMiddleware(authService AuthInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		token := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := authService.VerifyToken(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		authInfo := AuthInfo{
			Username: claims.Username,
		}
		c.Set("AuthInfo", authInfo)

		c.Next()
	}
}
