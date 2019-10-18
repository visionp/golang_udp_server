package server

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

type PoolClients struct {
	list map[string]*Client
}

func (pool PoolClients) clean() int {
	timeUnix := time.Now().Unix()
	countCleaned := 0
	for key, Client := range pool.list {
		diff := timeUnix - Client.getLastActiveTime()

		if diff > 60 {
			err := pool.RemoveByAddr(key)
			if err == nil {
				countCleaned++
			}
		}
	}
	return countCleaned
}

func (pool PoolClients) HasClient(addrStr string) bool {
	_, ok := pool.list[addrStr]
	return ok
}

func (pool PoolClients) AddClient(c *Client) bool {
	addr := c.addr.String()
	has := pool.HasClient(addr)

	if !has {
		pool.list[addr] = c
		fmt.Printf("Added new Client to pool %s \n", c.token)
	}

	return has
}

func (pool PoolClients) GetClient(addrStr string) *Client {
	if !pool.HasClient(addrStr) {
		panic("Client not found")
	}
	return pool.list[addrStr]
}

func (pool PoolClients) RemoveByAddr(addrStr string) error {
	if pool.HasClient(addrStr) {
		delete(pool.list, addrStr)
		fmt.Printf("Remove Client: %s \n", addrStr)
		return nil
	}

	return errors.New("Client not found")
}

func (pool PoolClients) RemoveByClient(Client *Client) error {
	return pool.RemoveByAddr(Client.addr.String())
}
