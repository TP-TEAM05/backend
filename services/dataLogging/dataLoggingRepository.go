package dataLogging

import (
	"recofiit/models"
)

func LogData(datagram models.UpdatePositionVehicleDatagram) {
	var carController = NewCarController()
	var measurementController = NewMeasurementController()

	// Creating a session and car session
	carSessionID, err := carController.CreateSessionAndCarSession(1, "Test Session")
	if err != nil {
		panic("Failed to create session and car session")
	}

	// Getting a sensor ID
	sensorID, err := measurementController.GetSensorID("BaseSensor", "FRONT_LIDAR")
	if err != nil {
		panic("Failed to get sensor ID")
	}

	// Inserting a measurement
	measurement := models.Measurement{CarSessionID: carSessionID, SensorID: sensorID, Data1: "23.45", Data2: "67.89"}
	err = measurementController.InsertMeasurement(measurement)
	if err != nil {
		panic("Failed to insert measurement")
	}
}
