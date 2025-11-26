package communication

import (
	"encoding/json"
	"fmt"
	"net"
	ws_session_namespace "recofiit/controllers/ws/Session"
	"recofiit/services/dataLogging"
	"sync"
	"time"

	api "github.com/TP-TEAM05/integration-api"
	"github.com/getsentry/sentry-go"
)

type IConnection interface {
	Establish()
	ProcessDatagram(data []byte)
	WriteDatagram(datagram api.IDatagram)
	WriteAcknowledgedDatagram(datagram api.IDatagram, retries int) bool
}

type IntegrationModuleConnection struct {
	sync.Mutex
	UDPConn           *net.UDPConn
	ServerAddress     *net.UDPAddr
	NextSendIndex     int
	LastReceivedIndex int

	// To check whether the subscription works
	LastOnUpdateVehiclesDatagram      *api.UpdateVehiclesDatagram
	LastOnUpdateNotificationsDatagram *api.UpdateNotificationsDatagram

	// Received datagrams
	OnUpdateVehicles      chan api.UpdateVehiclesDatagram
	OnUpdateNotifications chan api.UpdateNotificationsDatagram
	OnArea                chan api.AreaDatagram

	// Maps indexes of sent diagrams to channels waiting for their acknowledgement
	AcknowledgeWaiters map[int]chan bool

	Quit chan bool
}

func NewIntegrationModuleConnection(addr *net.UDPAddr) *IntegrationModuleConnection {
	connection := &IntegrationModuleConnection{
		ServerAddress:         addr,
		NextSendIndex:         0,
		LastReceivedIndex:     -1,
		AcknowledgeWaiters:    make(map[int]chan bool),
		OnUpdateVehicles:      make(chan api.UpdateVehiclesDatagram, 20),      // TODO constant
		OnUpdateNotifications: make(chan api.UpdateNotificationsDatagram, 20), // TODO constant
		OnArea:                make(chan api.AreaDatagram, 20),                // TODO constant
		Quit:                  make(chan bool),
	}
	var err error
	connection.UDPConn, err = net.DialUDP("udp", nil, connection.ServerAddress)

	if err != nil {
		sentry.CaptureException(err)
		fmt.Printf("Error initializing UDP connection: %v\n", err)
		return nil
	}

	return connection
}

func (connection *IntegrationModuleConnection) Establish() {
	go func() {
		readBuffer := make([]byte, 65536)

		for {
			readBufferLength, err := connection.UDPConn.Read(readBuffer)
			data := readBuffer[:readBufferLength]

			if err != nil {
				sentry.CaptureException(err)
				fmt.Printf("Error reading message %v\n", err)
				continue
			}

			connection.ProcessDatagram(data, true)

			select {
			case <-connection.Quit:
				break
			default:
			}
		}
	}()
}

func (connection *IntegrationModuleConnection) WriteDatagram(datagram api.IDatagram, safe bool) {
	if safe {
		connection.Lock()
		defer connection.Unlock()
	}

	datagram.SetTimestamp(time.Now().UTC().Format(api.TimestampFormat))
	datagram.SetIndex(connection.NextSendIndex)

	data, err := json.Marshal(datagram)
	if err != nil {
		sentry.CaptureException(err)
		fmt.Printf("Error writing diagram %v\n", err)
	}
	_, _ = connection.UDPConn.Write(data)
	connection.NextSendIndex++
}

