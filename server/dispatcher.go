package server

type dispatcher struct {
	handler     handlerInterface
	poolClients *poolClients
}

func (dis dispatcher) Dispatch(requestCh chan request, responseCh chan response) {
	for {
		r := <-requestCh
		responseCh <- dis.handler.Handle(r, dis.poolClients)
	}
}
