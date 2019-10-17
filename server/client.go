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
	countRequests  int
}

func (c *client) UpdateState() {
	c.incCountRequest()
	c.updateLastActiveTime()
}

func (c *client) UpdateToken(token string) {
	c.token = token
}

func (c *client) updateLastActiveTime() {
	c.lastActiveTime = time.Now().Unix()
}

func (c *client) incCountRequest() {
	c.countRequests++
}

func (c client) GetLastActiveTimeAsString() string {
	return strconv.FormatInt(c.lastActiveTime, 10)
}

func (c *client) GetCountRequestsAsString() string {
	return strconv.Itoa(c.countRequests)
}

func (c client) GetAddress() *net.UDPAddr {
	return c.addr
}
