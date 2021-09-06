package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Network struct {
}
type Message struct {
	Name string
	Body string
	Time int64
	/*
		Structure of network packets
			m := Message{
			"ping",
			"",
			time.Now().Unix(),
		}
	*/
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
	addr := contact.Address + ":8080"
	fmt.Println(addr)
	l, err := net.Dial("udp", addr)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	m := Message{
		"ping",
		"",
		time.Now().Unix(),
	}
	msg, _ := json.Marshal(m)
	for {
		_, err := l.Write([]byte(msg))
		if err != nil {
			fmt.Println("Could not send msg", err)
		}
		//remote, err := l.Read(message)
		//data := strings.TrimSpace(string(message))
		//fmt.Println("Received: ", data, remote)
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
