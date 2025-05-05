package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		log.Printf(
			"%s %s %s %d %s",
			c.Request.Method,
			c.Request.RequestURI,
			c.ClientIP(),
			c.Writer.Status(),
			time.Since(start),
		)
	}
}
