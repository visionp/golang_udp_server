package server

import (
	"fmt"
	"net"
	"os"
	"sync"
)

func Listen() {
	list := make(map[string]*client)
	mutex := &sync.Mutex{}
	handlerFunc := handler{mutex}

	poolClients := &poolClients{list, false}
	err := poolClients.Init(mutex)
	if err != nil {
		fmt.Println(err)
		return
	}

	PORT := ":3030"
	s, err := net.ResolveUDPAddr("udp4", PORT)
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

	requestCh := make(chan request, 1024)
	responseCh := make(chan response, 1024)
	disp := dispatcher{handlerFunc, poolClients}

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

		requestCh <- request{addr, buffer[:n]}
	}

}
