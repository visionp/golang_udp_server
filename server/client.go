package server

import (
	"net"
	"strconv"
	"sync/atomic"
	"time"
)

type Client struct {
	addr           *net.UDPAddr
	token          string
	lastActiveTime int64
	countRequests  int64
}

func (c Client) GetLastActiveTimeAsString() string {
	return strconv.FormatInt(c.getLastActiveTime(), 10)
}

func (c *Client) GetCountRequestsAsString() string {
	return strconv.FormatInt(c.getCountRequests(), 10)
}

func (c Client) GetAddress() *net.UDPAddr {
	return c.addr
}

func (c Client) GetToken() string {
	t := c.token
	return t
}

func (c *Client) updateState() {
	c.countRequests += 1
	c.lastActiveTime = time.Now().Unix()
}

func (c *Client) updateToken(token string) {
	c.token = token
}

func (c *Client) getCountRequests() int64 {
	return atomic.LoadInt64(&c.countRequests)
}

func (c *Client) getLastActiveTime() int64 {
	return atomic.LoadInt64(&c.lastActiveTime)
}

func (c Client) isValidToken(token string) bool {
	return c.token == token
}
