package ws_controller_namespace

import (
	"fmt"
	"recofiit/models"
	"recofiit/services/database"
	wsservice "recofiit/services/wsService"
	"strconv"
)

type WsControllerController struct{}

func (w WsControllerController) Get(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("GET CONTROLLER")
	type Body struct {
		ID uint `json:"id"`
	}

	var Req wsservice.WsRequestPrepared[Body]
	Req.Parse(req)

	db := database.GetDB()
	var ctrl models.Controller
	db.Find(&ctrl, Req.Body.ID)

	return wsservice.WsResponse[interface{}]{
		Namespace: "controller",
		Endpoint:  "get",
		Body:      ctrl,
	}
}
func (w WsControllerController) List(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("LIST CONTROLLER")
	type Body struct {
		CarVin *string `json:"vin"`
	}

	var Req wsservice.WsRequestPrepared[Body]
	Req.Parse(req)

	db := database.GetDB()
	var ctrls []models.Controller
	if Req.Body.CarVin == nil {
		db.Find(&ctrls)
	} else {
		var cars []models.Car
		db.Where("vin = ?", Req.Body.CarVin).Find(&cars)

		db.Where("car_id = ?", cars[0].ID).Find(&ctrls)
	}

	return wsservice.WsResponse[interface{}]{
		Namespace: "controller",
		Endpoint:  "list",
		Body:      ctrls,
	}
}
func (w WsControllerController) Create(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("CREATE CONTROLLER")
	type Body struct {
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		FirmwareID  uint   `json:"firmware_id"`
	}
	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var ctrl models.Controller
	ctrl.Name = Req.Body.Name
	ctrl.Type = Req.Body.Type
	ctrl.Description = Req.Body.Description
	ctrl.ControllerInstances = []models.ControllerInstace{}
	db.Create(&ctrl)

	var ci models.ControllerInstace
	ci.ControllerID = ctrl.ID
	ci.FirmwareID = Req.Body.FirmwareID
	db.Create(&ci)

	return wsservice.WsResponse[interface{}]{
		Namespace: "controller",
		Endpoint:  "create",
		Body:      ctrl,
	}
}
func (w WsControllerController) Update(req []byte) wsservice.WsResponse[interface{}] {
	return wsservice.WsResponse[interface{}]{} // TODO add update
}
func (w WsControllerController) Delete(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("DELETE CONTROLLER")
	type Body struct {
		ID uint `json:"id"`
	}
	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var ctrl models.Controller
	db.Find(&ctrl, Req.Body.ID)

	db.Delete(&ctrl)

	fmt.Println("DELETED CONTROLLER " + strconv.Itoa(int(ctrl.ID)))

	return wsservice.WsResponse[interface{}]{
		Namespace: "controller",
		Endpoint:  "delete",
		Body:      ctrl,
	}
}
