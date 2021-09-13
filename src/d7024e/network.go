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
	ReturnContact Contact
}
type InternalMessage struct {
	msg        Message
	conn       net.UDPConn
	remoteAddr net.UDPAddr
}

func handleErr(err error) {
	//TODO
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
func Listen(ip string, port int, msgChan chan InternalMessage) {
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
			m = cliParser(string(recv[:n]))
			//Handle string parser to get syntax
		}
		//fmt.Printf("\nReceived Listen: %#v \n", m)
		//go network.handleMessage(&m, l, remoteAddr)
		fmt.Printf("\nReceived Listen: %#v \n", m)
		msgChan <- InternalMessage{m, *l, *remoteAddr}
	}
}
func cliParser(msg string) Message {
	fmt.Println("Incoming", msg)
	var resp Message
	cmds := strings.Fields(msg)
	switch cmds[0] {
	case "put":
		fmt.Println("Received put")
		resp = Message{
			Type:          "put",
			SenderContact: NewContact(nil, ""),
			Data:          "",
		}
	case "get":
		fmt.Println("Received get")
		resp = Message{
			Type:          "get",
			SenderContact: NewContact(nil, ""),
			Data:          "",
		}
	case "exit":
		fmt.Println("Received exit")
		resp = Message{
			Type:          "exit",
			SenderContact: NewContact(nil, ""),
			Data:          "",
		}
	}
	return resp
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
		Type:          "ping",
		SenderContact: network.table.me,
		Data:          "",
	}
	msg, _ := json.Marshal(m)
	_, writeErr := l.Write(msg)
	//Handle err
	if writeErr != nil {
		fmt.Println("Could not send msg", err)
	} else {
		fmt.Println("SENT PING TO: " + addr)
	}

	for {
		n, _ := l.Read(recv)
		var m Message
		json.Unmarshal([]byte(string(recv[:n])), &m)

		fmt.Println("Confirmed alive",m)
		if m.Type == "ping" {
			fmt.Println("ID: not set", m.SenderContact.ID)
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

func (network *Network) SendFindContactMessage(contact *Contact, knownContact *Contact, contactChan chan Contact) { //contact is the contact to "find"
	// FIND_NODE request to bootstrap node
	var MessageRecv Message
	recv := make([]byte, 2048)

	fmt.Println("SENDING FIND CONTACT", knownContact, contact.ID)

	l, err := net.Dial("udp", knownContact.Address)
	defer l.Close()
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	m := Message{
		Type:          "LookUpNode",
		SenderContact: network.table.me, //Self
		TargetContact: *contact,         //Bootstrap node
	}
	msg, _ := json.Marshal(m)
	_, writeErr := l.Write(msg)
	//Handle err
	if writeErr != nil {
		fmt.Println("Could not send msg", err)
	} else {
		fmt.Println("SENT LOOKUPNODE: "+string(msg)+" TO : ", knownContact)
	}
	//**Listen for response**
	n, _ := l.Read(recv)

	json.Unmarshal([]byte(string(recv[:n])), &MessageRecv)

	fmt.Println("Received FIND_NODE response", MessageRecv)
	contactChan <- MessageRecv.ReturnContact
	//Lyssna efter svar
	//Returnera grannar
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
