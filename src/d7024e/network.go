package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
)

type Network struct {
	self Contact
}
type Message struct {
	Type          string
	SenderContact Contact
	Body          string

	/*
		Structure of network packets
			m := Message{
			"ping",
			"",
			time.Now().Unix(),
		}
	*/
}

func handleErr() {
	//TODO
}
func (network *Network) Listen(ip string, port int) {
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
		recv := make([]byte, 2048)
		n, remoteAddr, err := l.ReadFromUDP(recv)
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		var m Message
		json.Unmarshal([]byte(string(recv[:n])), &m)

		fmt.Printf("\nReceived Listen: %#v \n", m)
		go network.handleMessage(&m, l, remoteAddr)
	}
}
func (network *Network) handleMessage(m *Message, l *net.UDPConn, remoteAddr *net.UDPAddr) {
	var resp Message
	switch m.Type {
	case "ping":
		resp = Message{
			"ping",
			NewContact(nil, ""),
			"ack",
		}
	}
	msg, _ := json.Marshal(resp)
	//fmt.Println("SENDING BACK", network.me, string(msg))

	l.WriteToUDP(msg, remoteAddr)
}

func Bootstrap(ip string, port int) (network *Network) {
	/* 	Create id, contact and network
	*	This node is the first node of the network.
	 */
	//buck := newBucket()

	id := NewRandomKademliaID()
	myContact := NewContact(id, (ip + ":" + strconv.Itoa(port)))
	net := Network{myContact}
	return &net

}
func JoinNetwork(knownIP string, myip string, port int) (network *Network) {
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
	net := Network{myContact}
	knownID, err := net.SendPingMessage(&knownContact)
	net.SendFindContactMessage(&knownContact)
	//Lookup
	if err != nil {
		bootstrapContact := NewContact(NewKademliaID(knownID), knownIP+":"+strconv.Itoa(port))
		buck.AddContact(bootstrapContact)
	}
	fmt.Println("Bucket: ", buck.list.Front())
	return &net
}

func (network *Network) SendPingMessage(contact *Contact) (string, error) {
	// contact is assumed to be another node, not self.
	recv := make([]byte, 2048)
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
		NewContact(nil, ""),
		"",
	}
	msg, _ := json.Marshal(m)
	_, writeErr := l.Write(msg)
	//Handle err
	if writeErr != nil {
		fmt.Println("Could not send msg", err)
	} else {
		fmt.Println("SENT PING TO: " + string(msg))
	}

	for {
		n, _ := l.Read(recv)
		var m Message
		json.Unmarshal([]byte(string(recv[:n])), &m)

		fmt.Println("Received ping: ", m)
		if m.Type == "ping" {
			//fmt.Println("ID: ", m.SenderContact.ID)
			return m.SenderContact.ID.String(), nil
		}
		break
	}
	return "", writeErr
}

func (network *Network) SendFindContactMessage(knownContact *Contact) {
	// FIND_NODE request to bootstrap node

}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
