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

func (r *MeasurementController) GetSensorID(name string, sensorType string) (uint, error) {
	var sensor models.Sensor
	if err := database.GetDB().Where("name = ? AND sensor_type = ?", name, sensorType).First(&sensor).Error; err != nil {
		return 0, err
	}
	return sensor.ID, nil
}
