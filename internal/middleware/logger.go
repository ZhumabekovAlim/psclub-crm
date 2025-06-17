package middleware

import (
	"bytes"
	"io"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger logs HTTP requests with method, path, status, latency and body.
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		status := c.Writer.Status()
		latency := time.Since(start)
		log.Printf("%s %s%s | %d | %v | %s", method, path, func() string {
			if query == "" {
				return ""
			}
			return "?" + query
		}(), status, latency, string(bodyBytes))
	}
}
