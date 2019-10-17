package main

import (
	"./server"
	"fmt"
	"sync"
)

type handlerPing struct {
	mutex *sync.Mutex
}

func (h handlerPing) Handle(req server.Request, clients *server.PoolClients) server.Payload {
	var cl *server.Client

	payload := make(server.Payload)
	data, _ := req.GetPayload()

	action, ok := data["action"]

	pid, _ := data["pid"]
	token := fmt.Sprintf("%v", pid)
	payload["pid"] = token

	h.mutex.Lock()
	hasClient := clients.HasClient(req.GetAddrAsString())
	h.mutex.Unlock()

	if ok && action == "handshake" {
		if hasClient {
			h.mutex.Lock()
			err := clients.RemoveByAddr(req.GetAddrAsString())
			h.mutex.Unlock()
			if err == nil {
				fmt.Printf("Remove client: %s \n", req.GetAddrAsString())
			}
		}
	}

	h.mutex.Lock()
	if clients.HasClient(req.GetAddrAsString()) {
		cl = clients.GetClient(req.GetAddrAsString())
	} else {
		cl = &server.Client{req.GetAddr(), token, 0, 0}
		clients.AddClient(cl)
	}
	cl.UpdateState()
	h.mutex.Unlock()

	h.mutex.Lock()
	ct := cl.GetCountRequestsAsString()
	la := cl.GetLastActiveTimeAsString()
	tk := cl.GetToken()
	h.mutex.Unlock()

	if tk != payload["pid"] {
		fmt.Printf("Diff tokens %s != %s \n", cl.GetToken(), payload["pid"])
	}
	fmt.Printf("Received: %s -> %s Requests: %s  \n", req.GetPayloadAsString(), tk, ct)

	payload["last_active"] = la
	payload["token"] = tk
	payload["count_requests"] = ct
	payload["addr"] = cl.GetAddress().String()

	return payload
}
