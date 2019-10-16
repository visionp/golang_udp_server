package server

import (
	"fmt"
	"net"
	"os"
)

func Listen() {
	fmt.Println("Start listen")
	handlerFunc := handler{}

	list := make(map[string]client)
	poolClients := poolClients{list, false}
	err := poolClients.Init()
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

	buffer := make([]byte, 1024)

	requestCh := make(chan request, 200)
	responseCh := make(chan response, 200)
	disp := dispatcher{requestCh, responseCh, handlerFunc, poolClients}

	defer func() {
		close(requestCh)
		close(responseCh)
	}()

	for w := 0; w < 100; w++ {
		go disp.Dispatch()
		go func() {
			for {
				res := <-responseCh
				_, err = connection.WriteToUDP(res.GetPayload(), res.addr)
				if err != nil {
					fmt.Println(err)
				}
			}
		}()
	}

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		req := request{addr, buffer[:n]}

		requestCh <- req
	}

}
