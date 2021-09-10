package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
)

type Network struct {
	table *RoutingTable
}
type Message struct {
	Type          string
	SenderContact Contact
	TargetContact Contact
	TargetHash	  string
	Data          string
}
type InternelMessage struct {
	msg			Message
	conn 		*net.UDPConn
	remoteAddr	*net.UDPAddr
}
func (network *Network) Listen(ip string, port int, msgChan chan InternelMessage) {
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
		msgChan<-InternelMessage{m, l, remoteAddr}
	}
}

func (network *Network) SendPingMessage(contact *Contact) (string, error) {
	// contact is assumed to be another node, not self.
	recv := make([]byte, 2048)
	addr := contact.Address
	
	l, err := net.Dial("udp", addr)
	defer l.Close()
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	m := Message{
		Type: "ping",
		SenderContact: NewContact(nil, ""),
		Data: "",
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

func (network *Network) SendPingAckMessage(l *net.UDPConn, remoteAddr *net.UDPAddr){
	response := Message{
		Type: "ping",
		SenderContact: network.table.me,
		Data: "ack",
	}
	msg, _ := json.Marshal(response)
	l.WriteToUDP(msg, remoteAddr)
}

func (network *Network) SendFindContactMessage(contact *Contact) { //contact is the contact to "find"
	// FIND_NODE request to bootstrap node
	addr := contact.Address
	l, err := net.Dial("udp", addr)
	defer l.Close()
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	m := Message{
		Type: "LookUpNode",
		SenderContact: network.table.me,
		TargetContact: *contact,
	}
	msg, _ := json.Marshal(m)
	_, writeErr := l.Write(msg)
	//Handle err
	if writeErr != nil {
		fmt.Println("Could not send msg", err)
	} else {
		fmt.Println("SENT LOOKUPNODE: " + string(msg))
	}
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
