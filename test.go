package main

import (
	"encoding/json"
	"github.com/tidwall/evio"
	"udp/realmetric"
	"udp/server"
)

func main() {
	var events evio.Events
	handler := realmetric.Start()
	events.NumLoops = 36000
	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		payload := server.Payload{}
		payload["status"] = "ok"
		payload["time"], payload["error"] = handler(in)
		out, _ = json.Marshal(payload)

		return
	}
	if err := evio.Serve(events, "udp://127.0.0.1:3050"); err != nil {
		panic(err.Error())
	}
}
