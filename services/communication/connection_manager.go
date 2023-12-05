package communication

import (
	"fmt"
	_ "net/http/pprof"
	"recofiit/models"
	"time"

	"github.com/jftuga/geodist"
	_ "github.com/joho/godotenv/autoload"
)

// Keeps the connection with Integration module alive and ensures there is a connection
func sendKeepAlives(connection *IntegrationModuleConnection, interval float32) {
	for {
		var datagram models.IDatagram

		// Connect
		datagram = &models.ConnectDatagram{
			BaseDatagram: models.BaseDatagram{Type: "connect"},
		}

		acknowledged := connection.WriteAcknowledgedDatagram(datagram, 3, true)
		if !acknowledged {
			fmt.Printf("Could not connect to %v\n", connection.ServerAddress)
			continue
		}

		// Start sending keep-alives
		for range time.Tick(time.Second * time.Duration(interval)) {
			// Keep Alive
			datagram = &models.KeepAliveDatagram{
				BaseDatagram: models.BaseDatagram{Type: "keepalive"},
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

	var datagram models.IDatagram

	datagram = &models.SubscribeDatagram{
		BaseDatagram: models.BaseDatagram{
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

			var lastDatagram models.IDatagram
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

			lastUpdateTime, err := time.Parse(models.TimestampFormat, lastDatagram.GetTimestamp())
			if err != nil {
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
					geodist.Coord{Lat: float64(vehicleA.Position.Lat), Lon: float64(vehicleA.Position.Lon)},
					geodist.Coord{Lat: float64(vehicleB.Position.Lat), Lon: float64(vehicleB.Position.Lon)})

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

						datagram := &models.HeadCollisionNotifyDatagram{
							NotifyDatagram: models.NotifyDatagram{
								BaseDatagram: models.BaseDatagram{
									Type: "notify",
								},
								VehicleId:   vehicleA.Id,
								Level:       level,
								ContentType: "head_collision",
							},
							Content: models.HeadCollisionNotificationContent{
								TargetVehicleId:      vehicleB.Id,
								TimeToCollision:      float32(metersDistance / 100.0),
								MaxSpeedExceededBy:   10,
								BreakingDistanceDiff: float32(metersDistance),
							},
						}
						connection.WriteAcknowledgedDatagram(datagram, 2, true)

						datagram = &models.HeadCollisionNotifyDatagram{
							NotifyDatagram: models.NotifyDatagram{
								BaseDatagram: models.BaseDatagram{
									Type: "notify",
								},
								VehicleId:   vehicleB.Id,
								Level:       level,
								ContentType: "head_collision",
							},
							Content: models.HeadCollisionNotificationContent{
								TargetVehicleId:      vehicleA.Id,
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

						datagram := &models.ChainCollisionNotifyDatagram{
							NotifyDatagram: models.NotifyDatagram{
								BaseDatagram: models.BaseDatagram{
									Type: "notify",
								},
								VehicleId:   vehicleA.Id,
								Level:       level,
								ContentType: "chain_collision",
							},
							Content: models.ChainCollisionNotificationContent{
								TargetVehicleId:     vehicleB.Id,
								CurrentDistance:     float32(metersDistance),
								RecommendedDistance: SafeDistance,
							},
						}
						connection.WriteAcknowledgedDatagram(datagram, 2, true)

						datagram = &models.ChainCollisionNotifyDatagram{
							NotifyDatagram: models.NotifyDatagram{
								BaseDatagram: models.BaseDatagram{
									Type: "notify",
								},
								VehicleId:   vehicleB.Id,
								Level:       level,
								ContentType: "chain_collision",
							},
							Content: models.ChainCollisionNotificationContent{
								TargetVehicleId:     vehicleA.Id,
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

						datagram := &models.CrossroadNotifyDatagram{
							NotifyDatagram: models.NotifyDatagram{
								BaseDatagram: models.BaseDatagram{
									Type: "notify",
								},
								VehicleId:   vehicleA.Id,
								Level:       level,
								ContentType: "crossroad",
							},
							Content: models.CrossroadNotificationContent{
								Text:       "Pojdeš prvý.",
								Order:      1,
								RightOfWay: true,
							},
						}
						connection.WriteAcknowledgedDatagram(datagram, 2, true)

						datagram = &models.CrossroadNotifyDatagram{
							NotifyDatagram: models.NotifyDatagram{
								BaseDatagram: models.BaseDatagram{
									Type: "notify",
								},
								VehicleId:   vehicleB.Id,
								Level:       level,
								ContentType: "crossroad",
							},
							Content: models.CrossroadNotificationContent{
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

						datagram := &models.GenericNotifyDatagram{
							NotifyDatagram: models.NotifyDatagram{
								BaseDatagram: models.BaseDatagram{
									Type: "notify",
								},
								VehicleId:   vehicleA.Id,
								Level:       level,
								ContentType: "generic",
							},
							Content: models.GenericNotificationContent{
								Text: "Prajeme príjemnú jazdu.",
							},
						}
						connection.WriteAcknowledgedDatagram(datagram, 2, true)

						datagram = &models.GenericNotifyDatagram{
							NotifyDatagram: models.NotifyDatagram{
								BaseDatagram: models.BaseDatagram{
									Type: "notify",
								},
								VehicleId:   vehicleB.Id,
								Level:       level,
								ContentType: "generic",
							},
							Content: models.GenericNotificationContent{
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
