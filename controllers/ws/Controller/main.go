package ws_controller_namespace

import (
	"recofiit/models"
	"recofiit/services/database"
	wsservice "recofiit/services/wsService"
)

type WsControllerController struct{}

func (w WsControllerController) Get(req []byte) wsservice.WsResponse[interface{}] {
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
		var car models.Car
		db.Where("vin = ?", *Req.Body.CarVin).First(&car)

		var carControllers []models.CarController
		db.Where("car_id = ?", car.ID).Find(&carControllers)

		ctrlInstanceIds := make([]uint, 0, len(carControllers))
		for _, cc := range carControllers {
			ctrlInstanceIds = append(ctrlInstanceIds, cc.ControllerInstanceID)
		}

		var controllerInstances []models.ControllerInstance
		db.Find(&controllerInstances, ctrlInstanceIds)

		ctrlIds := make([]uint, 0, len(controllerInstances))
		for _, ci := range controllerInstances {
			ctrlIds = append(ctrlIds, ci.ControllerID)
		}

		db.Find(&ctrls, ctrlIds)
	}

	return wsservice.WsResponse[interface{}]{
		Namespace: "controller",
		Endpoint:  "list",
		Body:      ctrls,
	}
}
func (w WsControllerController) Create(req []byte) wsservice.WsResponse[interface{}] {
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
	ctrl.ControllerInstances = []models.ControllerInstance{}
	db.Create(&ctrl)

	var ci models.ControllerInstance
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
	type Body struct {
		ID          uint   `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		FirmwareID  uint   `json:"firmware_id"`
	}
	var Req wsservice.WsRequestPrepared[Body]
	Req.Parse(req)

	db := database.GetDB()

	var ctrl models.Controller
	db.Find(&ctrl, Req.Body.ID)

	var ci models.ControllerInstance
	db.Where("controller_id = ?", ctrl.ID).Where("deleted_at is null").First(&ci)

	ctrl.Name = Req.Body.Name
	ctrl.Type = Req.Body.Type
	ctrl.Description = Req.Body.Description

	db.Save(&ctrl)

	// IF the firmwareID is the same, change only the controller
	// OTHERWISE change also the controller instance and copy car_controllers and sensors
	if ci.FirmwareID != Req.Body.FirmwareID {
		w.RefreshInstance(ctrl, ci, Req.Body.FirmwareID, 0, 0)

		db.Delete(&ci)
	}

	return wsservice.WsResponse[interface{}]{
		Namespace: "controller",
		Endpoint:  "update",
		Body:      ctrl,
	}
}
func (w WsControllerController) Delete(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		ID uint `json:"id"`
	}
	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var ctrl models.Controller
	db.Find(&ctrl, Req.Body.ID)

	db.Delete(&ctrl)

	var ci models.ControllerInstance
	db.Where("controller_id = ?", ctrl.ID).Where("deleted_at is null").First(&ci)

	var carControllers []models.CarController
	db.Where("controller_instance_id = ?", ci.ID).Find(&carControllers)

	db.Delete(&carControllers)

	var sensors []models.Sensor
	db.Where("controller_instance_id = ?", ci.ID).Find(&sensors)

	db.Delete(&sensors)

	db.Delete(&ci)

	return wsservice.WsResponse[interface{}]{
		Namespace: "controller",
		Endpoint:  "delete",
		Body:      ctrl,
	}
}

func (w WsControllerController) RefreshInstance(
	ctrl models.Controller,
	ci models.ControllerInstance,
	newFirmwareID uint,
	exceptSensorId uint,
	exceptCarControllerId uint,
) models.ControllerInstance {
	db := database.GetDB()

	var newCi models.ControllerInstance
	newCi.ControllerID = ctrl.ID
	if newFirmwareID == 0 {
		newCi.FirmwareID = ci.FirmwareID
	} else {
		newCi.FirmwareID = newFirmwareID
	}
	db.Create(&newCi)

	var carControllers []models.CarController
	db.Where("controller_instance_id = ?", ci.ID).Find(&carControllers)

	for _, cc := range carControllers {
		if cc.ID == exceptCarControllerId {
			continue
		}
		var newCc models.CarController
		newCc.CarID = cc.CarID
		newCc.ControllerInstanceID = newCi.ID
		db.Save(&newCc)

		db.Delete(&cc)
	}

	var sensors []models.Sensor
	db.Where("controller_instance_id = ?", ci.ID).Find(&sensors)

	for _, s := range sensors {
		if s.ID == exceptSensorId {
			continue
		}
		var newSensor models.Sensor
		newSensor.Name = s.Name
		newSensor.SensorType = s.SensorType
		newSensor.ControllerInstanceID = newCi.ID
		db.Save(&newSensor)

		db.Delete(&s)
	}
	return newCi
}
