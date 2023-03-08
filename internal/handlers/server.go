package handlers

import (
	"crypto/tls"
	_ "embed"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/RunawayVPN/PROX/temp2/internal/types"

	"github.com/gin-gonic/gin"
)

var (
	//go:embed config.json
	config_file []byte
	Config      types.Config
)

// config returns the config.json file as a Config struct.
func init() {
	Config = types.Config{}
	if json.Unmarshal(config_file, &Config) != nil {
		log.Fatal("Error unmarshalling config.json")
	}
}

func Proxy(c *gin.Context) {
	// Check if Authorization header is set
	if c.Request.Header.Get("Authorization") == "" {
		c.JSON(401, gin.H{"message": "Unauthorized"})
		return
	}
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	// Add CORS headers
	c.Header("Access-Control-Allow-Origin", "images.duti.tech")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
	c.Header("Access-Control-Allow-Credentials", "true")
	// Proxy all requests directly to endpoint
	url := Config.Endpoint + "/api" + c.Param("path")
	// POST request with all data and headers
	var req *http.Request
	var err error
	if c.Request.Method == "POST" {
		req, err = http.NewRequest("POST", url, c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"message": "Internal server error"})
			return
		}
	} else if c.Request.Method == "GET" {
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			c.JSON(500, gin.H{"message": "Internal server error"})
			return
		}
	} else if c.Request.Method == "PATCH" {
		req, err = http.NewRequest("PATCH", url, c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"message": "Internal server error"})
			return
		}
	} else if c.Request.Method == "DELETE" {
		req, err = http.NewRequest("DELETE", url, c.Request.Body)
		if err != nil {
			c.JSON(500, gin.H{"message": "Internal server error"})
			return
		}
	} else if c.Request.Method == "OPTIONS" {
		c.JSON(200, gin.H{"message": "OK"})
		return
	} else {
		c.JSON(500, gin.H{"message": "Internal server error", "error": "Invalid HTTP method"})
		return
	}
	// Add headers
	for key, value := range c.Request.Header {
		req.Header.Set(key, value[0])
	}
	// Add content type JSON
	req.Header.Set("Content-Type", "application/json")
	// Set keep alive and timeout
	req.Close = false
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Keep-Alive", "timeout=360")
	// Send request
	client := &http.Client{Timeout: time.Second * 360}
	resp, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}
	// Stream response to client
	defer resp.Body.Close()
	// Set content type as text/event-stream
	c.Header("Content-Type", "text/event-stream")

	// Return stream of data to client
	c.Stream(func(w io.Writer) bool {
		// Write data to client
		io.Copy(w, resp.Body)
		return false
	})
}
