package communication

import (
	"encoding/json"
	"fmt"
	"net"
	ws_session_namespace "recofiit/controllers/ws/Session"
	"sync"
	"time"
)

type IConnection interface {
	Establish()
	ProcessDatagram(data []byte)
	WriteDatagram(datagram IDatagram)
	WriteAcknowledgedDatagram(datagram IDatagram, retries int) bool
}

type IntegrationModuleConnection struct {
	sync.Mutex
	UDPConn           *net.UDPConn
	ServerAddress     *net.UDPAddr
	NextSendIndex     int
	LastReceivedIndex int

	// To check whether the subscription works
	LastOnUpdateVehiclesDatagram      *UpdateVehiclesDatagram
	LastOnUpdateNotificationsDatagram *UpdateNotificationsDatagram

	// Received datagrams
	OnUpdateVehicles      chan UpdateVehiclesDatagram
	OnUpdateNotifications chan UpdateNotificationsDatagram
	OnArea                chan AreaDatagram

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
		OnUpdateVehicles:      make(chan UpdateVehiclesDatagram, 20),      // TODO constant
		OnUpdateNotifications: make(chan UpdateNotificationsDatagram, 20), // TODO constant
		OnArea:                make(chan AreaDatagram, 20),                // TODO constant
		Quit:                  make(chan bool),
	}
	var err error
	connection.UDPConn, err = net.DialUDP("udp", nil, connection.ServerAddress)

	if err != nil {
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

			fmt.Printf("Read a message %s ... \n", data[:min(len(data), 128)])
			if err != nil {
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

func (connection *IntegrationModuleConnection) WriteDatagram(datagram IDatagram, safe bool) {
	if safe {
		connection.Lock()
		defer connection.Unlock()
	}

	datagram.SetTimestamp(time.Now().UTC().Format(TimestampFormat))
	datagram.SetIndex(connection.NextSendIndex)

	data, err := json.Marshal(datagram)
	if err != nil {
		fmt.Printf("Error writing diagram %v\n", err)
	}
	_, _ = connection.UDPConn.Write(data)
	connection.NextSendIndex++
	fmt.Printf("Sending message to %v: %s\n", connection.ServerAddress, data[:min(len(data), 128)])
}

func (connection *IntegrationModuleConnection) ProcessDatagram(data []byte, safe bool) {
	if safe {
		connection.Lock()
		defer connection.Unlock()
	}

	// Parse data to JSON
	var datagram BaseDatagram
	err := json.Unmarshal(data, &datagram)
	if err != nil {
		fmt.Print("Parsing JSON failed.")
		return
	}
	// TODO uncomment this
	//if datagram.Index <= connection.LastReceivedIndex {
	//	return
	//}

	// TODO more cases
	switch datagram.Type {

	case "acknowledge":
		var acknowledgeDatagram AcknowledgeDatagram
		_ = json.Unmarshal(data, &acknowledgeDatagram)

		// Check if anyone is waiting for this
		channel, ok := connection.AcknowledgeWaiters[acknowledgeDatagram.AcknowledgingIndex]
		if ok {
			select {
			case channel <- true:
			default:
			}
		}

	case "update_vehicle_position":
		var updateVehiclesDatagram UpdateVehiclesDatagram
		_ = json.Unmarshal(data, &updateVehiclesDatagram)
		connection.LastOnUpdateVehiclesDatagram = &updateVehiclesDatagram

		// Send processed data to FE
		controller := ws_session_namespace.WsSessionController{}
		// time sleep 10 ms
		time.Sleep(10 * time.Millisecond)
		controller.SendLiveSessionData(&updateVehiclesDatagram)

		select {
		case connection.OnUpdateVehicles <- updateVehiclesDatagram:
		default:
			<-connection.OnUpdateVehicles // Discard oldest
			connection.OnUpdateVehicles <- updateVehiclesDatagram
		}

	case "update_notifications":
		var updateNotificationsDatagram UpdateNotificationsDatagram
		_ = json.Unmarshal(data, &updateNotificationsDatagram)
		connection.LastOnUpdateNotificationsDatagram = &updateNotificationsDatagram
		select {
		case connection.OnUpdateNotifications <- updateNotificationsDatagram:
		default:
			<-connection.OnUpdateNotifications // Discard oldest
			connection.OnUpdateNotifications <- updateNotificationsDatagram
		}

	case "area":
		var areaDatagram AreaDatagram
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

func (connection *IntegrationModuleConnection) WriteAcknowledgedDatagram(datagram IDatagram, retries int, safe bool) bool {
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
		case <-time.After(2 * time.Second): // TODO constant
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
