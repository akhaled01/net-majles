package funcs

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var mu sync.Mutex

type pool struct {
	conns []net.Conn
}

func (p *pool) Reply(reply []byte) {
	for _, c := range p.conns {
		c.Write(reply)
	}
}

var mainpool pool

func AddUser(listen net.Listener) {
	mu.Lock()
	conn, err := listen.Accept()
	mu.Unlock()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	mainpool.conns = append(mainpool.conns, conn)
	channel := make(chan string, 1)
	go HandleRequest(conn, channel)
	for {
		chans := <-channel
		fmt.Println(chans)
		go mainpool.Reply([]byte(chans))
	}
}

func HandleRequest(conn net.Conn, channel chan string) {
	for {
		// incoming request
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}

		// write data to response
		t := time.Now().Format(time.ANSIC)
		responseStr := fmt.Sprintf("[%v] %v", t, string(buffer[:]))

		channel <- responseStr
	}
}
