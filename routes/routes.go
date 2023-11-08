package routes

import (
	"recofiit/controllers/rest"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {

	r.GET("/health", rest.GetHealth)
	r.GET("/ws", rest.WsHandler)
}
