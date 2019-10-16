package server

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

type poolClients struct {
	list   map[string]client
	isInit bool
}

func (pool poolClients) Init() error {
	if pool.isInit != true {
		ticker := time.NewTicker(time.Minute * 10)
		go func() {
			for t := range ticker.C {
				countRemoved := pool.clean()
				fmt.Println("Pool cleaned: ", countRemoved, ", time ", t)
			}
		}()
		pool.isInit = true
	}
	return nil
}

func (pool poolClients) clean() int {
	timeUnix := time.Now().Unix()
	countCleaned := 0
	for key, client := range pool.list {
		diff := timeUnix - client.lastActiveTime

		if diff > 500 {
			err := pool.RemoveByAddr(key)
			if err == nil {
				countCleaned++
			}
		}
	}
	return countCleaned
}

func (pool poolClients) HasClient(addrStr string) bool {
	_, ok := pool.list[addrStr]
	return ok
}

func (pool poolClients) AddClient(c client) bool {
	addr := c.addr.String()
	has := pool.HasClient(addr)

	if !has {
		pool.list[addr] = c
		fmt.Println("Added new client to pool")
	}

	return has
}

func (pool poolClients) GetClient(addrStr string) client {
	if !pool.HasClient(addrStr) {
		panic("Client not found")
	}
	return pool.list[addrStr]
}

func (pool poolClients) RemoveByAddr(addrStr string) error {
	if pool.HasClient(addrStr) {
		delete(pool.list, addrStr)
		return nil
	}
	return errors.New("Client by this address is not found")
}

func (pool poolClients) RemoveByClient(client client) error {
	return pool.RemoveByAddr(client.addr.String())
}
