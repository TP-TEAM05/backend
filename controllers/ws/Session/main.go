package ws_session_namespace

import (
	"fmt"
	"recofiit/models"
	"recofiit/services/database"
	wsservice "recofiit/services/wsService"
)

type WsSessionController struct{}

func (w WsSessionController) Get(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		ID string `json:"id"`
	}

	var Req wsservice.WsRequestPrepared[Body]
	Req.Parse(req)

	return wsservice.WsResponse[interface{}]{}
}
func (w WsSessionController) List(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("LIST SESSION")
	db := database.GetDB()
	var sessions []models.Session
	db.Find(&sessions)
	return wsservice.WsResponse[interface{}]{}
}
func (w WsSessionController) Create(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		Cars []string `json:"cars"`
		Name string   `json:"name"`
	}

	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	fmt.Println("CREATE SESSION", Req.Body.Cars)

	fmt.Println("CREATE SESSION")
	return wsservice.WsResponse[interface{}]{}

}
func (w WsSessionController) Update(req []byte) wsservice.WsResponse[interface{}] {
	return wsservice.WsResponse[interface{}]{}
}
func (w WsSessionController) Delete(req []byte) wsservice.WsResponse[interface{}] {
	return wsservice.WsResponse[interface{}]{}
}
