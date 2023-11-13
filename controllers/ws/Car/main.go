package ws_car_namespace

import (
	"fmt"
	"recofiit/models"
	"recofiit/services/database"
	wsservice "recofiit/services/wsService"
)

type WsCarController struct{}

func (w WsCarController) Get(c *wsservice.Client, req *wsservice.WsRequest, res *wsservice.WsResponse) {
	fmt.Println("GET CAR")
}
func (w WsCarController) List(c *wsservice.Client, req *wsservice.WsRequest, res *wsservice.WsResponse) {
	fmt.Println("LIST CAR")
	db := database.GetDB()
	var cars []models.Car
	db.Find(&cars)

	res.Body = cars
}
func (w WsCarController) Create(c *wsservice.Client, req *wsservice.WsRequest, res *wsservice.WsResponse) {
	fmt.Println("CREATE CAR")
}
func (w WsCarController) Update(c *wsservice.Client, req *wsservice.WsRequest, res *wsservice.WsResponse) {
	fmt.Println("UPDATE CAR")
}
func (w WsCarController) Delete(c *wsservice.Client, req *wsservice.WsRequest, res *wsservice.WsResponse) {
	fmt.Println("DELETE CAR")
}
