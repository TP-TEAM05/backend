package ws_session_namespace

import (
	"recofiit/models"
	"recofiit/services/database"
	"recofiit/services/redis"
	"recofiit/services/statistics"
	wsservice "recofiit/services/wsService"
	"strconv"

	api "github.com/ReCoFIIT/integration-api"
	"github.com/getsentry/sentry-go"
)

type ExtendedUpdateVehicleDatagram struct {
	api.BaseDatagram
	Vehicle api.UpdateVehicleVehicle `json:"vehicle"`
	Network statistics.NetworkStats  `json:"network"`
}

func (w WsSessionController) SendLiveSessionData(data *api.UpdateVehicleDatagram) {

	db := database.GetDB()
	var session models.Session
	result := db.Where("started_at is not null").Where("ended_at is null").First(&session)

	if result.Error != nil {
		// Session not found
		sentry.CaptureMessage("Session not found, live data cannot be send for vehicle: " + data.Vehicle.Vin)
		return
	}

	var dataExtended = &ExtendedUpdateVehicleDatagram{
		BaseDatagram: api.BaseDatagram{
			Index:     data.Index,
			Type:      data.Type,
			Timestamp: data.Timestamp,
		},
		Vehicle: data.Vehicle,
		Network: *redis.GetNetworkStats(data.Vehicle.Vin),
	}

	var carSessionID = session.ID

	endpointResponse := wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "live/" + strconv.Itoa(int(carSessionID)),
		Body:      dataExtended,
	}

	res := endpointResponse.ToJSON()

	wsservice.Manager.Broadcast <- res
}
