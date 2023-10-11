package main

import (
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

type UDPRead struct {
	Addr   net.Addr
	Buffer []byte
}

func main() {
	conn, err := net.ListenPacket("udp", ":8080")
	if conn != nil {
		defer conn.Close()
	}
	if err != nil {
		log.Fatal(err)
	}

	jobs := make(chan UDPRead)
	var noOfWorkers int32 = 4

	for i := 0; i < int(noOfWorkers); i++ {
		go func(i int) {
			defer func() {
				if atomic.AddInt32(&noOfWorkers, -1) == 0 {
					close(jobs)
				}
			}()
			response(conn, i, jobs)
		}(i)
	}

	for {
		log.Printf("reading from UDP server...\n")
		buffer := make([]byte, 1024)
		_, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			log.Println(err)
			continue
		}
		jobs <- UDPRead{
			Addr:   addr,
			Buffer: buffer,
		}
	}
}

func response(conn net.PacketConn, wID int, jobs chan UDPRead) {
	for r := range jobs {
		res := fmt.Sprintf("worker id: %d. Your message: %s!", wID, r.Buffer)
		log.Printf("%s\n", res)
		conn.WriteTo([]byte(res), r.Addr)
	}
}
