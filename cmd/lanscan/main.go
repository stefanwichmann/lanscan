package main

import (
	"github.com/stefanwichmann/lanscan"
	"log"
	"time"
)

func main() {
	network := "tcp4"
	port := 80
	timeout := 5 * time.Second
	threads := 20

	log.Printf("Scanning link local network for %v services on port %d.", network, port)
	start := time.Now()
	hosts, err := lanscan.ScanLinkLocal(network, port, threads, timeout)
	duration := time.Now().Sub(start)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Scan results:")
	for _, host := range hosts {
		log.Printf("Host %v responded on port %d", host, port)
	}
	log.Printf("Scan duration: %v", duration)

}
