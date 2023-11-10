package routes

import (
	ws_car_namespace "recofiit/controllers/ws/Car"
	ws_session_namespace "recofiit/controllers/ws/Session"
	wsservice "recofiit/services/wsService"
)

func SetupWs(m *wsservice.ClientManager) {

	CarController := ws_car_namespace.WsCarController{}
	SessionController := ws_session_namespace.WsSessionController{}
	r := m.Router

	// Car routes
	r.Register("car", "get", CarController.Get)
	r.Register("car", "list", CarController.List)
	r.Register("car", "create", CarController.Create)
	r.Register("car", "update", CarController.Update)
	r.Register("car", "delete", CarController.Delete)

	// Session routes

	r.Register("session", "get", SessionController.Get)
	r.Register("session", "list", SessionController.List)
	r.Register("session", "create", SessionController.Create)
	r.Register("session", "update", SessionController.Update)
	r.Register("session", "delete", SessionController.Delete)
}
