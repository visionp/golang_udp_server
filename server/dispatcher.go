package server

import (
	"sync"
)

type dispatcher struct {
	mutex      *sync.Mutex
	handlers   *HandlersCollection
	responseCh chan Response
}

func (dis dispatcher) Dispatch(request Request, client *Client) {
	dis.debug("Start dispatch")
	dis.responseCh <- dis.resolveRequest(request, client)
	dis.debug("Dispatched")
}

func (dis dispatcher) resolveRequest(req Request, client *Client) Response {
	dis.debug("Start resolveRequest")
	actionName := req.GetAction()

	if actionName != "ping" {
		dis.debug("Update client state")
		dis.mutex.Lock()
		client.updateState()
		dis.debug("Client state updated")
		dis.mutex.Unlock()
	}

	payload := Payload{}

	dis.debug("Get handler")
	handler := dis.handlers.GetHandler(actionName)
	dis.debug("Start handle: " + actionName)
	payload = handler.Handle(req, client)
	dis.debug("Start handled")

	return Response{req.GetAddr(), payload}
}

func (dis dispatcher) debug(mess string) {
	//fmt.Println(mess)
}
