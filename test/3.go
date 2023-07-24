package main

import (
    "fmt"
    "log"
    "net"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal("Listen error:", err)
	}

	address := listener.Addr()
	fmt.Println("Listening on:", address.String())
}