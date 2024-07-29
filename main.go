package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var keysMap map[string]int

func init() {
	keysMap = map[string]int{
		"key1": 1,
		"key2": 2,
		"key3": 3,
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		fmt.Println("API_KEY not set in environment variables")
		return
	}

	router := gin.Default()

	router.Use(apiKeyMiddleware(apiKey))

	router.GET("/getKey", func(c *gin.Context) {
		key := c.Query("key")
		if key == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "key parameter is required"})
			return
		}

		value, ok := keysMap[key]
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "Key not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"value": value})
	})

	err = router.Run(":" + port)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func apiKeyMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestAPIKey := c.GetHeader("X-Prime-API-Key")
		if requestAPIKey != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}
		c.Next()
	}
}
