package server

import (
	"encoding/json"
	"net"
)

type Response struct {
	addr    *net.UDPAddr
	payload Payload
}

func (res Response) GetPayload() []byte {
	jsonData, err := json.Marshal(res.payload)
	if err != nil {
		return []byte(err.Error())
	}

	return jsonData
}
