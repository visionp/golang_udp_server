package server

type HandlerInterface interface {
	Handle(req Request, client *Client) Payload
}

type HandlersCollection struct {
	handlers map[string]HandlerInterface
}

func (collection HandlersCollection) GetHandler(action string) HandlerInterface {
	handler, ok := collection.handlers[action]

	if ok {
		return handler
	}

	return defaultHandler{}
}

func (collection HandlersCollection) Add(action string, handler HandlerInterface) {
	collection.handlers[action] = handler
}
