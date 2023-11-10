package wsservice

import (
	"encoding/json"
	"fmt"
)

type WsRequest struct {
	Namespace string      `json:"namespace"`
	Endpoint  string      `json:"endpoint"`
	Body      interface{} `json:"body"`
}

type WsResponse struct {
	Namespace string      `json:"namespace"`
	Endpoint  string      `json:"endpoint"`
	Body      interface{} `json:"body"`
	Error     string      `json:"error"`
}

type Endpoint = map[string]func()

type Namespace = map[string]Endpoint

type WsRouter struct {
	Routes Namespace
}

func (m *WsRequest) ParseJSON(message []byte) {

	err := json.Unmarshal([]byte(message), &m)

	if err != nil {
		fmt.Println("error:", err)
	}
}

func (r *WsRouter) Register(namespace string, endpoint string, handler func()) {

	if r.Routes[namespace] == nil {
		r.Routes[namespace] = make(Endpoint)
	}

	r.Routes[namespace][endpoint] = handler
}

func (r *WsRouter) GetHandler(namespace string, endpoint string) func() {

	if r.Routes[namespace] == nil {
		return nil
	}

	return r.Routes[namespace][endpoint]
}

func (r *WsRouter) Error() []byte {
	res := WsResponse{
		Namespace: "error",
		Endpoint:  "error",
		Body:      nil,
		Error:     "No handler found for this endpoint",
	}

	resJSON, _ := json.Marshal(&res)

	return resJSON
}
