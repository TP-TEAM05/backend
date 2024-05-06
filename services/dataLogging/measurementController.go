package dataLogging

import (
	"recofiit/models"
	"recofiit/services/database"
)

type MeasurementController struct{}

func NewMeasurementController() *MeasurementController {
	return &MeasurementController{}
}

func (r *MeasurementController) InsertMeasurement(measurement models.Measurement) error {
	return database.GetDB().Create(&measurement).Error
}

func (r *MeasurementController) GetSensorID(sensors []models.Sensor, sensorType models.SensorType) *uint {
	for _, sensor := range sensors {
		if sensor.SensorType == sensorType {
			return &sensor.ID
		}
	}

	return nil
}
