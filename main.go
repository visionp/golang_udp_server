package main

import (
	"udp/server"
)

func main() {
	handlersCollection := &server.HandlersCollection{}

	var serverUdp = server.Server{Handlers: handlersCollection}
	serverUdp.Start(":3030")
}
