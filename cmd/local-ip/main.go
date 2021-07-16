package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
)

var stunServerName = "stun.l.google.com:19302"

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

	dst, err := net.ResolveUDPAddr("udp", stunServerName)
	if err != nil {
		fmt.Printf("ResolveUDPAddr error  %v", err)
		return
	}

	_, err = ser.WriteToUDP(getData(), dst)
	if err != nil {
		fmt.Printf("WriteToUDP error  %v", err)
		return
	}

	p := make([]byte, 2048)
	var n int
	n, err = bufio.NewReader(ser).Read(p)
	if err == nil {
		fmt.Println(p[:n])

		ip := make([]byte, 4)
		binary.BigEndian.PutUint32(ip, 0x2112A442)
		start := 20 + 4 + 4

		for i := 0; i < 4; i++ {
			ip[i] ^= p[start+i]
		}

		fmt.Println(ip)
	} else {
		fmt.Printf("Reader error %v\n", err)
	}
	ser.Close()
}

func getData() []byte {
	res := make([]byte, 20)

	// query type
	binary.BigEndian.PutUint16(res[0:], 0x0001)
	// msg length
	binary.BigEndian.PutUint16(res[2:], 0x0000)
	// magic cookie
	binary.BigEndian.PutUint32(res[4:], 0x2112A442)

	// transaction id
	for i := 5; i < 20; i++ {
		res[i] = byte(i)
	}

	return res
}
