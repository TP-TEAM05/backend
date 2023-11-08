package rest

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

	client := &wsservice.Client{ID: uuid.Must(uuid.NewV4(), nil).String(), Socket: conn, Send: make(chan []byte)}

	wsservice.Manager.Register <- client

	// Start read and write at selected namespace
	go client.Read()
	go client.Write()
}
