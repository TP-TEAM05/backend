package ws_car_namespace

import (
	"fmt"
	"recofiit/models"
	"recofiit/services/database"
	wsservice "recofiit/services/wsService"
)

type WsCarController struct{}

func (w WsCarController) Get(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("GET CAR")
	return wsservice.WsResponse[interface{}]{}
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
	return wsservice.WsResponse[interface{}]{}
}
func (w WsCarController) Update(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("UPDATE CAR")
	return wsservice.WsResponse[interface{}]{}
}
func (w WsCarController) Delete(req []byte) wsservice.WsResponse[interface{}] {
	fmt.Println("DELETE CAR")
	return wsservice.WsResponse[interface{}]{}
}
