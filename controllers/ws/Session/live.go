package ws_session_namespace

import (
	"recofiit/models"
	"recofiit/services/database"
	wsservice "recofiit/services/wsService"
	"strconv"
)

func (w WsSessionController) SendLiveSessionData(data interface{}) {

	db := database.GetDB()
	var session models.Session
	result := db.Where("started_at is not null").Where("ended_at is null").First(&session)

	if result.Error != nil {
		panic("Failed to find started session")
	}
	
	var carSessionID = session.ID

	endpointResponse := wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "live/" + strconv.Itoa(int(carSessionID)),
		Body:      data,
	}

	res := endpointResponse.ToJSON()

	wsservice.Manager.Broadcast <- res
}