func (connection *IntegrationModuleConnection) ProcessDatagram(data []byte, safe bool) {
	if safe {
		connection.Lock()
		defer connection.Unlock()
	}

	// Parse data to JSON
	var datagram api.BaseDatagram
	err := json.Unmarshal(data, &datagram)
	if err != nil {
		sentry.CaptureException(err)
		fmt.Print("Parsing JSON failed: ", err)
		return
	}

	// DEBUG: Log all received datagram types
	if datagram.Type != "update_vehicle_position" {
		fmt.Printf("[BACKEND-RX-ALL] Received %s datagram (index: %d)\n", datagram.Type, datagram.Index)
	}
	// TODO uncomment this
	//if datagram.Index <= connection.LastReceivedIndex {
	//	return
	//}

	// TODO more cases
	switch datagram.Type {

	case "acknowledge":
		var acknowledgeDatagram api.AcknowledgeDatagram
		_ = json.Unmarshal(data, &acknowledgeDatagram)

		fmt.Printf("[BACKEND-RX-ACK] Received ACK for index: %d\n", acknowledgeDatagram.AcknowledgingIndex)

		// Check if anyone is waiting for this
		channel, ok := connection.AcknowledgeWaiters[acknowledgeDatagram.AcknowledgingIndex]
		if ok {
			fmt.Printf("[BACKEND-ACK-MATCHED] Found waiter for index: %d\n", acknowledgeDatagram.AcknowledgingIndex)
			select {
			case channel <- true:
				fmt.Printf("[BACKEND-ACK-SENT] Notified waiter for index: %d\n", acknowledgeDatagram.AcknowledgingIndex)
			default:
				fmt.Printf("[BACKEND-ACK-FAILED] Could not notify waiter for index: %d\n", acknowledgeDatagram.AcknowledgingIndex)
			}
		} else {
			fmt.Printf("[BACKEND-ACK-NO-WAITER] No waiter found for index: %d\n", acknowledgeDatagram.AcknowledgingIndex)
		}

	case "update_vehicles":
		var updateVehiclesDatagram api.UpdateVehiclesDatagram
		_ = json.Unmarshal(data, &updateVehiclesDatagram)
		connection.LastOnUpdateVehiclesDatagram = &updateVehiclesDatagram

		// Send processed data to FE
		//controller := ws_session_namespace.WsSessionController{}
		// time sleep 10 ms
		//time.Sleep(10 * time.Millisecond)
		//controller.SendLiveSessionData(&updateVehiclesDatagram)

		select {
		case connection.OnUpdateVehicles <- updateVehiclesDatagram:
		default:
			<-connection.OnUpdateVehicles // Discard oldest
			connection.OnUpdateVehicles <- updateVehiclesDatagram
		}

	case "update_vehicle_position":
		// live updates served directly to a database
		var updateVehiclesDatagram api.UpdateVehicleDatagram
		_ = json.Unmarshal(data, &updateVehiclesDatagram)

		// DEBUG: Log received packet
		fmt.Printf("[BACKEND-RX] Received update_vehicle_position for VIN: %s (index: %d)\n",
			updateVehiclesDatagram.Vehicle.Vin, datagram.Index)

		dataLogging.LogData(updateVehiclesDatagram)
		// Send processed data to FE
		controller := ws_session_namespace.WsSessionController{}
		controller.SendLiveSessionData(&updateVehiclesDatagram)

	case "update_notifications":
		var updateNotificationsDatagram api.UpdateNotificationsDatagram
		_ = json.Unmarshal(data, &updateNotificationsDatagram)
		connection.LastOnUpdateNotificationsDatagram = &updateNotificationsDatagram
		select {
		case connection.OnUpdateNotifications <- updateNotificationsDatagram:
		default:
			<-connection.OnUpdateNotifications // Discard oldest
			connection.OnUpdateNotifications <- updateNotificationsDatagram
		}

	case "area":
		var areaDatagram api.AreaDatagram
		_ = json.Unmarshal(data, &areaDatagram)
		select {
		case connection.OnArea <- areaDatagram:
		default:
			<-connection.OnArea // Discard oldest
			connection.OnArea <- areaDatagram
		}
	}

	connection.LastReceivedIndex = datagram.Index
}

func (connection *IntegrationModuleConnection) WriteAcknowledgedDatagram(datagram api.IDatagram, retries int, safe bool) bool {
	remainingTries := retries
	for remainingTries > 0 {
		remainingTries--

		// Send datagram
		connection.WriteDatagram(datagram, safe)

		acknowledgeWaiterChannel := make(chan bool)

		if safe {
			connection.Lock()
		}
		connection.AcknowledgeWaiters[datagram.GetIndex()] = acknowledgeWaiterChannel
		if safe {
			connection.Unlock()
		}

		acknowledged := false

		select {
		case acknowledged = <-acknowledgeWaiterChannel:
		case <-time.After(5 * time.Second): // Increased from 2s to 5s due to network latency
		}

		// Cleanup waiter channel
		if safe {
			connection.Lock()
		}
		delete(connection.AcknowledgeWaiters, datagram.GetIndex())
		if safe {
			connection.Unlock()
		}

		return acknowledged
	}

	return false
}
