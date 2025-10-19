package wsservice

import (
	"encoding/json"
	"fmt"

	"github.com/getsentry/sentry-go"
)

type WsRequest[B any] struct {
	Namespace string `json:"namespace"`
	Endpoint  string `json:"endpoint"`
	Body      B      `json:"-"`
}

type WsRequestPrepared[B any] struct {
	Namespace string `json:"namespace"`
	Endpoint  string `json:"endpoint"`
	Body      B      `json:"body"`
	Message   []byte
}

type WsResponse[B any] struct {
	Namespace string `json:"namespace"`
	Endpoint  string `json:"endpoint"`
	Body      B      `json:"body"`
	Error     string `json:"error"`
}

// Function type with optional parameters
type Handler func(message []byte) WsResponse[interface{}]

type Endpoint map[string]Handler

type Namespace = map[string]Endpoint

type WsRouter struct {
	Routes Namespace
}

func (m *WsRequest[B]) Parse(message []byte) {

	err := json.Unmarshal([]byte(message), &m)

	if err != nil {
		sentry.CaptureException(err)
		fmt.Println("[Error] [WsRequest]:", err)
	}
}

func (r *WsRequestPrepared[B]) Parse(message []byte) {

	err := json.Unmarshal([]byte(message), &r)

	if err != nil {
		sentry.CaptureException(err)
		fmt.Println("[Error] [WsRequestPrepared]:", err)
	}
}

func (m *WsResponse[B]) ToJSON() []byte {

	resJSON, _ := json.Marshal(&m)

	return resJSON
}

func (r *WsRouter) Register(namespace string, endpoint string, handler Handler) {

	if r.Routes[namespace] == nil {
		r.Routes[namespace] = make(Endpoint)
	}

	r.Routes[namespace][endpoint] = handler
}

func (r *WsRouter) GetHandler(namespace string, endpoint string) Handler {

	if r.Routes[namespace] == nil {
		return nil
	}

	return r.Routes[namespace][endpoint]
}
