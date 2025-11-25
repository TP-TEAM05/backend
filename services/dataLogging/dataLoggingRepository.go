package dataLogging

import (
	"recofiit/models"
	"recofiit/services/database"
	"recofiit/services/redis"
	"strconv"
	"time"

	api "github.com/TP-TEAM05/integration-api"
	"github.com/getsentry/sentry-go"
)

func LogData(datagram api.UpdateVehicleDatagram) {
	var measurementController = NewMeasurementController()

	db := database.GetDB()
	var session models.Session
	result := db.Where("started_at is not null").Where("ended_at is null").First(&session)

	if result.Error != nil {
		// Session not found
		sentry.CaptureMessage("Session not found for vehicle: " + datagram.Vehicle.Vin)
		return
	}

	var vehicleConfigKey = datagram.Vehicle.Vin + "-session-" + strconv.Itoa(int(session.ID))

	vehicleConfig := redis.GetVehicleConfig(vehicleConfigKey)
	if vehicleConfig == nil {

		var count int64
		db.Model(&models.Car{}).Where("vin", datagram.Vehicle.Vin).Where("deleted_at is null").Count(&count)

		var car models.Car
		if count == 0 {
			// create new car
			car.Vin = datagram.Vehicle.Vin
			car.Name = "Car"
			car.Color = "#ff0000"

			db.Create(&car)
		} else {
			// find existing car
			db.Where("vin = ?", datagram.Vehicle.Vin).Where("deleted_at is null").First(&car)
		}

		var carSession models.CarSession
		db.Model(&models.CarSession{}).Where("car_id", car.ID).Where("session_id", session.ID).Where("deleted_at is null").Count(&count)

		if count == 0 {
			// create new car-session
			carSession = models.CarSession{
				CarID:     car.ID,
				SessionID: session.ID,
			}
			db.Create(&carSession)
		} else {
			// find existing car-session
			db.Where("car_id", car.ID).Where("session_id", session.ID).Where("deleted_at is null").First(&carSession)
		}

		var carSessionID = carSession.ID

		var controllerInstanceIDs []uint

		db.Model(&models.CarSessionController{}).Where("car_session_id", carSessionID).Where("deleted_at is null").Count(&count)
		if count == 0 {
			// no car controller is registered for the given session -> add new
			var carController models.CarController
			var controllerInstance models.ControllerInstance
			result := db.Where("car_id", car.ID).Where("deleted_at is null").First(&carController)
			if result.Error != nil {
				// no controller exist for the given car -> add new
				var controller = models.Controller{
					Name: "Controller",
					Type: "Base Controller",
				}
				db.Create(&controller)

				var firmware models.Firmware
				result := db.Last(&firmware)
				if result.Error != nil {
					// No firmware exists, create a default one
					firmware = models.Firmware{
						Version:     "Default 1.0",
						Description: "Auto-created default firmware",
					}
					db.Create(&firmware)
				}

				controllerInstance = models.ControllerInstance{
					FirmwareID:   firmware.ID,
					ControllerID: controller.ID,
				}
				db.Create(&controllerInstance)

				carController = models.CarController{
					CarID:                car.ID,
					ControllerInstanceID: controllerInstance.ID,
				}
				db.Create(&carController)
			} else {
				// some controller exists on the car
				db.Where("id", carController.ControllerInstanceID).First(&controllerInstance)
			}
			controllerInstanceIDs = []uint{controllerInstance.ID}

			// add the controller to this car in the given session
			var carSessionController = models.CarSessionController{
				CarSessionID:         carSessionID,
				ControllerInstanceID: controllerInstance.ID,
			}
			db.Create(&carSessionController)

			db.Model(&models.Sensor{}).Where("controller_instance_id", controllerInstance.ID).Where("deleted_at is null").Count(&count)
			if count == 0 {
				// no sensors are defined -> add new
				for _, sensorType := range models.SensorTypes {
					var sensor = models.Sensor{
						ControllerInstanceID: controllerInstance.ID,
						Name:                 "BaseSensor",
						SensorType:           sensorType,
					}
					db.Create(&sensor)
				}
			}
		} else {
			// some controller is already registered -> use that one
			var carSessionController []models.CarSessionController
			db.Where("car_session_id", carSessionID).Where("deleted_at is null").Find(&carSessionController)

			controllerInstanceIDs = make([]uint, len(carSessionController))
			for _, csc := range carSessionController {
				controllerInstanceIDs = append(controllerInstanceIDs, csc.ControllerInstanceID)
			}
		}
		var sensors []models.Sensor
		database.GetDB().Where("controller_instance_id in ?", controllerInstanceIDs).Find(&sensors)

		vehicleConfig = &models.VehicleConfig{
			Car:        car,
			CarSession: carSession,
			Sensors:    sensors,
		}

		err := redis.SaveVehicleConfig(vehicleConfigKey, vehicleConfig)
		if err != nil {
			sentry.CaptureException(err)
			return
		}
	}

	SaveMeasurement(*measurementController, vehicleConfig.CarSession.ID, vehicleConfig.Sensors, "GPS_LOCATION", datagram.Vehicle.Latitude, &datagram.Vehicle.Longitude)
	SaveMeasurement(*measurementController, vehicleConfig.CarSession.ID, vehicleConfig.Sensors, "GPS_DIRECTION", datagram.Vehicle.GpsDirection, nil)
	SaveMeasurement(*measurementController, vehicleConfig.CarSession.ID, vehicleConfig.Sensors, "FRONT_ULTRASONIC", datagram.Vehicle.FrontUltrasonic, nil)
	SaveMeasurement(*measurementController, vehicleConfig.CarSession.ID, vehicleConfig.Sensors, "REAR_ULTRASONIC", datagram.Vehicle.RearUltrasonic, nil)
	SaveMeasurement(*measurementController, vehicleConfig.CarSession.ID, vehicleConfig.Sensors, "FRONT_LIDAR", datagram.Vehicle.FrontLidar, nil)
	SaveMeasurement(*measurementController, vehicleConfig.CarSession.ID, vehicleConfig.Sensors, "SPEED_FRONT_LEFT", datagram.Vehicle.SpeedFrontLeft, nil)
	SaveMeasurement(*measurementController, vehicleConfig.CarSession.ID, vehicleConfig.Sensors, "SPEED_FRONT_RIGHT", datagram.Vehicle.SpeedFrontRight, nil)
	SaveMeasurement(*measurementController, vehicleConfig.CarSession.ID, vehicleConfig.Sensors, "SPEED_REAR_LEFT", datagram.Vehicle.SpeedRearLeft, nil)
	SaveMeasurement(*measurementController, vehicleConfig.CarSession.ID, vehicleConfig.Sensors, "SPEED_REAR_RIGHT", datagram.Vehicle.SpeedRearRight, nil)
	SaveMeasurement(*measurementController, vehicleConfig.CarSession.ID, vehicleConfig.Sensors, "VOLTAGE0", datagram.Vehicle.Voltage0, nil)
	SaveMeasurement(*measurementController, vehicleConfig.CarSession.ID, vehicleConfig.Sensors, "VOLTAGE1", datagram.Vehicle.Voltage1, nil)
	SaveMeasurement(*measurementController, vehicleConfig.CarSession.ID, vehicleConfig.Sensors, "VOLTAGE2", datagram.Vehicle.Voltage2, nil)
}

func SaveMeasurement(measurementController MeasurementController, carSessionID uint, sensors []models.Sensor, sensorType models.SensorType, data1 float32, data2 *float32) {
	sensorID := measurementController.GetSensorID(sensors, sensorType)
	if sensorID == nil {
		return
	}

	var data2Value float32 = 0

	if data2 != nil {
		data2Value = *data2
	}

	currentTime := time.Now()
	measurement := models.Measurement{CarSessionID: carSessionID, CreatedAt: &currentTime, SensorID: *sensorID, Data1: data1, Data2: data2Value}
	err := measurementController.InsertMeasurement(measurement)
	if err != nil {
		sentry.CaptureException(err)
		return
	}
}
