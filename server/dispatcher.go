package server

import (
	"fmt"
	"sync"
)

type dispatcher struct {
	handlers    *HandlersCollection
	poolClients *PoolClients
	mutex       *sync.Mutex
}

func (dis dispatcher) Dispatch(RequestCh chan Request, ResponseCh chan Response) {
	for {
		req := <-RequestCh
		ResponseCh <- dis.resolveRequest(req)
	}
}

func (dis dispatcher) resolveRequest(req Request) Response {
	actionName := req.GetAction()
	isHandShake := actionName == "handshake"

	if isHandShake {
		dis.handShake(req)
	}

	client := dis.resolveClient(actionName, req)
	payload := Payload{}

	if isHandShake {
		payload["token"] = client.GetToken()
	} else {
		handler := dis.getAction(actionName)
		payload = handler.Handle(req, client)
	}

	return Response{req.GetAddr(), payload}
}

func (dis dispatcher) getAction(action string) HandlerInterface {
	return dis.handlers.GetHandler(action)
}

func (dis dispatcher) handShake(req Request) {
	dis.mutex.Lock()
	if dis.poolClients.HasClient(req.addr.String()) {
		dis.mutex.Lock()
		err := dis.poolClients.RemoveByAddr(req.addr.String())
		dis.mutex.Unlock()
		if err == nil {
			fmt.Printf("Remove Client: %s \n", req.addr.String())
		}
	}
	dis.mutex.Unlock()
}

func (dis dispatcher) resolveClient(action string, req Request) *Client {
	var cl *Client

	dis.mutex.Lock()
	if dis.poolClients.HasClient(req.addr.String()) {
		cl = dis.poolClients.GetClient(req.addr.String())
	} else {
		data, _ := req.GetPayload()
		pid, _ := data["pid"]
		token := fmt.Sprintf("%v", pid)
		cl = &Client{req.addr, token, 0, 0}
		dis.poolClients.AddClient(cl)
	}
	cl.UpdateState()
	dis.mutex.Unlock()

	return cl
}
