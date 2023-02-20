package main

import (
	"os"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/acheong08/ChatGPT-V2/internal/handlers"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

var limit_middleware gin.HandlerFunc
var limit_store ratelimit.Store

func init() {
	limit_store = ratelimit.InMemoryStore(
		&ratelimit.InMemoryOptions{
			Rate:  time.Minute,
			Limit: 40,
		},
	)
	limit_middleware = ratelimit.RateLimiter(
		limit_store,
		&ratelimit.Options{
			ErrorHandler: func(c *gin.Context, info ratelimit.Info) {
				c.JSON(
					429,
					gin.H{
						"message": "Too many requests",
					},
				)
				c.Abort()
			},
			KeyFunc: func(c *gin.Context) string {
				// Get Authorization header
				return c.ClientIP()
			},
		},
	)
}

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	handler := gin.Default()
	handler.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})
	if !handlers.Config.Private {
		handler.Any("/api/*path", limit_middleware, handlers.Proxy)
	} else {
		handler.Any("/api/*path", handlers.Proxy)
	}

	endless.ListenAndServe(":"+PORT, handler)
}
