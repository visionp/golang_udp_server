package server

import (
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

type Server struct {
	Handlers *HandlersCollection
}

func (server Server) Start(port string) {
	fmt.Printf("Server listening on port %s \n", port)

	mutex := &sync.Mutex{}
	list := make(map[string]*Client)
	poolClients := &PoolClients{list: list}

	disp := dispatcher{server.Handlers, poolClients, mutex}
	requestCh := make(chan Request, 1024)
	responseCh := make(chan Response, 1024)

	ticker := time.NewTicker(time.Second * 20)

	go func() {
		for t := range ticker.C {
			mutex.Lock()
			countRemoved := poolClients.clean()
			mutex.Unlock()
			fmt.Println("Pool cleaned: ", countRemoved, ", time ", t)
		}
	}()

	s, err := net.ResolveUDPAddr("udp4", port)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer func() {
		err = connection.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	for i := 0; i < 50; i++ {
		go disp.Dispatch(requestCh, responseCh)
	}

	go func() {
		for {
			res := <-responseCh
			_, err = connection.WriteToUDP(res.GetPayload(), res.addr)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	for {
		buffer := make([]byte, 1024)
		n, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		requestCh <- Request{addr, buffer[:n]}
	}

}
