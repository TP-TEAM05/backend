package ws_sensor_namespace

import (
	ws_controller_namespace "recofiit/controllers/ws/Controller"
	"recofiit/models"
	"recofiit/services/database"
	wsservice "recofiit/services/wsService"
)

type WsSensorController struct{}

func (w WsSensorController) Get(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		ID uint `json:"id"`
	}

	var Req wsservice.WsRequestPrepared[Body]
	Req.Parse(req)

	db := database.GetDB()
	var sensor models.Sensor
	db.Find(&sensor, Req.Body.ID)

	return wsservice.WsResponse[interface{}]{
		Namespace: "sensor",
		Endpoint:  "get",
		Body:      sensor,
	}
}
func (w WsSensorController) List(req []byte) wsservice.WsResponse[interface{}] {
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
	type Body struct {
		ControllerID string            `json:"controller_id"`
		Name         string            `json:"name"`
		Type         models.SensorType `json:"type"`
	}

	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var ctrl models.Controller
	db.Find(&ctrl, Req.Body.ControllerID)

	// refresh ControllerInstance
	var ci models.ControllerInstace
	db.Where("controller_id = ?", ctrl.ID).Where("deleted_at is null").First(&ci)
	ControllerController := ws_controller_namespace.WsControllerController{}
	var newci = ControllerController.RefreshInstance(ctrl, ci, 0, 0, 0)

	var sensor models.Sensor
	sensor.ControllerInstaceID = newci.ID
	sensor.Name = Req.Body.Name
	sensor.SensorType = Req.Body.Type

	db.Create(&sensor)

	return wsservice.WsResponse[interface{}]{
		Namespace: "sensor",
		Endpoint:  "create",
		Body:      sensor,
	}
}
func (w WsSensorController) Update(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		ID   string            `json:"id"`
		Name string            `json:"name"`
		Type models.SensorType `json:"type"`
	}

	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var sensor models.Sensor
	db.Find(&sensor, Req.Body.ID)

	sensor.Name = Req.Body.Name
	sensor.SensorType = Req.Body.Type

	db.Save(&sensor)

	return wsservice.WsResponse[interface{}]{
		Namespace: "sensor",
		Endpoint:  "update",
		Body:      sensor,
	}
}
func (w WsSensorController) Delete(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		ID string `json:"id"`
	}

	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var sensor models.Sensor
	db.Find(&sensor, Req.Body.ID)

	var ci models.ControllerInstace
	db.Find(&ci, sensor.ControllerInstaceID)

	var ctrl models.Controller
	db.Find(&ctrl, ci.ControllerID)

	// refresh ControllerInstance without this sensor
	ControllerController := ws_controller_namespace.WsControllerController{}
	ControllerController.RefreshInstance(ctrl, ci, 0, sensor.ID, 0)

	db.Delete(&sensor)

	return wsservice.WsResponse[interface{}]{
		Namespace: "sensor",
		Endpoint:  "delete",
		Body:      sensor,
	}
}
