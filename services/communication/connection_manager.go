package communication

import (
	"fmt"
	_ "net/http/pprof"
	"time"

	api "github.com/ReCoFIIT/integration-api"
	"github.com/getsentry/sentry-go"
	"github.com/jftuga/geodist"
	_ "github.com/joho/godotenv/autoload"
)

// Keeps the connection with Integration module alive and ensures there is a connection
func sendKeepAlives(connection *IntegrationModuleConnection, interval float32) {
	for {
		var datagram api.IDatagram

		// Connect
		datagram = &api.ConnectDatagram{
			BaseDatagram: api.BaseDatagram{Type: "connect"},
		}

		acknowledged := connection.WriteAcknowledgedDatagram(datagram, 3, true)
		if !acknowledged {
			fmt.Printf("Could not connect to %v\n", connection.ServerAddress)
			continue
		}

		// Start sending keep-alives
		for range time.Tick(time.Second * time.Duration(interval)) {
			// Keep Alive
			datagram = &api.KeepAliveDatagram{
				BaseDatagram: api.BaseDatagram{Type: "keepalive"},
			}

			acknowledged = connection.WriteAcknowledgedDatagram(datagram, 3, true)
			if !acknowledged {
				fmt.Printf("Could not send keep-alive to %v\n", connection.ServerAddress)
				break
			}
		}
	}
}

func maintainSubscription(connection *IntegrationModuleConnection, subscriptionContent string, subscriptionTopic string, subscriptionInterval float32, checkInterval float32) {
	if connection == nil {
		panic("Connection cannot be nil")
	}

	var datagram api.IDatagram

	datagram = &api.SubscribeDatagram{
		BaseDatagram: api.BaseDatagram{
			Type: "subscribe",
		},
		Content:  subscriptionContent,
		Topic:    subscriptionTopic,
		Interval: subscriptionInterval,
	}

	for {
		// Subscribe
		acknowledged := connection.WriteAcknowledgedDatagram(datagram, 3, true)
		if !acknowledged {
			fmt.Printf("Could not subscribe to %v\n", connection.ServerAddress)
		}

		// Time to get first updates
		time.Sleep(time.Second * time.Duration(subscriptionInterval*2))

		// Check subscription is active
		for range time.Tick(time.Second * time.Duration(checkInterval)) {

			var lastDatagram api.IDatagram
			switch subscriptionContent {
			case "vehicles":
				if connection.LastOnUpdateVehiclesDatagram != nil {
					lastDatagram = connection.LastOnUpdateVehiclesDatagram
				}
			case "notifications":
				if connection.LastOnUpdateNotificationsDatagram != nil {
					lastDatagram = connection.LastOnUpdateNotificationsDatagram
				}
			}

			if lastDatagram == nil {
				break
			}

			lastUpdateTime, err := time.Parse(api.TimestampFormat, lastDatagram.GetTimestamp())
			if err != nil {
				sentry.CaptureException(err)
				break
			}
			secondsSinceLastUpdate := float32(time.Now().Sub(lastUpdateTime).Seconds())
			if secondsSinceLastUpdate > subscriptionInterval*1.1 {
				break
			}
		}
	}
}

