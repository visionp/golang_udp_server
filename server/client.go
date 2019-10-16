package server

import (
	"net"
	"strconv"
	"time"
)

type client struct {
	addr           *net.UDPAddr
	token          string
	lastActiveTime int64
	lastResponseId int
}

func (c *client) UpdateLastActiveTime() {
	c.lastActiveTime = time.Now().Unix()
}

func (c client) GetLastActiveTimeAsString() string {
	return strconv.FormatInt(c.lastActiveTime, 10)
}
