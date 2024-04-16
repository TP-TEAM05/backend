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

func (r *MeasurementController) GetSensorID(controllerInstanceIDs []uint, sensorType string) (uint, error) {
	var sensor models.Sensor
	if err := database.GetDB().Where("controller_instance_id in ? AND sensor_type = ?", controllerInstanceIDs, sensorType).First(&sensor).Error; err != nil {
		return 0, err
	}
	return sensor.ID, nil
}
