package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

const apiTokenHeader = "X-API-Token"

func ApiTokenAuth(authToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiToken := c.GetHeader(apiTokenHeader)

		if apiToken != authToken {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		c.Next()
	}
}

func AllowOriginMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		allowHeaders := "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Api-Token"

		c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Methods", "POST, PUT, PATCH, GET, DELETE")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Headers", allowHeaders)

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.WithFields(log.Fields{
			"clientIP": c.ClientIP(),
		}).Infof("[%s] %s %d", c.Request.Method, c.Request.RequestURI, c.Writer.Status())
	}
}
