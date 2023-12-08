package rest

import (
	"net/http"
	database "recofiit/services/database"
	redis "recofiit/services/redis"

	"github.com/gin-gonic/gin"
)

func GetHealth(c *gin.Context) {
	redisStatus := redis.HealthCheck()
	dbStatus := database.HealthCheck()
	c.JSON(http.StatusOK, gin.H{
		"API":   true,
		"DB":    dbStatus,
		"Redis": redisStatus,
	})
}
