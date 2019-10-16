package server

import (
	"fmt"
)

type handlerInterface interface {
	Handle(req request, clients poolClients) response
}

type handler struct {
}

func (h handler) Handle(req request, clients poolClients) response {
	var cl client

	payload := make(map[string]string)
	data, _ := req.GetPayload()
	action, ok := data["action"]

	if ok && action == "ping" {
		pid, _ := data["pid"]
		payload["pid"] = fmt.Sprintf("%v", pid)
	}

	if clients.HasClient(req.addr.String()) {
		cl = clients.GetClient(req.addr.String())
	} else {
		token, _ := data["pid"]

		cl = client{req.addr, fmt.Sprintf("%v", token), 0, 0}
		clients.AddClient(cl)
	}

	cl.UpdateState()

	if cl.token != payload["pid"] {
		fmt.Printf("%s != %s \n", cl.token, payload["pid"])
	}

	payload["last_active"] = cl.GetLastActiveTimeAsString()
	payload["token"] = cl.token
	payload["addr"] = cl.GetAddress().String()

	return response{req.addr, payload}
}
