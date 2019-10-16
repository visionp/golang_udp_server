package server

import (
	"encoding/json"
	"net"
)

type response struct {
	addr   *net.UDPAddr
	id     int64
	status string
}

func (res response) GetPayload() []byte {
	data := make(map[string]string)
	data["status"] = res.status

	jsonData, err := json.Marshal(data)
	if err != nil {
		return []byte(err.Error())
	}

	return jsonData
}
