package server

import (
	"fmt"
	"sync"
)

type dispatcher struct {
	handlers    *HandlersCollection
	poolClients *PoolClients
	clientsCh   chan *Client
	mutex       *sync.Mutex
}

func (dis dispatcher) Dispatch(RequestCh chan Request, ResponseCh chan Response) {
	for {
		select {
		case req := <-RequestCh:
			ResponseCh <- dis.resolveRequest(req)
		case cl := <-dis.clientsCh:
			cl.updateState()
		}
	}
}

func (dis dispatcher) resolveRequest(req Request) Response {
	actionName := req.GetAction()
	isHandShake := actionName == "handshake"

	if isHandShake {
		dis.handShake(req)
	}

	client := dis.resolveClient(actionName, req)

	if actionName != "ping" {
		dis.clientsCh <- client
	}

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
	defer dis.mutex.Unlock()

	err := dis.poolClients.RemoveByAddr(req.addr.String())
	if err == nil {
	}
}

func (dis dispatcher) resolveClient(action string, req Request) *Client {
	var cl *Client

	dis.mutex.Lock()
	defer dis.mutex.Unlock()

	if dis.poolClients.HasClient(req.addr.String()) {
		cl = dis.poolClients.GetClient(req.addr.String())
	} else {
		data, _ := req.GetPayload()
		pid, _ := data["pid"]
		token := fmt.Sprintf("%v", pid)
		cl = &Client{
			addr:           req.addr,
			token:          token,
			lastActiveTime: 0,
			countRequests:  0,
		}
		dis.poolClients.AddClient(cl)
	}

	return cl
}
