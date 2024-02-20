package wsservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"recofiit/utils"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	//User ID
	ID string
	//Connected socket
	Socket *websocket.Conn
	//Message
	Send chan []byte
}

type ClientManager struct {
	//The client map stores and manages all long connection clients, online is TRUE, and those who are not there are FALSE
	Clients map[*Client]bool
	//Web side MESSAGE we use Broadcast to receive, and finally distribute it to all clients
	Broadcast chan []byte
	//Newly created long connection client
	Register chan *Client
	//Newly canceled long connection client
	Unregister chan *Client
	// WS Router
	Router WsRouter

	ActiveConnections int

	TotalConnections int
}

type Message struct {
	//Message Struct
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
	ServerIP  string `json:"serverIp,omitempty"`
	SenderIP  string `json:"senderIp,omitempty"`
}

var Manager = ClientManager{
	Broadcast:         make(chan []byte),
	Register:          make(chan *Client),
	Unregister:        make(chan *Client),
	Clients:           make(map[*Client]bool),
	Router:            WsRouter{Routes: make(Namespace)},
	ActiveConnections: 0,
	TotalConnections:  0,
}

func (manager *ClientManager) Start() {
	for {
		select {
		//If there is a new connection access, pass the connection to conn through the channel
		case conn := <-manager.Register:
			//Set the client connection to true
			manager.Clients[conn] = true
			manager.ActiveConnections++
			manager.TotalConnections++
			//Format the message of returning to the successful connection JSON
			jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected. ", ServerIP: utils.LocalIp(), SenderIP: conn.Socket.RemoteAddr().String()})
			//Call the client's send method and send messages
			manager.Send(jsonMessage, conn)
			//If the connection is disconnected
		case conn := <-manager.Unregister:
			manager.ActiveConnections--
			if _, ok := manager.Clients[conn]; ok {
				close(conn.Send)
				delete(manager.Clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected. ", ServerIP: utils.LocalIp(), SenderIP: conn.Socket.RemoteAddr().String()})
				manager.Send(jsonMessage, conn)
			}
			if (manager.TotalConnections - manager.ActiveConnections) > 10 {
				newClientsMap := make(map[*Client]bool)
				for conn := range manager.Clients {
					newClientsMap[conn] = true
				}
				manager.Clients = newClientsMap
				manager.TotalConnections = manager.ActiveConnections
			}
		//broadcast
		case message := <-manager.Broadcast:
			//Traversing the client that has been connected, send the message to them
			for conn := range manager.Clients {
				select {
				case conn.Send <- message:
				default:
					close(conn.Send)
					delete(manager.Clients, conn)
				}
			}
		}
	}
}

func (manager *ClientManager) Send(message []byte, ignore *Client) {
	for conn := range manager.Clients {
		//Send messages not to the shielded connection
		if conn != ignore {
			conn.Send <- message
		}
	}
}

// Define the read method of the client structure
func (c *Client) Read() {
	defer func() {
		Manager.Unregister <- c
		c.Socket.Close()
	}()
	c.Socket.SetReadLimit(maxMessageSize)
	c.Socket.SetReadDeadline(time.Now().Add(pongWait))
	c.Socket.SetPongHandler(func(string) error { c.Socket.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		Manager.Broadcast <- message

		ReqJson := &WsRequest[interface{}]{}
		ReqJson.Parse(message)

		Handle := Manager.Router.GetHandler(ReqJson.Namespace, ReqJson.Endpoint)

		if Handle == nil {
			fmt.Println("No handler found for this endpoint")
			return
		}

		res := Handle(message)

		c.Send <- res.ToJSON()
	}
}

func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()

	for {
		select {
		//Read the message from send
		case message, ok := <-c.Send:
			c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Socket.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		}
	}
}
