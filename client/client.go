package main

import (
	"fmt"
	"log"
	"net"
)

type UDPWrite struct {
	Code    string
	Message string
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 1024)

	for i := 0; i < 10000; i++ {
		_, err = conn.Write([]byte(fmt.Sprintf("message no: %d", i)))
		if err != nil {
			log.Fatal(err)
		}
		n, err := conn.Read(buffer)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("%s!\n", buffer[:n])
	}
}
