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

	SaveMeasurement(*measurementController, carSessionID, "GPS", datagram.Vehicle.Latitude, &datagram.Vehicle.Longitude)
	SaveMeasurement(*measurementController, carSessionID, "CAR_DIRECTION", datagram.Vehicle.CarDirection, nil)
	SaveMeasurement(*measurementController, carSessionID, "DISTANCE_ULTRASONIC", datagram.Vehicle.DistanceUltrasonic, nil)
	SaveMeasurement(*measurementController, carSessionID, "REAR_ULTRASONIC", datagram.Vehicle.RearUltrasonic, nil)
	SaveMeasurement(*measurementController, carSessionID, "DISTANCE_LIDAR", datagram.Vehicle.DistanceLidar, nil)
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
