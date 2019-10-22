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

	mutex := &sync.Mutex{}
	list := make(map[string]*Client)
	poolClients := &PoolClients{mutex: mutex, list: list}

	disp := &dispatcher{mutex, server.Handlers}

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

	for {
		buffer := make([]byte, 65536)
		n, addr, err := connection.ReadFromUDP(buffer)

		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		go func(pc *net.UDPConn, address *net.UDPAddr, in []byte, pool *PoolClients, dispatcher *dispatcher, m *sync.Mutex) {
			request := Request{address, in}
			m.Lock()
			client := pool.resolveClient(request)
			m.Unlock()
			res := disp.Dispatch(request, client)
			_, _ = pc.WriteToUDP(res.GetPayload(), res.addr)
		}(connection, addr, buffer[:n], poolClients, disp, mutex)

		buffer = nil
	}

}
