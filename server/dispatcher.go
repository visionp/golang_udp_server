package server

type dispatcher struct {
	requestCh   chan request
	responseCh  chan response
	handler     handlerInterface
	poolClients poolClients
}

func (dis dispatcher) Dispatch() {
	for {
		r := <-dis.requestCh
		dis.responseCh <- dis.handler.Handle(r, dis.poolClients)
	}
}
