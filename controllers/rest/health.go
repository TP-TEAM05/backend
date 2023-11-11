package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHealth(c *gin.Context) {
	fmt.Println("GetHealth")
	c.JSON(http.StatusOK, gin.H{
		"message": "ReCo web API is running",
	})
}
