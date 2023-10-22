package main

import (
	"net/http"
	"recofiit/service"

	"github.com/gin-gonic/gin"
)

func health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ReCo web API is running",
	})
}

func main() {
	router := gin.Default()
	service.Register()
	router.GET("/health", health)

	router.Run("localhost:8080")
}
