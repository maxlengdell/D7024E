package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	//"strconv"
	"strings"
)

type Network struct {
	table *RoutingTable
}

type Message struct {
	Type           string
	SenderContact  Contact
	TargetContact  Contact
	TargetHash     string
	Data           []byte
	ReturnContacts []Contact
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

func GetNetworkMessageChannel(port int) (chan InternalMessage, error) {
	msgChan := make(chan InternalMessage)
	return msgChan, nil
}

func ShovelMessages(conn *net.UDPConn, msgChan chan InternalMessage) error {
	recv := make([]byte, 2048)
	for {
		n, remoteAddr, err := conn.ReadFromUDP(recv)
		if err != nil {
			fmt.Errorf("Failed to read from connection: %v", err)
			return err
		}
		var m Message
		err = json.Unmarshal([]byte(string(recv[:n])), &m)
		if err != nil {
			// FIXME: assuming that unmarshal failure means CLI command.
			m = cliParser(string(recv[:n]))
		}
		msgChan <- InternalMessage{m, *conn, *remoteAddr}
	}
	return nil
}

// TODO: Replace with GetNetworkMessageChannel and ShovelMessages??
func Listen(ip string, port int, msgChan chan InternalMessage) error {
	// TODO
	//Port 8080 för ping -> besvara meddelande
	//Port 4000 för lookup
	//nc -u IP PORT för att testa denna funktion

	l, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	})
	if err != nil {
		//fmt.Println("Error listening:", err.Error())
		fmt.Errorf("Error listening: %v", err)
		return err
		//os.Exit(1)
	}
	defer l.Close()
	//fmt.Println("Listening on " + ip + ":" + strconv.Itoa(port))
	fmt.Printf("Listening on %v:%d\n", ip, port)
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
		fmt.Println("\nReceived Listen:", m.Type, m.SenderContact.Address)
		msgChan <- InternalMessage{m, *l, *remoteAddr}
	}
	return nil
}

func cliParser(msg string) Message {
	fmt.Println("Incoming", msg)
	var resp Message
	cmds := strings.Fields(msg)
	switch cmds[0] {
	case "put":
		fmt.Println("Received put")
		resp = Message{
			Type: "put",
			Data: []byte(cmds[1]),
		}
	case "get":
		fmt.Println("Received get")
		resp = Message{
			Type: "get",
		}
	case "exit":
		fmt.Println("Received exit")
		resp = Message{
			Type: "exit",
		}
	}
	return resp
}

// ContactConnection converts a Contact to a UDPAddr.
func ContactUDPAddress(contact *Contact) (*net.UDPAddr, error) {
	addr := contact.Address
	return net.ResolveUDPAddr("udp", addr)
}

// ContactConnection returns a network connection to the contact.
// The caller is responsible for closing the connection when done.
func ContactConnection(contact *Contact) (*net.UDPConn, error) {
	addr, err := ContactUDPAddress(contact)
	if err != nil {
		fmt.Errorf("Could not get the UDP address of contact %v: %v", contact, err)
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	return conn, err
}

// SendMessage opens a network connection to the contact and sends the given
// message over that connection, then it returns the response Message. The
// connection is closed. This function blocks while waiting for the response.
func SendMessage(contact *Contact, msg Message) (Message, error) { // TODO: return connection?
	conn, err := ContactConnection(contact)
	defer conn.Close()
	if err != nil {
		fmt.Errorf("Failed to open connection to %v: %v", contact, err)
		return Message{}, err
	}
	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Errorf("Could not convert Message %v to JSON: %v", msg, err)
		return Message{}, err
	}
	//return SendJSONMessage(conn, jsonMsg)
	response, _ := SendJSONMessage(conn, jsonMsg)
	var m Message
	json.Unmarshal(response, &m)
	return m, nil
}

// Sending self contact
func (network *Network) SendContactNode(conn *net.UDPConn, dest *net.UDPAddr, returnContacts []Contact) {
	m := Message{
		Type:           "LookUpNode-response",
		SenderContact:  network.table.me,
		ReturnContacts: returnContacts,
	}
	msg, _ := json.Marshal(m)
	fmt.Println("Returned contact manifest", m)
	SendJSONViaUDP(conn, dest, msg)

}
func SendJSONViaUDP(conn *net.UDPConn, dest *net.UDPAddr, msg []byte) error {
	_, err := conn.WriteToUDP(msg, dest)
	//fmt.Println("Connection: ", msg, " TO: ", dest)
	return err
}

// SendJSONMessage sends a JSON-encoded message on the given connection and
// then returns the (possibly JSON-encoded) response.
// This function does not close the connection. This function blocks while
// waiting for the response.
func SendJSONMessage(conn *net.UDPConn, msg []byte) ([]byte, error) {
	_, err := conn.Write(msg)
	fmt.Println("Connection: ", conn.LocalAddr(), conn.RemoteAddr())
	response := make([]byte, 2048)
	return response, err
}

// PingMessage takes a contact and returns a ping Message.
func PingMessage(contact *Contact) Message {
	return Message{
		Type:          "ping",
		SenderContact: *contact,
	}
}

// SendPingMessage sends a ping Message to the given contact and then
// returns the response message. This function blocks while waiting for
// the response.
func SendPingMessage(contact *Contact) (Message, error) {
	msg := PingMessage(contact)
	return SendMessage(contact, msg)
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

		fmt.Println("Confirmed alive", string(recv))

		if m.Type == "ping" {
			//fmt.Println("ID: not set", m.SenderContact.ID)
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
	}
	fmt.Println("PING ack: ", response.SenderContact)
	msg, _ := json.Marshal(response)
	l.WriteToUDP(msg, remoteAddr)
}

func (network *Network) SendFindContactMessage(contact *Contact, knownContact *Contact, contactChan chan []Contact) { //contact is the contact to "find"
	//fmt.Printf("Routing table:\n%s\n", network.table.String())
	// FIND_NODE request to bootstrap node
	var MessageRecv Message
	recv := make([]byte, 2048)

	//fmt.Println("SENDING FIND CONTACT", knownContact, contact.ID)

	l, err := net.Dial("udp", knownContact.Address)
	defer l.Close()
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	m := Message{
		Type:          "LookUpNode",
		SenderContact: network.table.me, //Self
		//BootstrapNode
		TargetContact: *contact, //Bootstrap node
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

	fmt.Println("Received FIND_NODE response", MessageRecv, contactChan)
	contactChan <- MessageRecv.ReturnContacts
	//Lyssna efter svar
	//Returnera grannar
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(recipient *Contact, data []byte) {
	// TODO
	m := Message{
		Type:          "StoreData",
		SenderContact: network.table.me,
		Data:          data,
	}
	SendMessage(recipient, m)
}
