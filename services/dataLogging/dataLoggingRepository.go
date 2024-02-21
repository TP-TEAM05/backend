package dataLogging

import (
	"recofiit/models"
	"time"
)

func LogData(datagram models.UpdateVehicleDatagram) {
	var carController = NewCarController()
	var measurementController = NewMeasurementController()

	carSessionID, err := carController.CreateSessionAndCarSession(1, "Test Session")
	if err != nil {
		panic("Failed to create session and car session")
	}

	SaveMeasurement(*measurementController, carSessionID, "LONGITUDE", datagram.Vehicle.Latitude)
	SaveMeasurement(*measurementController, carSessionID, "LATITUDE", datagram.Vehicle.Longitude)
	SaveMeasurement(*measurementController, carSessionID, "DISTANCE_ULTRASONIC", datagram.Vehicle.DistanceUltrasonic)
	SaveMeasurement(*measurementController, carSessionID, "DISTANCE_LIDAR", datagram.Vehicle.DistanceLidar)
	SaveMeasurement(*measurementController, carSessionID, "SPEED_FRONT_LEFT", datagram.Vehicle.SpeedFrontLeft)
	SaveMeasurement(*measurementController, carSessionID, "SPEED_FRONT_RIGHT", datagram.Vehicle.SpeedFrontRight)
	SaveMeasurement(*measurementController, carSessionID, "SPEED_REAR_LEFT", datagram.Vehicle.SpeedRearLeft)
	SaveMeasurement(*measurementController, carSessionID, "SPEED_REAR_RIGHT", datagram.Vehicle.SpeedRearRight)
}

func SaveMeasurement(measurementController MeasurementController, carSessionID uint, sensorName string, data float32) {
	sensorID, err := measurementController.GetSensorID("BaseSensor", sensorName)
	if err != nil {
		return
	}

	measurement := models.Measurement{CarSessionID: carSessionID, CreatedAt: &time.Time{}, SensorID: sensorID, Data: data}
	err = measurementController.InsertMeasurement(measurement)
	if err != nil {
		return
	}
}
