package ws_car_namespace

import (
	"gorm.io/gorm"
	"recofiit/models"
	"recofiit/services/database"
	wsservice "recofiit/services/wsService"
	"time"
)

type WsCarController struct{}

type ExtendedCar struct {
	ID                 uint           `json:"id"`
	Vin                string         `json:"vin"`
	Name               string         `json:"name"`
	Color              string         `json:"color"`
	IsControlledByUser bool           `json:"is_controlled_by_user"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

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
		Body:      w.ExtendCar(car, db),
	}
}
func (w WsCarController) List(req []byte) wsservice.WsResponse[interface{}] {
	db := database.GetDB()
	var cars []models.Car
	db.Find(&cars)

	var ecs = make([]ExtendedCar, 0)
	for _, c := range cars {
		ecs = append(ecs, w.ExtendCar(c, db))
	}

	return wsservice.WsResponse[interface{}]{
		Namespace: "car",
		Endpoint:  "list",
		Body:      ecs,
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
		Body:      w.ExtendCar(car, db),
	}
}
func (w WsCarController) Update(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		Vin                string `json:"vin"`
		Name               string `json:"name"`
		Color              string `json:"color"`
		IsControlledByUser bool   `json:"is_controlled_by_user"`
		SessionID          *uint  `json:"session_id"`
	}
	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var car models.Car
	db.Where("vin = ?", Req.Body.Vin).First(&car)

	car.Name = Req.Body.Name
	car.Color = Req.Body.Color

	db.Save(&car)

	if Req.Body.SessionID != nil {
		var carSession models.CarSession
		db.Where("session_id = ?", *Req.Body.SessionID).Where("car_id", car.ID).First(&carSession)
		carSession.IsControlledByUser = Req.Body.IsControlledByUser
		db.Save(&carSession)
	}

	return wsservice.WsResponse[interface{}]{
		Namespace: "car",
		Endpoint:  "update",
		Body:      w.ExtendCar(car, db),
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

func (w WsCarController) ExtendCar(
	car models.Car, db *gorm.DB) ExtendedCar {
	var count int64
	db.Model(&models.CarSession{}).Where("car_id", car.ID).Where("is_controlled_by_user = true").Count(&count)

	return ExtendedCar{
		ID:                 car.ID,
		Vin:                car.Vin,
		Name:               car.Name,
		Color:              car.Color,
		IsControlledByUser: count > 0,
		CreatedAt:          car.CreatedAt,
		UpdatedAt:          car.UpdatedAt,
		DeletedAt:          car.DeletedAt,
	}
}
