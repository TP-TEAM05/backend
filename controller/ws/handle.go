package ws

import (
	"net/http"
	wsservice "recofiit/services/wsService"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 1024 * 1024,
	WriteBufferSize: 1024 * 1024 * 1024,
	//Solving cross-domain problems
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(c *gin.Context) {
	//Upgrade the HTTP protocol to the websocket protocol
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	//Every connection will open a new client, client.id generates through UUID to ensure that each time it is different

	client := &wsservice.Client{ID: uuid.Must(uuid.NewV4(), nil).String(), Socket: conn, Send: make(chan []byte)}
	//Register a new link
	wsservice.Manager.Register <- client

	//Start the message to collect the news from the web side
	go client.Read()
	//Start the corporation to return the message to the web side
	go client.Write()
}
