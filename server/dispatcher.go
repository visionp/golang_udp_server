package server

import "fmt"

type dispatcher struct {
	requestCh   chan request
	responseCh  chan response
	handler     handlerInterface
	poolClients poolClients
}

func (dis dispatcher) Dispatch() {
	for {
		r := <-dis.requestCh
		fmt.Println("Dispatch")
		dis.responseCh <- dis.handler.Handle(r)
	}
}
