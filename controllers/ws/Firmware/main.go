package ws_firmware_namespace

import (
	"fmt"
	"recofiit/models"
	"recofiit/services/database"
	wsservice "recofiit/services/wsService"
)

type WsFirmwareController struct{}

func (w WsFirmwareController) Get(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("GET FIRMWARE")
	type Body struct {
		ID uint `json:"id"`
	}

	var Req wsservice.WsRequestPrepared[Body]
	Req.Parse(req)

	db := database.GetDB()
	var fw models.Firmware
	db.Find(&fw, Req.Body.ID)

	return wsservice.WsResponse[interface{}]{
		Namespace: "firmware",
		Endpoint:  "get",
		Body:      fw,
	}
}
func (w WsFirmwareController) List(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("LIST FIRMWARE")
	db := database.GetDB()
	var fws []models.Firmware
	db.Find(&fws)

	return wsservice.WsResponse[interface{}]{
		Namespace: "firmware",
		Endpoint:  "list",
		Body:      fws,
	}
}
func (w WsFirmwareController) Create(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("CREATE FIRMWARE")
	type Body struct {
		Version     string `json:"version"`
		Description string `json:"description"`
	}
	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var fw models.Firmware
	fw.Version = Req.Body.Version
	fw.Description = Req.Body.Description
	db.Create(&fw)

	return wsservice.WsResponse[interface{}]{
		Namespace: "firmware",
		Endpoint:  "create",
		Body:      fw,
	}
}
func (w WsFirmwareController) Update(req []byte) wsservice.WsResponse[interface{}] {
	return wsservice.WsResponse[interface{}]{}
}
func (w WsFirmwareController) Delete(req []byte) wsservice.WsResponse[interface{}] {
	return wsservice.WsResponse[interface{}]{}
}
