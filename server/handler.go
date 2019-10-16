package server

import (
	"strconv"
	"time"
)

type handlerInterface interface {
	Handle(req request, clients poolClients) response
}

type handler struct {
}

func (h handler) Handle(req request, clients poolClients) response {
	var cl client
	if clients.HasClient(req.addr.String()) {
		cl = clients.GetClient(req.addr.String())
		cl.UpdateLastActiveTime()
	} else {
		token := "token_" + req.addr.String()
		cl = client{req.addr, token, time.Now().Unix(), 1}
		clients.AddClient(cl)
	}

	payload := make(map[string]string)
	payload["token"] = cl.token
	payload["last_active_time"] = cl.GetLastActiveTimeAsString()
	payload["time"] = strconv.FormatInt(time.Now().Unix(), 10)
	res := response{req.addr, payload}

	return res
}