// notificationsType can be "head_collision" or "chain_collision" or "crossroad"
func simulateNotifications(connection *IntegrationModuleConnection, notificationsType string) {
	for {
		vehiclesDatagram := <-connection.OnUpdateVehicles
		vehicles := vehiclesDatagram.Vehicles
		for i := 0; i < len(vehicles)-1; i++ {
			for j := i + 1; j < len(vehicles); j++ {
				vehicleA := vehicles[i]
				vehicleB := vehicles[j]

				_, km := geodist.HaversineDistance(
					geodist.Coord{Lat: float64(vehicleA.Latitude), Lon: float64(vehicleA.Longitude)},
					geodist.Coord{Lat: float64(vehicleB.Latitude), Lon: float64(vehicleB.Longitude)})

				metersDistance := km * 1000.0

				const SafeDistance = 100

				if notificationsType == "head_collision" {
					if metersDistance < SafeDistance {

						level := "info"
						if metersDistance < SafeDistance/2.0 {
							level = "warning"
						}
						if metersDistance < SafeDistance/4.0 {
							level = "danger"
						}

						datagram := &api.HeadCollisionNotifyDatagram{
							NotifyDatagram: api.NotifyDatagram{
								BaseDatagram: api.BaseDatagram{
									Type: "notify",
								},
								VehicleVin:  vehicleA.Vin,
								Level:       level,
								ContentType: "head_collision",
							},
							Content: api.HeadCollisionNotificationContent{
								TargetVehicleVin:     vehicleB.Vin,
								TimeToCollision:      float32(metersDistance / 100.0),
								MaxSpeedExceededBy:   10,
								BreakingDistanceDiff: float32(metersDistance),
							},
						}
						connection.WriteAcknowledgedDatagram(datagram, 2, true)

						datagram = &api.HeadCollisionNotifyDatagram{
							NotifyDatagram: api.NotifyDatagram{
								BaseDatagram: api.BaseDatagram{
									Type: "notify",
								},
								VehicleVin:  vehicleB.Vin,
								Level:       level,
								ContentType: "head_collision",
							},
							Content: api.HeadCollisionNotificationContent{
								TargetVehicleVin:     vehicleA.Vin,
								TimeToCollision:      float32(metersDistance / 100.0),
								MaxSpeedExceededBy:   10,
								BreakingDistanceDiff: float32(metersDistance),
							},
						}
						connection.WriteAcknowledgedDatagram(datagram, 2, true)
					}

				} else if notificationsType == "chain_collision" {
					if metersDistance < SafeDistance {

						level := "info"
						if metersDistance < SafeDistance/2.0 {
							level = "warning"
						}
						if metersDistance < SafeDistance/4.0 {
							level = "danger"
						}

						datagram := &api.ChainCollisionNotifyDatagram{
							NotifyDatagram: api.NotifyDatagram{
								BaseDatagram: api.BaseDatagram{
									Type: "notify",
								},
								VehicleVin:  vehicleA.Vin,
								Level:       level,
								ContentType: "chain_collision",
							},
							Content: api.ChainCollisionNotificationContent{
								TargetVehicleVin:    vehicleB.Vin,
								CurrentDistance:     float32(metersDistance),
								RecommendedDistance: SafeDistance,
							},
						}
						connection.WriteAcknowledgedDatagram(datagram, 2, true)

						datagram = &api.ChainCollisionNotifyDatagram{
							NotifyDatagram: api.NotifyDatagram{
								BaseDatagram: api.BaseDatagram{
									Type: "notify",
								},
								VehicleVin:  vehicleB.Vin,
								Level:       level,
								ContentType: "chain_collision",
							},
							Content: api.ChainCollisionNotificationContent{
								TargetVehicleVin:    vehicleA.Vin,
								CurrentDistance:     float32(metersDistance),
								RecommendedDistance: SafeDistance,
							},
						}
						connection.WriteAcknowledgedDatagram(datagram, 2, true)
					}
				} else if notificationsType == "crossroad" {
					if metersDistance < SafeDistance {

						level := "info"
						if metersDistance < SafeDistance/2.0 {
							level = "warning"
						}
						if metersDistance < SafeDistance/4.0 {
							level = "danger"
						}

						datagram := &api.CrossroadNotifyDatagram{
							NotifyDatagram: api.NotifyDatagram{
								BaseDatagram: api.BaseDatagram{
									Type: "notify",
								},
								VehicleVin:  vehicleA.Vin,
								Level:       level,
								ContentType: "crossroad",
							},
							Content: api.CrossroadNotificationContent{
								Text:       "Pojdeš prvý.",
								Order:      1,
								RightOfWay: true,
							},
						}
						connection.WriteAcknowledgedDatagram(datagram, 2, true)

						datagram = &api.CrossroadNotifyDatagram{
							NotifyDatagram: api.NotifyDatagram{
								BaseDatagram: api.BaseDatagram{
									Type: "notify",
								},
								VehicleVin:  vehicleB.Vin,
								Level:       level,
								ContentType: "crossroad",
							},
							Content: api.CrossroadNotificationContent{
								Text:       "Pojdeš druhý.",
								Order:      2,
								RightOfWay: false,
							},
						}
						connection.WriteAcknowledgedDatagram(datagram, 2, true)
					}
				} else if notificationsType == "generic" {
					if metersDistance < SafeDistance {

						level := "info"
						if metersDistance < SafeDistance/2.0 {
							level = "warning"
						}
						if metersDistance < SafeDistance/4.0 {
							level = "danger"
						}

						datagram := &api.GenericNotifyDatagram{
							NotifyDatagram: api.NotifyDatagram{
								BaseDatagram: api.BaseDatagram{
									Type: "notify",
								},
								VehicleVin:  vehicleA.Vin,
								Level:       level,
								ContentType: "generic",
							},
							Content: api.GenericNotificationContent{
								Text: "Prajeme príjemnú jazdu.",
							},
						}
						connection.WriteAcknowledgedDatagram(datagram, 2, true)

						datagram = &api.GenericNotifyDatagram{
							NotifyDatagram: api.NotifyDatagram{
								BaseDatagram: api.BaseDatagram{
									Type: "notify",
								},
								VehicleVin:  vehicleB.Vin,
								Level:       level,
								ContentType: "generic",
							},
							Content: api.GenericNotificationContent{
								Text: "Prajeme príjemnú jazdu.",
							},
						}
						connection.WriteAcknowledgedDatagram(datagram, 2, true)
					}
				}
			}
		}
	}
}
