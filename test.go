package main

import (
	"encoding/json"
	"github.com/tidwall/evio"
	"udp/realmetric"
)

type Payload map[string]interface{}

func main() {
	var events evio.Events
	handler := realmetric.Start()

	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		payload := Payload{}
		payload["status"] = "ok"
		payload["time"], payload["error"] = handler(in)
		out, _ = json.Marshal(payload)

		return
	}
	if err := evio.Serve(events, "udp://127.0.0.1:3030"); err != nil {
		panic(err.Error())
	}
}
