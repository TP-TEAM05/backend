package main

import (
	"recofiit/routes"
	"recofiit/services"
	wsservice "recofiit/services/wsService"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.SetTrustedProxies([]string{""})
	services.Register()

	routes.Setup(router)
	routes.SetupWs(&wsservice.Manager)

	go wsservice.Manager.Start()

	router.Run("localhost:8080")
}
