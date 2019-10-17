package server

import (
	"fmt"
	"sync"
)

type handlerInterface interface {
	Handle(req request, clients *poolClients) response
}

type handler struct {
	mutex *sync.Mutex
}

func (h handler) Handle(req request, clients *poolClients) response {
	var cl *client

	payload := make(map[string]string)
	data, _ := req.GetPayload()

	action, ok := data["action"]

	pid, _ := data["pid"]
	token := fmt.Sprintf("%v", pid)
	payload["pid"] = token

	h.mutex.Lock()
	hasClient := clients.HasClient(req.addr.String())
	h.mutex.Unlock()

	if ok && action == "handshake" {
		if hasClient {
			h.mutex.Lock()
			err := clients.RemoveByAddr(req.addr.String())
			h.mutex.Unlock()
			if err == nil {
				hasClient = false
				fmt.Printf("Remove client: %s \n", req.addr.String())
			}
		}
	}

	h.mutex.Lock()
	if clients.HasClient(req.addr.String()) {
		cl = clients.GetClient(req.addr.String())
	} else {
		cl = &client{req.addr, token, 0, 0}
		clients.AddClient(cl)
	}
	cl.UpdateState()
	h.mutex.Unlock()

	if cl.token != payload["pid"] {
		fmt.Printf("Diff tokens %s != %s \n", cl.token, payload["pid"])
	}

	h.mutex.Lock()
	ct := cl.GetCountRequestsAsString()
	la := cl.GetLastActiveTimeAsString()
	tk := cl.token
	h.mutex.Unlock()
	fmt.Printf("Received: %s -> %s Requests: %s  \n", req.GetPayloadAsString(), tk, ct)

	payload["last_active"] = la
	payload["token"] = tk
	payload["count_requests"] = ct
	payload["addr"] = cl.GetAddress().String()

	return response{req.addr, payload}
}
