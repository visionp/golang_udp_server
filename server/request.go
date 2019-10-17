package server

import (
	"encoding/json"
	"fmt"
	"net"
)

type Request struct {
	addr    *net.UDPAddr
	Payload []byte
}

type Payload map[string]interface{}

func (req Request) GetPayload() (Payload, error) {
	var data Payload
	err := json.Unmarshal(req.Payload, &data)

	return data, err
}

func (req Request) GetPayloadAsString() string {
	return string(req.Payload)
}

func (req Request) GetAction() string {
	data, _ := req.GetPayload()

	action, ok := data["action"]

	if ok {
		return fmt.Sprintf("%v", action)
	}

	return "undefined"
}

func (req Request) GetAddr() *net.UDPAddr {
	return req.addr
}

func (req Request) GetAddrAsString() string {
	return req.addr.String()
}
