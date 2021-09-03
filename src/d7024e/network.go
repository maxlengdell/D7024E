package d7024e

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Network struct {
}

func Listen(ip string, port int) {
	// TODO
	//Port 8080 för ping -> besvara meddelande
	//Port 4000 för lookup
	l, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	})
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + ip + ":" + strconv.Itoa(port))
	for {
		message := make([]byte, 20)
		rlen, remote, err := l.ReadFromUDP(message[:])
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		data := strings.TrimSpace(string(message[:rlen]))
		fmt.Println("Received: ", data, remote)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
