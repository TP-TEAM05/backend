package dataLogging

import (
	"recofiit/models"
	"time"
)

func LogData(datagram models.UpdatePositionVehicleDatagram) {
	var carController = NewCarController()
	var measurementController = NewMeasurementController()

	carSessionID, err := carController.CreateSessionAndCarSession(1, "Test Session")
	if err != nil {
		panic("Failed to create session and car session")
	}

	SaveMeasurement(*measurementController, carSessionID, "FRONT_LIDAR", datagram.Vehicle.FrontLidar, nil)
	SaveMeasurement(*measurementController, carSessionID, "FRONT_ULTRASONIC", datagram.Vehicle.FrontUltrasonic, nil)
	SaveMeasurement(*measurementController, carSessionID, "REAR_ULTRASONIC", datagram.Vehicle.RearUltrasonic, nil)
	SaveMeasurement(*measurementController, carSessionID, "WHEEL_SPEED", datagram.Vehicle.WheelSpeed, nil)
	SaveMeasurement(*measurementController, carSessionID, "GPS_LOCATION", datagram.Vehicle.GPSLocation.Latitude, &datagram.Vehicle.GPSLocation.Longitude)
	SaveMeasurement(*measurementController, carSessionID, "GPS_SPEED", datagram.Vehicle.GPSSpeed, nil)
	SaveMeasurement(*measurementController, carSessionID, "GPS_DIRECTION", datagram.Vehicle.GPSDirection, nil)
	SaveMeasurement(*measurementController, carSessionID, "MAGNETOMETER_DIRECTION", datagram.Vehicle.MagnetometerDirection, nil)
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
