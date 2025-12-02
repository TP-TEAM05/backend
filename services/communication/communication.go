package communication

import (
	"fmt"
	"net"
	"os"

	"github.com/getsentry/sentry-go"
)

// func getIPs() []net.IP {
// 	ips, err := net.LookupIP("car-integration")
// 	if err != nil || len(ips) == 0 {
// 		sentry.CaptureException(err)
// 		fmt.Printf("Could not resolve or find hostname %v\n", err)
// 		os.Exit(1)
// 	}
// 	return ips
// }

func getIPs() []net.IP {
	host := os.Getenv("CAR_INTEGRATION_HOST")
	port := os.Getenv("CAR_INTEGRATION_PORT")

	if host == "" {
		host = "127.0.0.1" // fallback for local simulation
	}
	if port == "" {
		port = "5050"
	}

	fmt.Printf("Connecting to car-integration at %s:%s\n", host, port)

	ip := net.ParseIP(host)
	if ip == nil {
		// try DNS if host looks like a hostname
		ips, err := net.LookupIP(host)
		if err != nil || len(ips) == 0 {
			sentry.CaptureException(err)
			fmt.Printf("Could not resolve hostname '%s': %v\n", host, err)
			os.Exit(1)
		}
		return ips
	}

	return []net.IP{ip}
}

func subscribe() {
	var ips = getIPs()
	port := 5050
	serverAddress := net.UDPAddr{Port: port, IP: ips[0]}

	connection := NewIntegrationModuleConnection(&serverAddress)
	connection.Establish()

	// Set the global connection for manual control
	SetCarIntegrationConnection(connection)

	// Send keep-alives and check whether subscription is active periodically
	go sendKeepAlives(connection, 10)

	// Handle vehicle subscriptions
	go maintainSubscription(connection, "live-updates", "vehicles", 1, 10)
}

func Init() {
	subscribe()
}
