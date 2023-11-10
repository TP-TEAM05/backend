package ws_car_namespace

import "fmt"

type WsCarController struct{}

func (w WsCarController) Get()    { fmt.Println("GET CAR") }
func (w WsCarController) List()   { fmt.Println("LIST CAR") }
func (w WsCarController) Create() { fmt.Println("CREATE CAR") }
func (w WsCarController) Update() { fmt.Println("UPDATE CAR") }
func (w WsCarController) Delete() { fmt.Println("DELETE CAR") }
