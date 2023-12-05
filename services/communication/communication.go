package communication

import (
	"fmt"
	"net"
	"os"
)

func getIPs() []net.IP {
	ips, err := net.LookupIP("car-integration")
	if err != nil || len(ips) == 0 {
		fmt.Printf("Could not resolve or find hostname %v\n", err)
		os.Exit(1)
	}
	return ips
}

func subscribe() {
	var ips = getIPs()
	port := 5050
	serverAddress := net.UDPAddr{Port: port, IP: ips[0]}

	connection := NewIntegrationModuleConnection(&serverAddress)
	connection.Establish()

	// Send keep-alives and check whether subscription is active periodically
	go sendKeepAlives(connection, 10)

	// Handle vehicle subscriptions
	go maintainSubscription(connection, "live-updates", "vehicles", 1, 10)
}

func Init() {
	subscribe()
}
