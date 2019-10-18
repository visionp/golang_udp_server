package main

import "github.com/tidwall/evio"

func main() {
	var events evio.Events
	events.NumLoops = 36000
	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		out = in
		return
	}
	if err := evio.Serve(events, "udp://127.0.0.1:3030"); err != nil {
		panic(err.Error())
	}
}
