package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/usama1031/golang-ecommerce-store/helper"
)

func Authenicate() gin.HandlerFunc {
	return func(c *gin.Context) {

		clientToken, err := c.Cookie("token")

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Authorization Required"})
			c.Abort()
			return
		}

		claims, msg := helper.ValidateToken(clientToken)

		if msg != "" {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("uid", claims.Uid)

		c.Next()
	}
}
