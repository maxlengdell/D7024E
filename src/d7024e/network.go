package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
)

type Network struct {
}
type Message struct {
	Type   string
	Sender string
	Body   string
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
		recv := make([]byte, 100)
		n, remoteAddr, err := l.ReadFromUDP(recv)
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		var m Message
		json.Unmarshal([]byte(string(recv[:n])), &m)

		fmt.Println("Received Listen: ", m)
		go handleMessage(&m, l, remoteAddr)
	}
}
func handleMessage(m *Message, l *net.UDPConn, remoteAddr *net.UDPAddr) {
	var resp Message
	switch m.Type {
	case "ping":
		resp = Message{
			"ping",
			"",
			"ack",
		}

	}
	msg, _ := json.Marshal(resp)
	l.WriteToUDP(msg, remoteAddr)
	//fmt.Println("SENDING BACK", string(msg), remoteAddr.String())
}

func Bootstrap(ip string, port int) (network *Network, contact *Contact) {
	/* 	Create id, contact and network
	*	This node is the first node of the network.
	 */
	id := NewRandomKademliaID()
	myContact := NewContact(id, (ip + ":" + strconv.Itoa(port)))
	net := Network{}
	return &net, &myContact

}
func JoinNetwork(knownIP string, myip string, port int) (network *Network, contact *Contact) {
	/*	This Node is about to join a existing network.
		Create new bucket
		Create contact for known node
		Create contact for self
		Check so known node is alive
		if alive:
			Add to bucket
	*/
	buck := newBucket()

	knownContact := NewContact(nil, knownIP+":"+strconv.Itoa(port))
	myContact := NewContact(NewRandomKademliaID(), (myip + ":" + strconv.Itoa(port)))
	net := Network{}
	err := net.SendPingMessage(&knownContact)
	if err != nil {
		buck.AddContact(knownContact)
	}
	return &net, &myContact
}

func (network *Network) SendPingMessage(contact *Contact) error {
	// contact is assumed to be another node, not self.
	recv := make([]byte, 60)
	addr := contact.Address
	fmt.Println(addr)
	l, err := net.Dial("udp", addr)
	defer l.Close()
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	m := Message{
		"ping",
		"penis",
		"",
	}
	msg, _ := json.Marshal(m)
	_, writeErr := l.Write(msg)
	//fmt.Println("RAW TO BE SENT: ", msg)

	if writeErr != nil {
		fmt.Println("Could not send msg", err)
	} else {
		fmt.Println("SENT PING TO: " + addr)
	}

	for {
		n, _ := l.Read(recv)
		var m Message
		json.Unmarshal([]byte(string(recv[:n])), &m)

		fmt.Println("Received ping: ", m)
		if m.Type == "ping" {
			return nil
		}
	}
	return writeErr
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
