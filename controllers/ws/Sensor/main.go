package ws_sensor_namespace

import (
	"fmt"
	"recofiit/models"
	"recofiit/services/database"
	wsservice "recofiit/services/wsService"
)

type WsSensorController struct{}

func (w WsSensorController) Get(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("GET SENSOR")
	type Body struct {
		ID uint `json:"id"`
	}

	var Req wsservice.WsRequestPrepared[Body]
	Req.Parse(req)

	db := database.GetDB()
	var sensor models.Sensor
	sensor.ID = Req.Body.ID
	db.Find(&sensor)

	return wsservice.WsResponse[interface{}]{
		Namespace: "sensor",
		Endpoint:  "get",
		Body:      sensor,
	}
}
func (w WsSensorController) List(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("LIST SENSOR")
	type Body struct {
		ControllerInstanceID uint `json:"controller_instance_id"`
	}

	var Req wsservice.WsRequestPrepared[Body]
	Req.Parse(req)

	db := database.GetDB()
	var sensors []models.Sensor
	db.Where("controller_instance_id = ?", Req.Body.ControllerInstanceID).Find(&sensors)

	return wsservice.WsResponse[interface{}]{
		Namespace: "sensor",
		Endpoint:  "list",
		Body:      sensors,
	}
}
func (w WsSensorController) Create(req []byte) wsservice.WsResponse[interface{}] {
	return wsservice.WsResponse[interface{}]{} // TODO create sensor with new controller instance
}
func (w WsSensorController) Update(req []byte) wsservice.WsResponse[interface{}] {
	return wsservice.WsResponse[interface{}]{}
}
func (w WsSensorController) Delete(req []byte) wsservice.WsResponse[interface{}] {
	return wsservice.WsResponse[interface{}]{}
}
