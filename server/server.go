package server

import (
	"fmt"
	"net"
	"os"
	"sync"
)

type Server struct {
	Handlers *HandlersCollection
}

func (server Server) Start(port string) {
	fmt.Printf("Server listening on port %s \n", port)

	responseCh := make(chan Response, 4096)

	mutex := &sync.Mutex{}
	list := make(map[string]*Client)
	poolClients := &PoolClients{mutex: mutex, list: list}

	disp := dispatcher{mutex, server.Handlers, responseCh}

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
			return
		}
	}()

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
		buffer := make([]byte, 4096)
		n, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		go func(address *net.UDPAddr, in []byte, pool *PoolClients) {
			request := Request{addr, buffer[:n]}
			client := poolClients.resolveClient(request)
			disp.Dispatch(request, client)
		}(addr, buffer, poolClients)
	}

}
