package main

import (
	"log"
	"net/http"
	"os"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/acheong08/ChatGPT-V2/internal/api"
	"github.com/acheong08/ChatGPT-V2/internal/handlers"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

const (
	DEFAULT_PORT = "8080"
)

var (
	limitStore      ratelimit.Store
	limitMiddleware gin.HandlerFunc
)

func init() {
	limitStore = ratelimit.InMemoryStore(
		&ratelimit.InMemoryOptions{
			Rate:  time.Minute,
			Limit: 40,
		},
	)
	limitMiddleware = ratelimit.RateLimiter(
		limitStore,
		&ratelimit.Options{
			ErrorHandler: func(c *gin.Context, info ratelimit.Info) {
				c.JSON(
					http.StatusTooManyRequests,
					gin.H{
						"message": "Too many requests",
					},
				)
				c.Abort()
			},
			KeyFunc: func(c *gin.Context) string {
				return c.ClientIP()
			},
		},
	)
}

func secretAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if secret := os.Getenv("SECRET"); secret != "" {
			authHeader := c.GetHeader("Secret")
			if authHeader != secret {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = DEFAULT_PORT
	}

	handler := gin.Default()

	if !api.Config.Private {
		handler.Use(limitMiddleware)
	}

	handler.Use(secretAuth())

	// Proxy all requests to /* to proxy if not already handled
	handler.Any("/*path", handlers.Proxy)

	err := endless.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
