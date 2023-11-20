package main

import (
	"fmt"
	"net"
	"os"
	"recofiit/routes"
	"recofiit/services"
	wsservice "recofiit/services/wsService"

	"github.com/gin-gonic/gin"
)

func main() {
	// --------- Subscribe to Integration Module ---------------

	ips, err := net.LookupIP("car-integration")
	if err != nil || len(ips) == 0 {
		fmt.Printf("Could not resolve or find hostname %v\n", err)
		os.Exit(1)
	}

	port := 5050
	serverAddress := net.UDPAddr{Port: port, IP: ips[0]}

	connection := NewIntegrationModuleConnection(&serverAddress)
	connection.Establish()

	// Send keep-alives and check whether subscription is active periodically
	go sendKeepAlives(connection, 10)

	// Handle vehicle subscriptions
	go maintainSubscription(connection, "vehicles", 5, 10)

	// --------- Subscribe to Integration Module ---------------

	router := gin.Default()
	router.SetTrustedProxies([]string{""})
	services.Register()

	routes.Setup(router)
	routes.SetupWs(&wsservice.Manager)

	go wsservice.Manager.Start()

	router.Run(":8080")
}
