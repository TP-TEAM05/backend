package ws_manual_control_namespace

import (
	"encoding/json"
	"fmt"
	"recofiit/services/communication"
	wsservice "recofiit/services/wsService"
	"time"

	api "github.com/TP-TEAM05/integration-api"
	"github.com/getsentry/sentry-go"
)

type WsManualControlController struct{}

// SendControlCommand sends a manual control command to a specific vehicle via car-integration
func (w WsManualControlController) SendControlCommand(req []byte) wsservice.WsResponse[interface{}] {
	type Body struct {
		Vin       string  `json:"vin"`
		Direction float32 `json:"direction"` // -1 to 1 (left to right)
		Speed     float32 `json:"speed"`     // -1 to 1 (reverse to forward)
	}

	var Req wsservice.WsRequestPrepared[Body]
	Req.Parse(req)

	// Format the control message for the ROS2 nodes
	// The message will be passed through car-integration to the vehicle's UDP listener
	// which expects: {"updateVehicleDecision": {"message": "command_string", "vin": "VIN"}}
	controlMessage := fmt.Sprintf("<%.2f,%.2f,0>", Req.Body.Direction, Req.Body.Speed)

	// Get the UDP connection manager
	conn := communication.GetCarIntegrationConnection()
	if conn == nil {
		sentry.CaptureMessage("Car integration connection not available")
		return wsservice.WsResponse[interface{}]{
			Namespace: "manual_control",
			Endpoint:  "send_command",
			Body:      map[string]string{"error": "Car integration connection not available"},
		}
	}

	// Create the decision datagram
	datagram := &api.UpdateVehicleDecisionDatagram{
		BaseDatagram: api.BaseDatagram{
			Type:      "decision_update",
			Timestamp: time.Now().UTC().Format(api.TimestampFormat),
		},
		VehicleDecision: api.UpdateVehicleDecision{
			Vin:     Req.Body.Vin,
			Message: controlMessage,
		},
	}

	// Marshal to JSON for logging
	data, err := json.Marshal(datagram)
	if err != nil {
		sentry.CaptureException(err)
		fmt.Printf("Error marshalling control command: %v\n", err)
	} else {
		fmt.Printf("[MANUAL-CONTROL] Sending command to VIN %s: %s\n", Req.Body.Vin, string(data))
	}

	// Send via UDP to car-integration
	conn.WriteDatagram(datagram, true)

	return wsservice.WsResponse[interface{}]{
		Namespace: "manual_control",
		Endpoint:  "send_command",
		Body: map[string]interface{}{
			"success": true,
			"vin":     Req.Body.Vin,
			"command": controlMessage,
		},
	}
}
