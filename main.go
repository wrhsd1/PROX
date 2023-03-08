package main

import (
	"os"

	"github.com/RunawayVPN/PROX/temp2/internal/handlers"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func main() {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080"
	}
	handler := gin.Default()
	handler.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	handler.Any("/api/*path", handlers.Proxy)

	endless.ListenAndServe(os.Getenv("HOST")+":"+PORT, handler)
}
