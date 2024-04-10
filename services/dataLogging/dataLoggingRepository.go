package dataLogging

import (
	"recofiit/models"
	"recofiit/services/database"
	"strconv"
	"time"

	api "github.com/ReCoFIIT/integration-api"
)

func LogData(datagram api.UpdateVehicleDatagram) {
	var measurementController = NewMeasurementController()

	db := database.GetDB()
	var session models.Session
	result := db.Where("started_at is not null").Where("ended_at is null").First(&session)

	if result.Error != nil {
		panic("Failed to find started session")
	}

	var car models.Car
	result = db.Where("vin", datagram.Vehicle.Vin).Where("deleted_at is null").First(&car)
	if result.Error != nil {
		var count int64

		db.Model(&models.Car{}).Count(&count)

		car = models.Car{
			Vin:   datagram.Vehicle.Vin,
			Name:  "Car " + strconv.FormatInt(count+1, 10),
			Color: "#" + datagram.Vehicle.Vin[:6],
		}
		db.Create(&car)
	}

	var carSession models.CarSession
	result = db.Where("car_id", car.ID).Where("session_id", session.ID).Where("deleted_at is null").First(&carSession)

	if result.Error != nil {
		carSession = models.CarSession{
			CarID:     car.ID,
			SessionID: session.ID,
		}
		db.Create(&carSession)
	}

	var carSessionID = carSession.ID

	SaveMeasurement(*measurementController, carSessionID, "GPS_LOCATION", datagram.Vehicle.Latitude, &datagram.Vehicle.Longitude)
	SaveMeasurement(*measurementController, carSessionID, "GPS_DIRECTION", datagram.Vehicle.GpsDirection, nil)
	SaveMeasurement(*measurementController, carSessionID, "FRONT_ULTRASONIC", datagram.Vehicle.FrontUltrasonic, nil)
	SaveMeasurement(*measurementController, carSessionID, "REAR_ULTRASONIC", datagram.Vehicle.RearUltrasonic, nil)
	SaveMeasurement(*measurementController, carSessionID, "FRONT_LIDAR", datagram.Vehicle.FrontLidar, nil)
	SaveMeasurement(*measurementController, carSessionID, "SPEED_FRONT_LEFT", datagram.Vehicle.SpeedFrontLeft, nil)
	SaveMeasurement(*measurementController, carSessionID, "SPEED_FRONT_RIGHT", datagram.Vehicle.SpeedFrontRight, nil)
	SaveMeasurement(*measurementController, carSessionID, "SPEED_REAR_LEFT", datagram.Vehicle.SpeedRearLeft, nil)
	SaveMeasurement(*measurementController, carSessionID, "SPEED_REAR_RIGHT", datagram.Vehicle.SpeedRearRight, nil)
}

func SaveMeasurement(measurementController MeasurementController, carSessionID uint, sensorName string, data1 float32, data2 *float32) {
	sensorID, err := measurementController.GetSensorID("BaseSensor", sensorName)
	if err != nil {
		return
	}

	var data2Value float32 = 0

	if data2 != nil {
		data2Value = *data2
	}

	measurement := models.Measurement{CarSessionID: carSessionID, CreatedAt: &time.Time{}, SensorID: sensorID, Data1: data1, Data2: data2Value}
	err = measurementController.InsertMeasurement(measurement)
	if err != nil {
		return
	}
}
