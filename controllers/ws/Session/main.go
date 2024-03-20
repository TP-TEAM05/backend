package ws_session_namespace

import (
	"recofiit/models"
	"recofiit/services/database"
	wsservice "recofiit/services/wsService"
	"time"
)

type WsSessionController struct{}

func (w WsSessionController) Get(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		ID uint `json:"id"`
	}

	var Req wsservice.WsRequestPrepared[Body]
	Req.Parse(req)

	db := database.GetDB()
	var session models.Session
	db.First(&session, Req.Body.ID)

	return wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "get",
		Body:      session,
	}
}
func (w WsSessionController) List(req []byte) wsservice.WsResponse[interface{}] {
	db := database.GetDB()
	var sessions []models.Session
	db.Preload("Cars").Find(&sessions)
	return wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "list",
		Body:      sessions,
	}
}
func (w WsSessionController) Create(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		Cars []string `json:"cars"`
		Name string   `json:"name"`
	}

	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	var cars []models.Car

	db := database.GetDB()
	db.Where("vin IN ?", Req.Body.Cars).Find(&cars)

	var session models.Session

	session.Name = Req.Body.Name
	db.Create(&session)

	var carSessions []models.CarSession
	var carIds []uint
	for _, car := range cars {
		var cs models.CarSession
		cs.CarID = car.ID
		cs.SessionID = session.ID
		carSessions = append(carSessions, cs)
		carIds = append(carIds, car.ID)
	}
	db.Create(&carSessions)

	var carControllers []models.CarController
	db.Where("car_id IN ?", carIds).Where("deleted_at IS NULL").Find(&carControllers)

	var carSessionControllers []models.CarSessionController
	for _, carController := range carControllers {
		var csc models.CarSessionController
		for _, carSession := range carSessions {
			if carSession.CarID == carController.CarID {
				csc.CarSessionID = carSession.ID
			}
		}
		csc.ControllerInstanceID = carController.ControllerInstanceID
		carSessionControllers = append(carSessionControllers, csc)
	}
	db.Create(&carSessionControllers)

	var s models.Session
	db.Preload("Cars").First(&s, session.ID)

	return wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "create",
		Body:      s,
	}

}
func (w WsSessionController) Start(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		ID uint `json:"id"`
	}

	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var session models.Session
	db.First(&session, Req.Body.ID)

	var timenow = time.Now()

	session.StartedAt = &timenow

	db.Save(&session)

	return wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "update",
		Body:      session,
	}
}
func (w WsSessionController) End(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		ID uint `json:"id"`
	}

	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var session models.Session
	db.First(&session, Req.Body.ID)

	var timenow = time.Now()

	session.EndedAt = &timenow

	db.Save(&session)

	return wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "update",
		Body:      session,
	}
}
func (w WsSessionController) Delete(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		ID uint `json:"id"`
	}
	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var session models.Session
	db.First(&session, Req.Body.ID)

	db.Delete(&session)

	return wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "delete",
		Body:      session,
	}
}
