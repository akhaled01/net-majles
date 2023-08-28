package main

import (
	"fmt"
	"log"
	"nc/funcs"
	"net"
	"os"
	"time"
)

func main() {
	var PORT int
	TYPE := "tcp"
	// ConnArray := []net.Conn{}
	if len(os.Args[1:]) < 1 {
		PORT = 8989
	} else if len(os.Args[1:]) > 1 {
		log.Fatal("[USAGE]: ./server.go $PORT")
	} else {
		PORT, _ = funcs.Atoi(os.Args[1])
		_, err := funcs.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	}

	listen, err := net.Listen(TYPE, funcs.GetLocalIP()+":"+fmt.Sprint(PORT))
	// fmt.Println("Listening and serving on port ")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer listen.Close()
	for {
		go funcs.AddUser(listen)
		time.Sleep(1 * time.Microsecond)
	}
}
