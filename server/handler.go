package server

import (
	"fmt"
)

type handlerInterface interface {
	Handle(req request, clients *poolClients) response
}

type handler struct {
}

func (h handler) Handle(req request, clients *poolClients) response {
	var cl *client

	payload := make(map[string]string)
	data, _ := req.GetPayload()
	action, ok := data["action"]

	pid, _ := data["pid"]
	token := fmt.Sprintf("%v", pid)
	payload["pid"] = token

	if ok && action == "handshake" {
		if clients.HasClient(req.addr.String()) {
			err := clients.RemoveByAddr(req.addr.String())
			if err == nil {
				fmt.Printf("Remove client: %s \n", req.addr.String())
			}
		}
	}

	if clients.HasClient(req.addr.String()) {
		cl = clients.GetClient(req.addr.String())
	} else {
		cl = &client{req.addr, token, 0, 0}
		clients.AddClient(cl)
	}

	cl.UpdateState()

	if cl.token != payload["pid"] {
		fmt.Printf("Diff tokens %s != %s \n", cl.token, payload["pid"])
	}

	payload["last_active"] = cl.GetLastActiveTimeAsString()
	payload["token"] = cl.token
	payload["count_requests"] = cl.GetCountRequestsAsString()
	payload["addr"] = cl.GetAddress().String()

	return response{req.addr, payload}
}
