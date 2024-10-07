package main

import (
	"recofiit/routes"
	"recofiit/services"
	wsservice "recofiit/services/wsService"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
)

func main() {
	defer sentry.Flush(2 * time.Second)

	router := gin.Default()
	router.SetTrustedProxies([]string{""})
	services.Register()

	routes.Setup(router)
	routes.SetupWs(&wsservice.Manager)

	go wsservice.Manager.Start()

	router.Run(":8080")
}
