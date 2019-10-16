package server

import (
	"encoding/json"
	"net"
)

type response struct {
	addr    *net.UDPAddr
	payload map[string]string
}

func (res response) GetPayload() []byte {
	jsonData, err := json.Marshal(res.payload)
	if err != nil {
		return []byte(err.Error())
	}

	return jsonData
}
