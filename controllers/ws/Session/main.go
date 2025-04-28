package ws_session_namespace

import (
	"fmt"
	"recofiit/models"
	"recofiit/services/database"
	wsservice "recofiit/services/wsService"
	"time"
)

type WsSessionController struct{}

func (w WsSessionController) Get(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		ID uint `json:"id"`
	}

	var Req wsservice.WsRequestPrepared[Body]
	Req.Parse(req)

	db := database.GetDB()
	var session models.Session
	db.First(&session, Req.Body.ID)

	return wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "get",
		Body:      session,
	}
}
func (w WsSessionController) List(req []byte) wsservice.WsResponse[interface{}] {
	db := database.GetDB()
	var sessions []models.Session
	db.Preload("Cars").Find(&sessions)
	return wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "list",
		Body:      sessions,
	}
}
func (w WsSessionController) Create(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		Cars []string `json:"cars"`
		Name string   `json:"name"`
	}

	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	var cars []models.Car

	db := database.GetDB()
	db.Where("vin IN ?", Req.Body.Cars).Find(&cars)

	var session models.Session

	session.Name = Req.Body.Name
	db.Create(&session)

	var carSessions []models.CarSession
	var carIds []uint
	for _, car := range cars {
		var cs models.CarSession
		cs.CarID = car.ID
		cs.SessionID = session.ID
		carSessions = append(carSessions, cs)
		carIds = append(carIds, car.ID)
	}
	db.Create(&carSessions)

	var carControllers []models.CarController
	db.Where("car_id IN ?", carIds).Where("deleted_at IS NULL").Find(&carControllers)

	var carSessionControllers []models.CarSessionController
	for _, carController := range carControllers {
		var csc models.CarSessionController
		for _, carSession := range carSessions {
			if carSession.CarID == carController.CarID {
				csc.CarSessionID = carSession.ID
			}
		}
		csc.ControllerInstanceID = carController.ControllerInstanceID
		carSessionControllers = append(carSessionControllers, csc)
	}
	db.Create(&carSessionControllers)

	var s models.Session
	db.Preload("Cars").First(&s, session.ID)

	return wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "create",
		Body:      s,
	}

}
func (w WsSessionController) Start(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		ID uint `json:"id"`
	}

	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var session models.Session
	db.First(&session, Req.Body.ID)

	var timenow = time.Now()

	session.StartedAt = &timenow

	db.Save(&session)

	var s models.Session
	db.Preload("Cars").First(&s, session.ID)

	return wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "start",
		Body:      s,
	}
}
func (w WsSessionController) End(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		ID uint `json:"id"`
	}

	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var session models.Session
	db.First(&session, Req.Body.ID)

	var timenow = time.Now()

	session.EndedAt = &timenow

	db.Save(&session)

	var s models.Session
	db.Preload("Cars").First(&s, session.ID)

	return wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "end",
		Body:      s,
	}
}
func (w WsSessionController) GetMeasurements(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		SessionID uint `json:"session_id"`
	}

	var Req wsservice.WsRequestPrepared[Body]
	Req.Parse(req)

	db := database.GetDB()

	// Get all CarSessions for this session
	var carSessions []models.CarSession
	db.Where("session_id = ?", Req.Body.SessionID).Find(&carSessions)

	var carSessionIDs []uint
	for _, cs := range carSessions {
		carSessionIDs = append(carSessionIDs, cs.ID)
	}

	// Return early if there are no car sessions
	if len(carSessionIDs) == 0 {
		return wsservice.WsResponse[interface{}]{
			Namespace: "session",
			Endpoint:  "measurements",
			Body:      map[string]interface{}{"measurements": []interface{}{}},
		}
	}

	// Count the total number of measurements first
	var totalCount int64
	db.Table("measurements").
		Where("car_session_id IN ?", carSessionIDs).
		Count(&totalCount)

	// Set maximum number of measurements to return
	const maxMeasurements = 1000

	// Calculate sampling factor if needed
	var samplingFactor int = 1
	var sampleInfo string = "all"

	if totalCount > maxMeasurements {
		samplingFactor = int(totalCount)/maxMeasurements + 1
		sampleInfo = fmt.Sprintf("sampled 1/%d (total: %d)", samplingFactor, totalCount)
	}

	// Structure to hold the response data
	type MeasurementData struct {
		CarSessionID uint       `json:"car_session_id"`
		SensorType   string     `json:"sensor_type"`
		SensorID     uint       `json:"sensor_id"`
		CarID        uint       `json:"car_id"`
		VIN          string     `json:"vin"`
		Data1        float64    `json:"data1"`
		Data2        float64    `json:"data2"`
		CreatedAt    *time.Time `json:"created_at"`
	}

	// Get measurements with appropriate sampling
	var measurements []struct {
		models.Measurement
		SensorType string `json:"sensor_type"`
		CarID      uint   `json:"car_id"`
		RowNum     int    `json:"row_num"`
		VIN        string `json:"vin"`
	}

	if samplingFactor == 1 {
		// No sampling needed, get all measurements
		db.Table("measurements").
			Select("measurements.*, sensors.sensor_type, car_sessions.car_id, cars.vin").
			Joins("JOIN sensors ON measurements.sensor_id = sensors.id").
			Joins("JOIN car_sessions ON measurements.car_session_id = car_sessions.id").
			Joins("JOIN cars ON car_sessions.car_id = cars.id").
			Where("measurements.car_session_id IN ?", carSessionIDs).
			Order("measurements.created_at").
			Find(&measurements)
	} else {
		// For sampling, we need to use a more complex query
		// Get a subset of measurements by using row_number() window function
		// to select every Nth row, grouped by sensor to retain data quality across all sensors
		// Using a Common Table Expression to get row numbers

		sampleQuery := `
			WITH numbered_measurements AS (
				SELECT 
					measurements.*,
					sensors.sensor_type,
					car_sessions.car_id,
					cars.vin,
					ROW_NUMBER() OVER (PARTITION BY measurements.sensor_id ORDER BY measurements.created_at) as row_num
				FROM 
					measurements
				JOIN 
					sensors ON measurements.sensor_id = sensors.id
				JOIN 
					car_sessions ON measurements.car_session_id = car_sessions.id
				JOIN
					cars ON car_sessions.car_id = cars.id
				WHERE 
					measurements.car_session_id IN (?)
			)
			SELECT * FROM numbered_measurements
			WHERE row_num % ? = 0
			ORDER BY created_at;
		`

		db.Raw(sampleQuery, carSessionIDs, samplingFactor).Scan(&measurements)
	}

	var measurementData []MeasurementData
	for _, m := range measurements {
		measurementData = append(measurementData, MeasurementData{
			CarSessionID: m.CarSessionID,
			SensorType:   m.SensorType,
			SensorID:     m.SensorID,
			CarID:        m.CarID,
			VIN:          m.VIN,
			Data1:        float64(m.Data1),
			Data2:        float64(m.Data2),
			CreatedAt:    m.CreatedAt,
		})
	}

	return wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "measurements",
		Body: map[string]interface{}{
			"measurements": measurementData,
			"count":        len(measurementData),
			"total_count":  totalCount,
			"sampling":     sampleInfo,
		},
	}
}
func (w WsSessionController) Delete(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		ID uint `json:"id"`
	}
	var Req wsservice.WsRequestPrepared[Body]

	Req.Parse(req)

	db := database.GetDB()

	var session models.Session
	db.First(&session, Req.Body.ID)

	db.Delete(&session)

	return wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "delete",
		Body:      session,
	}
}
