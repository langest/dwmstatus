package main

import (
	"log"
	"net"
)

func networkConn() string {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		log.Println(err)
		return "Network: err"
	}
	if len(addr) > 2 { // We found more connections than loopback
		return "Net: conn"
	} else {
		return "Net: disc"
	}
}
