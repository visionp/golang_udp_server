package server

import "net"

type client struct {
	addr *net.UDPAddr
	token string
	lastActiveTime int64
	lastResponseId int
}