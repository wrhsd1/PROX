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
	// Set up the rate limiter middleware using an in-memory store
	limitStore = ratelimit.InMemoryStore(
		&ratelimit.InMemoryOptions{
			Rate:  time.Minute,
			Limit: 40,
		},
	)
	limitMiddleware = ratelimit.RateLimiter(
		limitStore,
		&ratelimit.Options{
			// Handler to use when the request rate is exceeded
			ErrorHandler: func(c *gin.Context, info ratelimit.Info) {
				c.JSON(
					http.StatusTooManyRequests,
					gin.H{
						"message": "Too many requests",
					},
				)
				c.Abort()
			},
			// Key function to use for rate limiting
			KeyFunc: func(c *gin.Context) string {
				return c.ClientIP()
			},
		},
	)
}

func secretAuth() gin.HandlerFunc {
	// Middleware function to authenticate secret header
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
	// Get the port to use from the PORT environment variable, or use the default port
	port := os.Getenv("PORT")
	if port == "" {
		port = DEFAULT_PORT
	}

	// Set up the Gin router with logging and middleware
	handler := gin.Default()

	// Add rate limiter middleware if the API is public
	if !api.Config.Private {
		handler.Use(limitMiddleware)
	}

	// Add secret authentication middleware
	handler.Use(secretAuth())

	// Proxy all requests to /* to proxy if not already handled
	handler.Any("/*path", handlers.Proxy)

	// Use endless package to listen on the specified port
	err := endless.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
