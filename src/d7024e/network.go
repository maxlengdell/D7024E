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
	//nc -u IP PORT för att testa denna funktion
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
		//go handleMessage(data)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	// contact is assumed to be another node, not self.
	message := make([]byte, 20)
	addr := net.UDPAddr{
		Port: 8080,
		IP:   net.IP(contact.Address),
	}
	l, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	msg := "ping"
	for {
		_, err := l.WriteToUDP([]byte(msg), &addr)
		if err != nil {
			fmt.Printf("Could not send msg" + msg)
		}
		rlen, remote, err := l.ReadFromUDP(message[:])
		data := strings.TrimSpace(string(message[:rlen]))
		fmt.Println("Received: ", data, remote)

	}
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
