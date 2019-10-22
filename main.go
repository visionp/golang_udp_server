package main

import (
	"fmt"
	"udp/realmetric"
	"udp/server"
)

type handlerPing struct {
}

func (h handlerPing) Handle(req server.Request, client *server.Client) server.Payload {
	fmt.Println("Start handle")
	payload := make(server.Payload)
	fmt.Println("GetPayload")
	data, _ := req.GetPayload()

	payload["action"] = "pong"
	payload["pid"] = fmt.Sprintf("%v", data["pid"])
	fmt.Println("GetToken")
	//payload["token"] = client.GetToken()
	fmt.Println("GetCountRequestsAsString")
	payload["count_requests"] = client.GetCountRequestsAsString()
	fmt.Println("Before return")

	return payload
}

func main() {
	handler := realmetric.Start()
	m := make(map[string]server.HandlerInterface)
	handlersCollection := &server.HandlersCollection{Handlers: m}

	handlersCollection.Add("handshake", handlerPing{})
	handlersCollection.Add("ping", handlerPing{})
	handlersCollection.Add("track", realMetricHandler{handler})
	handlersCollection.Add("undefined", realMetricHandler{handler})

	var serverUdp = server.Server{Handlers: handlersCollection}
	serverUdp.Start(":3030")
}

type realMetricHandler struct {
	realMetricTrack func(body []byte) (int64, error)
}

func (h realMetricHandler) Handle(req server.Request, client *server.Client) server.Payload {
	payload := make(server.Payload)

	//payload["_timing"], payload["error"] = h.realMetricTrack(req.Payload)

	payload["action"] = "track"
	payload["token"] = client.GetToken()
	payload["count_requests"] = client.GetCountRequestsAsString()

	return payload
}
