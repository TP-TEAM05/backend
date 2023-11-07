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

	go wsservice.Manager.Start()

	routes.Setup(router)

	router.Run("localhost:8080")
}
