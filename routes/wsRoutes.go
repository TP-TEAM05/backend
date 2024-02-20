package routes

import (
	ws_car_namespace "recofiit/controllers/ws/Car"
	ws_controller_namespace "recofiit/controllers/ws/Controller"
	ws_firmware_namespace "recofiit/controllers/ws/Firmware"
	ws_sensor_namespace "recofiit/controllers/ws/Sensor"
	ws_session_namespace "recofiit/controllers/ws/Session"
	wsservice "recofiit/services/wsService"
)

func SetupWs(m *wsservice.ClientManager) {

	CarController := ws_car_namespace.WsCarController{}
	ControllerController := ws_controller_namespace.WsControllerController{}
	FirmwareController := ws_firmware_namespace.WsFirmwareController{}
	SensorController := ws_sensor_namespace.WsSensorController{}
	SessionController := ws_session_namespace.WsSessionController{}
	r := m.Router

	// Car routes
	r.Register("car", "get", CarController.Get)
	r.Register("car", "list", CarController.List)
	r.Register("car", "create", CarController.Create)
	r.Register("car", "update", CarController.Update)
	r.Register("car", "delete", CarController.Delete)

	// Controller routes
	r.Register("controller", "get", ControllerController.Get)
	r.Register("controller", "list", ControllerController.List)
	r.Register("controller", "create", ControllerController.Create)
	r.Register("controller", "update", ControllerController.Update)
	r.Register("controller", "delete", ControllerController.Delete)

	// Firmware routes
	r.Register("firmware", "get", FirmwareController.Get)
	r.Register("firmware", "list", FirmwareController.List)
	r.Register("firmware", "create", FirmwareController.Create)
	r.Register("firmware", "update", FirmwareController.Update)
	//r.Register("firmware", "delete", FirmwareController.Delete)

	// Sensor routes
	r.Register("sensor", "get", SensorController.Get)
	r.Register("sensor", "list", SensorController.List)
	r.Register("sensor", "create", SensorController.Create)
	r.Register("sensor", "update", SensorController.Update)
	r.Register("sensor", "delete", SensorController.Delete)

	// Session routes
	r.Register("session", "get", SessionController.Get)
	r.Register("session", "list", SessionController.List)
	r.Register("session", "create", SessionController.Create)
	r.Register("session", "update", SessionController.Update)
	r.Register("session", "delete", SessionController.Delete)
}
