package main

import (
	"fmt"
	"net"

	"github.com/stormtrooper1859/p2p-chat/internal/stun"
)

func main() {
	addr := net.UDPAddr{
		Port: 12345,
		IP:   net.ParseIP("localhost"),
	}
	ser, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("ListenUDP error %v\n", err)
		return
	}

	ip, port := stun.GetIP(ser)
	fmt.Printf("%s:%d\n", ip.String(), port)

	ip, port = stun.GetIP(ser)
	fmt.Printf("%s:%d\n", ip.String(), port)
}
