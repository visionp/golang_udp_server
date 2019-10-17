package server

import (
	"fmt"
)

type defaultHandler struct {
}

func (h defaultHandler) Handle(req Request, client *Client) Payload {
	payload := make(Payload)
	data, _ := req.GetPayload()

	pid, _ := data["pid"]
	token := fmt.Sprintf("%v", pid)
	payload["pid"] = token

	if client.token != payload["pid"] {
		fmt.Printf("Diff tokens %s != %s \n", client.token, payload["pid"])
	}

	ct := client.GetCountRequestsAsString()
	la := client.GetLastActiveTimeAsString()
	tk := client.token
	fmt.Printf("Received: %s -> %s Requests: %s  \n", req.GetPayloadAsString(), tk, ct)

	payload["last_active"] = la
	payload["token"] = tk
	payload["count_Requests"] = ct
	payload["addr"] = client.GetAddress().String()

	return payload
}
