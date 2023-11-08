package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "ReCo web API is running",
	})
}
