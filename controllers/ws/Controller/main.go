package ws_controller_namespace

import (
	"gorm.io/gorm"
	"recofiit/models"
	"recofiit/services/database"
	wsservice "recofiit/services/wsService"
)

type WsControllerController struct{}

type ExtendedController struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	FirmwareID  uint    `json:"firmware_id"`
	Vin         *string `json:"vin"`
}

func (w WsControllerController) Get(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		ID uint `json:"id"`
	}

	var Req wsservice.WsRequestPrepared[Body]
	Req.Parse(req)

	db := database.GetDB()
	var ctrl models.Controller
	db.First(&ctrl, Req.Body.ID)

	var ec = w.ExtendController(ctrl, db)

	return wsservice.WsResponse[interface{}]{
		Namespace: "controller",
		Endpoint:  "get",
		Body:      ec,
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

	ctrls_extended := make([]ExtendedController, 0, len(ctrls))

	for _, c := range ctrls {
		var ec = w.ExtendController(c, db)

		ctrls_extended = append(ctrls_extended, ec)
	}

	return wsservice.WsResponse[interface{}]{
		Namespace: "controller",
		Endpoint:  "list",
		Body:      ctrls_extended,
	}
}
func (w WsControllerController) Create(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		Name                string `json:"name"`
		Type                string `json:"type"`
		Description         string `json:"description"`
		FirmwareID          uint   `json:"firmware_id"`
		FirmwareVersion     string `json:"firmware_version"`
		FirmwareDescription string `json:"firmware_description"`
		Vin                 string `json:"vin"`
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

	if Req.Body.FirmwareID == 0 {
		var f models.Firmware
		f.Version = Req.Body.FirmwareVersion
		f.Description = Req.Body.FirmwareDescription

		db.Create(&f)

		ci.FirmwareID = f.ID
	} else {
		ci.FirmwareID = Req.Body.FirmwareID
	}

	db.Create(&ci)

	if Req.Body.Vin != "" {
		var car models.Car
		db.Where("vin = ?", Req.Body.Vin).First(&car)

		var cc models.CarController
		cc.CarID = car.ID
		cc.ControllerInstanceID = ci.ID

		db.Create(&cc)
	}

	var ec = w.ExtendController(ctrl, db)

	return wsservice.WsResponse[interface{}]{
		Namespace: "controller",
		Endpoint:  "create",
		Body:      ec,
	}
}
func (w WsControllerController) Update(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		ID                  uint   `json:"id"`
		Name                string `json:"name"`
		Type                string `json:"type"`
		Description         string `json:"description"`
		FirmwareID          uint   `json:"firmware_id"`
		FirmwareVersion     string `json:"firmware_version"`
		FirmwareDescription string `json:"firmware_description"`
		Vin                 string `json:"vin"`
	}
	var Req wsservice.WsRequestPrepared[Body]
	Req.Parse(req)

	db := database.GetDB()

	var ctrl models.Controller
	db.First(&ctrl, Req.Body.ID)

	var ci models.ControllerInstance
	db.Where("controller_id = ?", ctrl.ID).Where("deleted_at is null").First(&ci)

	ctrl.Name = Req.Body.Name
	ctrl.Type = Req.Body.Type
	ctrl.Description = Req.Body.Description

	db.Save(&ctrl)

	var fid = Req.Body.FirmwareID

	if Req.Body.FirmwareID == 0 {
		var f models.Firmware
		f.Version = Req.Body.FirmwareVersion
		f.Description = Req.Body.FirmwareDescription

		db.Create(&f)

		fid = f.ID
	}

	var ccs []models.CarController
	db.Where("controller_instance_id = ?", ci.ID).Where("deleted_at is null").Find(&ccs)

	db.Delete(&ccs)

	// IF the firmwareID is the same, change only the controller
	// OTHERWISE change also the controller instance and copy car_controllers and sensors
	if ci.FirmwareID != fid {
		var newCi = w.RefreshInstance(ctrl, ci, fid, 0, 0)

		ci = newCi
	}

	if Req.Body.Vin != "" {
		var car models.Car
		db.Where("vin = ?", Req.Body.Vin).First(&car)

		var cc models.CarController
		cc.CarID = car.ID
		cc.ControllerInstanceID = ci.ID

		db.Create(&cc)
	}

	var ec = w.ExtendController(ctrl, db)

	return wsservice.WsResponse[interface{}]{
		Namespace: "controller",
		Endpoint:  "update",
		Body:      ec,
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
	db.First(&ctrl, Req.Body.ID)

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

	db.Delete(&ci, ci.ID)

	return newCi
}

func (w WsControllerController) ExtendController(
	ctrl models.Controller, db *gorm.DB) ExtendedController {

	var ec ExtendedController
	ec.ID = ctrl.ID
	ec.Name = ctrl.Name
	ec.Type = ctrl.Type
	ec.Description = ctrl.Description

	var ci models.ControllerInstance
	db.Where("controller_id = ?", ctrl.ID).Where("deleted_at is null").First(&ci)

	ec.FirmwareID = ci.FirmwareID

	var carControllers []models.CarController
	db.Where("controller_instance_id = ?", ci.ID).Where("deleted_at is null").Find(&carControllers)

	if len(carControllers) > 0 {
		var car models.Car
		db.First(&car, carControllers[0].CarID)

		ec.Vin = &car.Vin
	}

	return ec
}
