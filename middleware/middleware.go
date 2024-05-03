package middleware

import (
	"fmt"
	customerrors "jwt_use/customErrors"
	"jwt_use/tokens"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip authentication for signup and login endpoints
		if strings.HasPrefix(c.Request.URL.Path, "/auth/login") || strings.HasPrefix(c.Request.URL.Path, "/user/signup") {
			log.Print("chnage from here")
			c.Next()
			return
		}

		// Check for token in header
		ClientToken := c.Request.Header.Get("token")
		fmt.Println(ClientToken)
		if ClientToken == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": customerrors.ErrTokenEmpty})
			c.Abort()
			return
		}

		// Validate token
		claims, err := tokens.ValidateToken(ClientToken)
		if err != "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			c.Abort()
			return
		}

		// Set claims in context for later use
		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)
		c.Next()
	}
}
