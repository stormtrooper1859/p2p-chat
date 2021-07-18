package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/stormtrooper1859/p2p-chat/internal/stun"
)

var (
	ser     *net.UDPConn
	dst     *net.UDPAddr
	chclose = make(chan struct{})
)

func main() {
	addr := net.UDPAddr{
		// IP: net.ParseIP("0.0.0.0"),
	}

	var err error
	ser, err = net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Printf("ListenUDP error %v\n", err)
		return
	}

	// get ip and port
	ip, port := stun.GetIP(ser)
	go stun.SendIndication(ser)

	fmt.Printf("Your address %s:%d\n", ip.String(), port)

	// _, port, _ := net.SplitHostPort(ser.LocalAddr().String())
	// fmt.Printf("Local port: %s\n", port)

	var secondIP string
	fmt.Scanln(&secondIP)

	ip2, port2, err := net.SplitHostPort(secondIP)
	if err != nil {
		fmt.Printf("cant parse: %v", err)
	}

	p2, _ := strconv.Atoi(port2)

	dst = &net.UDPAddr{
		Port: p2,
		IP:   net.ParseIP(ip2),
	}

	go sender()
	go reciever()
	<-chclose
}

func sender() {
	in := bufio.NewReader(os.Stdin)
	for {
		line, _, _ := in.ReadLine()
		// fmt.Fprintln(os.Stderr, line)
		if string(line) == "exit" {
			chclose <- struct{}{}
		}
		_, err := ser.WriteToUDP(line, dst)
		if err != nil {
			fmt.Printf("WriteToUDP error  %v", err)
			return
		}
		// fmt.Fprintln(os.Stderr, "sended")
	}
}

func reciever() {
	p := make([]byte, 2048)
	for {
		n, addr, err := ser.ReadFromUDP(p)
		// fmt.Fprintln(os.Stderr, "readed", n)
		if err != nil {
			fmt.Printf("Reader error %v\n", err)
			return
		}
		if n == 0 || addr != nil {
			fmt.Printf("[%s:%d]: %s\n", addr.IP.String(), addr.Port, string(p[:n]))
		}
	}
}
