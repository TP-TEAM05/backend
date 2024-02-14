package ws_session_namespace

import (
	"fmt"
	"recofiit/models"
	"recofiit/services/database"
	wsservice "recofiit/services/wsService"
	"strconv"
)

type WsSessionController struct{}

func (w WsSessionController) Get(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("GET SESSION")
	type Body struct {
		ID uint `json:"id"`
	}

	var Req wsservice.WsRequestPrepared[Body]
	Req.Parse(req)

	db := database.GetDB()
	var session models.Session
	db.Find(&session, Req.Body.ID)

	return wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "get",
		Body:      session,
	}
}
func (w WsSessionController) List(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("LIST SESSION")
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
	fmt.Println("CREATE SESSION")
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
	session.Cars = cars

	db.Create(&session)

	var cis []models.ControllerInstace
	db.Where("car_id IN ?", Req.Body.Cars).Find(&cis)

	var cscs []models.CarSessionController
	for _, ci := range cis {
		var csc models.CarSessionController
		csc.CarSessionID = session.ID
		csc.ControllerInstanceID = ci.ID
		cscs = append(cscs, csc)
	}
	db.Create(&cscs)

	fmt.Println("CREATED SESSION " + strconv.Itoa(int(session.ID)))

	return wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "create",
		Body:      session,
	}

}
func (w WsSessionController) Update(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("UPDATE SESSION")
	type Body struct {
		ID        uint     `json:"id"`
		Cars      []string `json:"cars"`
		Name      string   `json:"name"`
		StartedAt *string  `json:"started_at"`
		EndedAt   *string  `json:"ended_at"`
	}

	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	var cars []models.Car

	db := database.GetDB()
	db.Where("vin IN ?", Req.Body.Cars).Find(&cars)

	var session models.Session
	db.Find(&session, Req.Body.ID)

	session.Name = Req.Body.Name
	session.Cars = cars
	session.StartedAt = Req.Body.StartedAt
	session.EndedAt = Req.Body.EndedAt

	db.Save(&session)

	fmt.Println("UPDATED SESSION " + strconv.Itoa(int(session.ID)))

	return wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "update",
		Body:      session,
	}
}
func (w WsSessionController) Delete(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("DELETE SESSION")
	type Body struct {
		ID uint `json:"id"`
	}
	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var session models.Session
	db.Find(&session, Req.Body.ID)

	db.Delete(&session)

	fmt.Println("DELETED SESSION " + strconv.Itoa(int(session.ID)))

	return wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "delete",
		Body:      session,
	}
}
