package server

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
)


type data struct{
	addr *net.UDPAddr
	id int
}

type Counter struct {
	count int
}

func (self Counter) currentValue() int {
	return self.count
}
func (self *Counter) increment() {
	self.count++
}

func h(d data, counter *Counter, mutex *sync.Mutex) data {
	mutex.Lock()
	fmt.Println("COUNTER " + strconv.Itoa(counter.currentValue()))
	counter.increment()
	mutex.Unlock()

	return data{d.addr, d.id}
}

func Listen() {
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

	defer func () {
		err = connection.Close()
		if err != nil {
			fmt.Println(err)
		}
	} ()

	//var currentId int
	//var lossPackets int
	var mutex = &sync.Mutex{}
	//m :=  make(map[string]int)
	//receivePackets :=  make(map[string]int)
	//
	//countPackets := 0
	buffer := make([]byte, 1024)

	ch := make(chan data)
	response := make(chan data)


	defer func() {
		close(ch);
		close(response);
	}()

	counter := Counter{1}

	for w := 0; w < 100; w++ {
		go func() {
			for {
				select {
				case msg := <-ch:
					response <- h(msg, &counter, mutex)
				case msg := <-response:
					d := strconv.Itoa(msg.id)
					_, err = connection.WriteToUDP([]byte("Hello -> " + d), msg.addr)
					if err != nil {
						fmt.Println(err)
					}
				}
			}
		} ()
	}

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		currentId, _ := strconv.Atoi(string(buffer[:n]))

		ch <- data{addr, currentId}
	}


}