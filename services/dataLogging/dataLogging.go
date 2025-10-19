package dataLogging

import (
	"fmt"
	"recofiit/models"
	"recofiit/services/database"

	"github.com/getsentry/sentry-go"
)

func populateControllers() {
	controllers := []models.Controller{
		{
			Name:        "TestController1",
			Type:        "TestType1",
			Description: "This is a test controller 1",
		},
		{
			Name:        "TestController2",
			Type:        "TestType2",
			Description: "This is a test controller 2",
		},
	}

	for _, controller := range controllers {
		database.GetDB().Create(&controller)
	}
}

func populateFirmwares() {
	firmwares := []models.Firmware{
		{
			Version:     "Version 1.0",
			Description: "First firmware version",
		},
		{
			Version:     "Version 2.0",
			Description: "Second firmware version",
		},
	}

	for _, firmware := range firmwares {
		database.GetDB().Create(&firmware)
	}
}

func populateControllerInstances() {
	controllerInstances := []models.ControllerInstance{
		{
			FirmwareID:   1,
			ControllerID: 1,
		},
		{
			FirmwareID:   2,
			ControllerID: 2,
		},
	}

	for _, instance := range controllerInstances {
		database.GetDB().Create(&instance)
	}
}

func getRandomControllerInstanceID() (uint, error) {
	var id uint
	if err := database.GetDB().Model(&models.ControllerInstance{}).Select("id").Order("RANDOM()").Limit(1).Pluck("id", &id).Error; err != nil {
		return 0, err
	}
	return id, nil
}

func populateSensors() {
	sensorTypes := []models.SensorType{
		models.GPS_LOCATION,
		models.GPS_SPEED,
		models.GPS_DIRECTION,
		models.GPS_ADDITIONAL,
		models.FRONT_ULTRASONIC,
		models.REAR_ULTRASONIC,
		models.FRONT_LIDAR,
		models.SPEED_FRONT_LEFT,
		models.SPEED_FRONT_RIGHT,
		models.SPEED_REAR_LEFT,
		models.SPEED_REAR_RIGHT,
		models.MAGNETOMETER_DIRECTION,
	}

	id, err := getRandomControllerInstanceID()
	if err != nil {
		sentry.CaptureException(err)
		// Handle error (log, return, etc.)
		fmt.Println("Error getting ControllerInstanceID:", err)
		return
	}

	for _, sensorType := range sensorTypes {
		sensor := models.Sensor{
			ControllerInstanceID: id,
			Name:                 "BaseSensor",
			SensorType:           sensorType,
		}
		database.GetDB().FirstOrCreate(&sensor, models.Sensor{SensorType: sensorType})
	}
}

func populateCars() {
	cars := []models.Car{
		{
			Vin:   "C4RF117S7U0000001",
			Name:  "Car 1",
			Color: "#0033ff",
		},
		{
			Vin:   "C4RF117S7U0000002",
			Name:  "Car 2",
			Color: "#009900",
		},
	}

	for _, car := range cars {
		database.GetDB().Create(&car)
	}

}

func Init() {
	// populateControllers()
	// populateFirmwares()
	// populateControllerInstances()
	// populateSensors()
	// populateCars()
}
