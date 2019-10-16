package server

type handlerInterface interface {
	Handle(req request) response
}

type handler struct {
}

func (h handler) Handle(req request) response {
	res := response{req.addr, 1, "ok"}

	return res
}
