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

// struct pool is how we store connections
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
	// lock the thread to accept connections, then unlock
	mu.Lock()
	conn, err := listen.Accept()
	mu.Unlock()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	mainpool.conns = append(mainpool.conns, conn)
	//TODO: Add username authentication here
	channel := make(chan string, 1)
	go HandleRequest(conn, channel)
	for message := range channel {
		fmt.Println(message)
		go mainpool.Reply([]byte(message))
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

		// write data to channel
		t := time.Now().Format(time.ANSIC)
		responseStr := fmt.Sprintf("[%v] %v", t, string(buffer[:]))

		channel <- responseStr
	}
}
