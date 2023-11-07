package routes

import (
	"recofiit/controller/rest"
	"recofiit/controller/ws"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {

	r.GET("/health", rest.GetHealth)
	r.GET("/ws", ws.WsHandler)
}
