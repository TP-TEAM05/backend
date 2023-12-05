package dataLogging

import (
	"fmt"
	"recofiit/models"
	"recofiit/services/database"
)

func populateSensors() {
	sensorTypes := []models.SensorType{
		models.FRONT_LIDAR,
		models.FRONT_ULTRASONIC,
		models.REAR_ULTASONIC,
		models.WHEEL_SPEED,
		models.GPS_LOCATION,
		models.GPS_SPEED,
		models.GPS_DIRECTION,
		models.MAGNETOMETER_DIRECTION,
	}

	for _, sensorType := range sensorTypes {
		sensor := models.Sensor{
			Name:       fmt.Sprintf("BaseSensor", sensorType),
			SensorType: sensorType,
			// You might need to add other fields depending on your actual model structure
		}
		database.GetDB().FirstOrCreate(&sensor, models.Sensor{SensorType: sensorType})
	}
}

func Init() {
	populateSensors()
}
