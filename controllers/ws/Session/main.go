package ws_session_namespace

import (
	"errors"
	"fmt"
	"github.com/stretchr/objx"
	"gorm.io/gorm"
	"recofiit/models"
	"recofiit/services/database"
	wsservice "recofiit/services/wsService"
)

type WsSessionController struct{}

type SessionGetRequestBody struct {
	id int
}

func (w WsSessionController) Get(c *wsservice.Client, req *wsservice.WsRequest, res *wsservice.WsResponse) {
	id := objx.New(req.Body).Get("id").Int()
	fmt.Println(id)
	fmt.Println("GET SESSION", id)
	db := database.GetDB()
	var session models.Session
	result := db.First(&session, id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		res.Error = "Not found"
	} else {
		res.Body = session
	}
}
func (w WsSessionController) List(c *wsservice.Client, req *wsservice.WsRequest, res *wsservice.WsResponse) {
	fmt.Println("LIST SESSION")
	db := database.GetDB()
	var sessions []models.Session
	db.Find(&sessions)

	res.Body = sessions
}
func (w WsSessionController) Create(c *wsservice.Client, req *wsservice.WsRequest, res *wsservice.WsResponse) {
}
func (w WsSessionController) Update(c *wsservice.Client, req *wsservice.WsRequest, res *wsservice.WsResponse) {
}
func (w WsSessionController) Delete(c *wsservice.Client, req *wsservice.WsRequest, res *wsservice.WsResponse) {
}
