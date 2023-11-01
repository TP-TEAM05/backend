package database

import (
	"recofiit/models"
)

func GetModels() []interface{} {
	return []interface{}{
		&models.Car{},
		&models.Controller{},
		&models.Firmware{},
		&models.Meassurement{},
		&models.Sensor{},
		&models.SensorData{},
		&models.Session{},
		&models.CarController{},
		&models.CarSession{},
		&models.CarSessionController{},
		&models.ControllerInstace{},
	}
}
