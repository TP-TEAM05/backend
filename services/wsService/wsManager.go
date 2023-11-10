package wsservice

import (
	"encoding/json"
	"fmt"
	"recofiit/utils"

	"github.com/gorilla/websocket"
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
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Clients:    make(map[*Client]bool),
	Router:     WsRouter{Routes: make(Namespace)},
}

func (manager *ClientManager) Start() {
	for {
		select {
		//If there is a new connection access, pass the connection to conn through the channel
		case conn := <-manager.Register:
			//Set the client connection to true
			manager.Clients[conn] = true
			//Format the message of returning to the successful connection JSON
			jsonMessage, _ := json.Marshal(&Message{Content: "/A new socket has connected. ", ServerIP: utils.LocalIp(), SenderIP: conn.Socket.RemoteAddr().String()})
			//Call the client's send method and send messages
			manager.Send(jsonMessage, conn)
			//If the connection is disconnected
		case conn := <-manager.Unregister:
			//Determine the state of the connection, if it is true, turn off Send and delete the value of connecting client
			if _, ok := manager.Clients[conn]; ok {
				close(conn.Send)
				delete(manager.Clients, conn)
				jsonMessage, _ := json.Marshal(&Message{Content: "/A socket has disconnected. ", ServerIP: utils.LocalIp(), SenderIP: conn.Socket.RemoteAddr().String()})
				manager.Send(jsonMessage, conn)
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
		_ = c.Socket.Close()
	}()

	for {
		//Read message
		_, message, err := c.Socket.ReadMessage()
		//If there is an error message, cancel this connection and then close it
		if err != nil {
			Manager.Unregister <- c
			_ = c.Socket.Close()
			break
		}

		messageJSON := WsRequest{}
		messageJSON.ParseJSON(message)

		handle := Manager.Router.GetHandler(messageJSON.Namespace, messageJSON.Endpoint)

		if handle == nil {
			fmt.Println("No handler found for this endpoint")
			return
		}

		handle()
	}
}

func (c *Client) Write() {
	defer func() {
		_ = c.Socket.Close()
	}()

	for {
		select {
		//Read the message from send
		case message, ok := <-c.Send:
			//If there is no message
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			//Write it if there is news and send it to the web side
			_ = c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
