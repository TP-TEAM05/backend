package ws_car_namespace

import (
	"fmt"
	"recofiit/models"
	"recofiit/services/database"
	wsservice "recofiit/services/wsService"
	"strconv"
)

type WsCarController struct{}

func (w WsCarController) Get(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("GET CAR")
	type Body struct {
		Vin string `json:"vin"`
	}

	var Req wsservice.WsRequestPrepared[Body]
	Req.Parse(req)

	db := database.GetDB()
	var car models.Car
	car.Vin = Req.Body.Vin
	db.Find(&car)

	return wsservice.WsResponse[interface{}]{
		Namespace: "car",
		Endpoint:  "get",
		Body:      car,
	}
}
func (w WsCarController) List(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("LIST CAR")
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
	fmt.Println("CREATE CAR")
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

	fmt.Println("CREATED CAR " + strconv.Itoa(int(car.ID)))

	return wsservice.WsResponse[interface{}]{
		Namespace: "car",
		Endpoint:  "create",
		Body:      car,
	}
}
func (w WsCarController) Update(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("UPDATE CAR")
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
	db.Find(&car)

	car.Vin = Req.Body.Vin
	car.Name = Req.Body.Name
	car.Color = Req.Body.Color

	db.Save(&car)

	fmt.Println("UPDATED CAR " + strconv.Itoa(int(car.ID)))

	return wsservice.WsResponse[interface{}]{
		Namespace: "car",
		Endpoint:  "update",
		Body:      car,
	}
}
func (w WsCarController) Delete(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("DELETE CAR")
	type Body struct {
		Vin string `json:"vin"`
	}
	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var car models.Car
	car.Vin = Req.Body.Vin
	db.Find(&car)

	db.Delete(&car)

	fmt.Println("DELETED CAR " + strconv.Itoa(int(car.ID)))

	return wsservice.WsResponse[interface{}]{
		Namespace: "car",
		Endpoint:  "delete",
		Body:      car,
	}
}
