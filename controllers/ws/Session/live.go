package ws_session_namespace

import (
	wsservice "recofiit/services/wsService"
)

func (w WsSessionController) SendLiveSessionData(data interface{}) {

	// empty object interface

	endpointResponse := wsservice.WsResponse[interface{}]{
		Namespace: "session",
		Endpoint:  "live/0",
		Body:      data,
	}

	res := endpointResponse.ToJSON()

	wsservice.Manager.Broadcast <- res
}
