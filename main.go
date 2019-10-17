package main

import (
	"fmt"
	"udp/server"
)

type handlerPing struct {
}

func (h handlerPing) Handle(req server.Request, client *server.Client) server.Payload {
	payload := make(server.Payload)
	data, _ := req.GetPayload()

	//fmt.Printf("Receive: %s \n\r", req.GetPayloadAsString())

	payload["action"] = "pong"
	payload["rand"] = fmt.Sprintf("%v", data["rand"])
	payload["pid"] = fmt.Sprintf("%v", data["pid"])
	payload["token"] = client.GetToken()
	payload["count_requests"] = client.GetCountRequestsAsString()

	return payload
}

func main() {
	m := make(map[string]server.HandlerInterface)
	handlersCollection := &server.HandlersCollection{Handlers: m}

	handlersCollection.Add("ping", handlerPing{})
	handlersCollection.Add("action", handlerPing{})
	var serverUdp = server.Server{Handlers: handlersCollection}
	serverUdp.Start(":3030")
}
