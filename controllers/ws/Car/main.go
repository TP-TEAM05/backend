package ws_car_namespace

import (
	"recofiit/models"
	"recofiit/services/database"
	wsservice "recofiit/services/wsService"
)

type WsCarController struct{}

func (w WsCarController) Get(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		Vin string `json:"vin"`
	}

	var Req wsservice.WsRequestPrepared[Body]
	Req.Parse(req)

	db := database.GetDB()
	var car models.Car
	db.Where("vin = ?", Req.Body.Vin).First(&car)

	return wsservice.WsResponse[interface{}]{
		Namespace: "car",
		Endpoint:  "get",
		Body:      car,
	}
}
func (w WsCarController) List(req []byte) wsservice.WsResponse[interface{}] {
	db := database.GetDB()
	var cars []models.Car
	db.Find(&cars)

	return wsservice.WsResponse[interface{}]{
		Namespace: "car",
		Endpoint:  "list",
		Body:      cars,
	}
}
func (w WsCarController) Create(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		Vin   string `json:"vin"`
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var car models.Car
	car.Vin = Req.Body.Vin
	car.Name = Req.Body.Name
	car.Color = Req.Body.Color

	db.Create(&car)

	return wsservice.WsResponse[interface{}]{
		Namespace: "car",
		Endpoint:  "create",
		Body:      car,
	}
}
func (w WsCarController) Update(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		Vin   string `json:"vin"`
		Name  string `json:"name"`
		Color string `json:"color"`
	}
	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var car models.Car
	db.Where("vin = ?", Req.Body.Vin).First(&car)

	car.Name = Req.Body.Name
	car.Color = Req.Body.Color

	db.Save(&car)

	return wsservice.WsResponse[interface{}]{
		Namespace: "car",
		Endpoint:  "update",
		Body:      car,
	}
}
func (w WsCarController) Delete(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		Vin string `json:"vin"`
	}
	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var car models.Car
	db.Where("vin = ?", Req.Body.Vin).First(&car)

	db.Delete(&car)

	return wsservice.WsResponse[interface{}]{
		Namespace: "car",
		Endpoint:  "delete",
		Body:      car,
	}
}
