package server

import (
	"encoding/json"
	"net"
)

type request struct {
	addr    *net.UDPAddr
	payload []byte
}

type payload interface{}

func (req request) GetPayload() (payload, error) {
	var data payload
	err := json.Unmarshal(req.payload, &data)

	return data, err
}
