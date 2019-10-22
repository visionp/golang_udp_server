package server

import (
	"fmt"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type PoolClients struct {
	mutex *sync.Mutex
	list  map[string]*Client
}

func (pool PoolClients) init() {
	ticker := time.NewTicker(time.Second * 600)

	go func() {
		for t := range ticker.C {
			countRemoved := pool.clean()
			fmt.Printf("Pool cleaned: %d, date: %s \n\r", countRemoved, t)
		}
	}()
}

func (pool PoolClients) clean() int {
	timeUnix := time.Now().Unix()
	countCleaned := 0
	for key, Client := range pool.list {
		diff := timeUnix - Client.getLastActiveTime()

		if diff > 600 {
			pool.mutex.Lock()
			err := pool.RemoveByAddr(key)
			pool.mutex.Unlock()
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
		//fmt.Printf("Added new Client to pool %s \n", c.token)
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
	pool.mutex.Lock()
	defer pool.mutex.Unlock()
	return pool.RemoveByAddr(Client.addr.String())
}

func (pool PoolClients) resolveClient(req Request) *Client {
	var cl *Client

	if pool.HasClient(req.addr.String()) {
		cl = pool.GetClient(req.addr.String())
	} else {
		data, _ := req.GetPayload()
		pid, _ := data["pid"]
		token := fmt.Sprintf("%v", pid)
		cl = &Client{
			addr:           req.addr,
			token:          token,
			lastActiveTime: 0,
			countRequests:  0,
		}
		pool.AddClient(cl)
	}

	return cl
}
