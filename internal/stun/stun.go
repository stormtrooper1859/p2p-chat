package stun

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

const (
	secretKey uint32 = 0x2112A442
)

var stunServerName = "stun.l.google.com:19302"
var stunServerAddr *net.UDPAddr

func getData() []byte {
	res := make([]byte, 20)

	// query type
	binary.BigEndian.PutUint16(res[0:], 0x0001)
	// msg length
	binary.BigEndian.PutUint16(res[2:], 0x0000)
	// magic cookie
	binary.BigEndian.PutUint32(res[4:], secretKey)

	// transaction id
	for i := 5; i < 20; i++ {
		res[i] = byte(i)
	}

	return res
}

func GetIP(conn *net.UDPConn) (resultIp net.IP, port uint16) {
	var err error
	stunServerAddr, err = net.ResolveUDPAddr("udp", stunServerName)
	if err != nil {
		fmt.Printf("ResolveUDPAddr error  %v", err)
		return
	}

	_, err = conn.WriteToUDP(getData(), stunServerAddr)
	if err != nil {
		fmt.Printf("WriteToUDP error  %v", err)
		return
	}

	responseBuffer := make([]byte, 2048)
	_, _, err = conn.ReadFromUDP(responseBuffer)
	ip := make([]byte, 4)
	if err == nil {
		// fmt.Println(responseBuffer[:n])

		binary.BigEndian.PutUint32(ip, secretKey)
		start := 20 + 4 + 4

		for i := 0; i < 4; i++ {
			ip[i] ^= responseBuffer[start+i]
		}

		port = binary.BigEndian.Uint16(responseBuffer[start-2:])

		// fmt.Println(ip)
	} else {
		fmt.Printf("Reader error %v\n", err)
		return nil, 0
	}

	return net.IPv4(ip[0], ip[1], ip[2], ip[3]), port
}

func getIndicationData() []byte {
	res := make([]byte, 20)

	// query type
	binary.BigEndian.PutUint16(res[0:], 0x0010)
	// msg length
	binary.BigEndian.PutUint16(res[2:], 0x0000)
	// magic cookie
	binary.BigEndian.PutUint32(res[4:], secretKey)

	// transaction id
	for i := 5; i < 20; i++ {
		res[i] = byte(i)
	}

	return res
}

func SendIndication(conn *net.UDPConn) {
	data := getIndicationData()
	for {
		_, err := conn.WriteToUDP(data, stunServerAddr)
		// fmt.Fprintln(os.Stderr, "indication sended")
		if err != nil {
			fmt.Printf("SendIndication WriteToUDP error  %v", err)
			return
		}
		time.Sleep(120 * time.Millisecond)
	}
}
