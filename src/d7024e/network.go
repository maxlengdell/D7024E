package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Network struct {
	table *RoutingTable
}
type Message struct {
	Type          string
	SenderContact Contact
	TargetContact Contact
	TargetHash    string
	Data          string
}
type InternelMessage struct {
	msg        Message
	conn       *net.UDPConn
	remoteAddr *net.UDPAddr
}

func handleErr(err error) {
	//TODO
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

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
		//os.Exit(1)
	}
	defer l.Close()
	fmt.Println("Listening on " + ip + ":" + strconv.Itoa(port))
	for {
		recv := make([]byte, 2048)
		n, remoteAddr, err := l.ReadFromUDP(recv)
		handleErr(err)
		var m Message

		MarErr := json.Unmarshal([]byte(string(recv[:n])), &m)
		if MarErr != nil {
			//Unable to Unmarshal message
			fmt.Println("CLI COMMAND received")
			m = network.cliParser(string(recv[:n]))
			//Handle string parser to get syntax
		}
		//fmt.Printf("\nReceived Listen: %#v \n", m)
		//go network.handleMessage(&m, l, remoteAddr)
		fmt.Printf("\nReceived Listen: %#v \n", m)
		msgChan <- InternelMessage{m, l, remoteAddr}
	}
}
func (network *Network) cliParser(msg string) Message {
	fmt.Println("Incoming", msg)
	var resp Message
	cmds := strings.Fields(msg)
	switch cmds[0] {
	case "put":
		fmt.Println("Received put")
		resp = Message{
			Type:          "ping",
			SenderContact: NewContact(nil, ""),
			Data:          "",
		}
	case "get":
		fmt.Println("Received get")
		resp = Message{
			Type:          "ping",
			SenderContact: NewContact(nil, ""),
			Data:          "",
		}
	case "exit":
		fmt.Println("Received exit")
		resp = Message{
			Type:          "ping",
			SenderContact: NewContact(nil, ""),
			Data:          "",
		}
	}
	return resp
}

// func (network *Network) handleMessage(m *Message, l *net.UDPConn, remoteAddr *net.UDPAddr) {
// 	var resp Message
// 	switch m.Type {
// 	case "ping":
// 		resp = Message{
// 			"ping",
// 			NewContact(nil, ""),
// 			"ack",
// 		}
// 	}
// 	msg, _ := json.Marshal(resp)
// 	//fmt.Println("SENDING BACK", network.me, string(msg))

// 	l.WriteToUDP(msg, remoteAddr)
// }

// func Bootstrap(ip string, port int) (network *Network) {
// 	/* 	Create id, contact and network
// 	*	This node is the first node of the network.
// 	 */
// 	//buck := newBucket()

// 	id := NewRandomKademliaID()
// 	myContact := NewContact(id, (ip + ":" + strconv.Itoa(port)))
// 	net := Network{myContact}
// 	return &net

// }
// func JoinNetwork(knownIP string, myip string, port int) (network *Network) {
// 	/*	This Node is about to join a existing network.
// 		Create new bucket
// 		Create contact for known node
// 		Create contact for self
// 		Check so known node is alive
// 		if alive:
// 			Add to bucket
// 	*/
// 	buck := newBucket()

// 	knownContact := NewContact(nil, knownIP+":"+strconv.Itoa(port))
// 	myContact := NewContact(NewRandomKademliaID(), (myip + ":" + strconv.Itoa(port)))
// 	net := Network{myContact}
// 	knownID, err := net.SendPingMessage(&knownContact)
// 	net.SendFindContactMessage(&knownContact)
// 	//Lookup
// 	if err != nil {
// 		bootstrapContact := NewContact(NewKademliaID(knownID), knownIP+":"+strconv.Itoa(port))
// 		buck.AddContact(bootstrapContact)
// 	}
// 	fmt.Println("Bucket: ", buck.list.Front())
// 	return &net
// }

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
		Type:          "ping",
		SenderContact: NewContact(nil, ""),
		Data:          "",
	}
	msg, _ := json.Marshal(m)
	_, writeErr := l.Write(msg)
	//Handle err
	if writeErr != nil {
		fmt.Println("Could not send msg", err)
	} else {
		fmt.Println("SENT PING TO: " + m.SenderContact.Address)
	}

	for {
		n, _ := l.Read(recv)
		var m Message
		json.Unmarshal([]byte(string(recv[:n])), &m)

		fmt.Println("Confirmed alive")
		if m.Type == "ping" {
			fmt.Println("ID: ", m.SenderContact.ID)
			return m.SenderContact.ID.String(), nil
		}
		break
	}
	return "", writeErr
}

func (network *Network) SendPingAckMessage(l *net.UDPConn, remoteAddr *net.UDPAddr) {
	response := Message{
		Type:          "ping",
		SenderContact: network.table.me,
		Data:          "ack",
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
		Type:          "LookUpNode",
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
