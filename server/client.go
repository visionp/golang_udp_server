package server

import (
	"fmt"
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

func (c client) UpdateLastActiveTime() {
	c.lastActiveTime = time.Now().Unix()
	fmt.Println("Updated last active time " + c.addr.String() + " Set: " + c.GetLastActiveTimeAsString())
}

func (c client) GetLastActiveTimeAsString() string {
	return strconv.FormatInt(c.lastActiveTime, 10)
}
