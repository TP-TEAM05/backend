package dataLogging

import (
	"recofiit/models"
	"recofiit/services/database"
)

type MeasurementController struct{}

func NewMeasurementController() *MeasurementController {
	return &MeasurementController{}
}

// InsertMeasurement inserts a new measurement record
func (r *MeasurementController) InsertMeasurement(measurement models.Measurement) error {
	return database.GetDB().Create(&measurement).Error
}

// GetSensorID retrieves a sensor ID based on its name and type
func (r *MeasurementController) GetSensorID(name string, sensorType string) (uint, error) {
	var sensor models.Sensor
	err := database.GetDB().Where("name = ? AND sensor_type = ?", name, sensorType).First(&sensor).Error
	if err != nil {
		return 0, err
	}
	return sensor.ID, nil
}
