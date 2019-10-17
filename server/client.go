package server

import (
	"net"
	"strconv"
	"time"
)

type Client struct {
	addr           *net.UDPAddr
	token          string
	lastActiveTime int64
	countRequests  int
}

func (c *Client) UpdateState() {
	c.incCountRequest()
	c.updateLastActiveTime()
}

func (c *Client) UpdateToken(token string) {
	c.token = token
}

func (c *Client) updateLastActiveTime() {
	c.lastActiveTime = time.Now().Unix()
}

func (c *Client) incCountRequest() {
	c.countRequests++
}

func (c Client) GetLastActiveTimeAsString() string {
	return strconv.FormatInt(c.lastActiveTime, 10)
}

func (c *Client) GetCountRequestsAsString() string {
	return strconv.Itoa(c.countRequests)
}

func (c Client) GetAddress() *net.UDPAddr {
	return c.addr
}

func (c Client) GetToken() string {
	return c.token
}

func (c Client) isValidToken(token string) bool {
	return c.token == token
}
