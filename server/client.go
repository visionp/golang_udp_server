package server

import (
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type Client struct {
	sync.Mutex
	addr           *net.UDPAddr
	token          string
	lastActiveTime int64
	countRequests  int64
}

func (c Client) GetLastActiveTimeAsString() string {
	c.Lock()
	c.Unlock()
	return strconv.FormatInt(c.getLastActiveTime(), 10)
}

func (c *Client) GetCountRequestsAsString() string {
	c.Lock()
	c.Unlock()
	return strconv.FormatInt(c.getCountRequests(), 10)
}

func (c Client) GetAddress() *net.UDPAddr {
	c.Lock()
	c.Unlock()
	return c.addr
}

func (c Client) GetToken() string {
	c.Lock()
	defer c.Unlock()
	return c.token
}

func (c *Client) updateState() {
	c.Lock()
	defer c.Unlock()
	atomic.AddInt64(&c.countRequests, 1)
	atomic.StoreInt64(&c.lastActiveTime, time.Now().Unix())
}

func (c *Client) updateToken(token string) {
	c.Lock()
	c.Unlock()
	c.token = token
}

func (c *Client) getCountRequests() int64 {
	return atomic.LoadInt64(&c.countRequests)
}

func (c *Client) getLastActiveTime() int64 {
	return atomic.LoadInt64(&c.lastActiveTime)
}

func (c Client) isValidToken(token string) bool {
	c.Lock()
	c.Unlock()
	return c.token == token
}
