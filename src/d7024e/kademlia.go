package d7024e

import (
	"fmt"
	"net"
	"strconv"
)

type Kademlia struct {
	Net Network
}

var numberOfParrallellRequests int = 3

func (kademlia *Kademlia) LookupContact(target *Contact, conn *net.UDPConn, addr *net.UDPAddr) *Contact {
	fmt.Println("TODO LOOKUPCONTRACT, target contact: ", target, " My id: ", kademlia.Net.table.me)
	//Locate k closest nodes
	contactChan := make(chan Contact)

	neighbours := kademlia.Net.table.FindClosestContacts(target.ID, numberOfParrallellRequests) //3 närmsta grannarna
	//Check if target contact is closest in neighbours. If so, return target contact.
	if neighbours[0].ID.Equals(kademlia.Net.table.me.ID) {
		return &kademlia.Net.table.me
	} else {
		for _, node := range neighbours {
			contactChan <- kademlia.Net.SendFindContactMessage(target, &node)
		}
	}
	return nil
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

func Bootstrap(ip string, port int) (kademlia *Kademlia) {
	/* 	Create id, contact and network
	*	This node is the first node of the network.
	 */
	id := NewRandomKademliaID()
	myContact := NewContact(id, (ip + ":" + strconv.Itoa(port)))

	table := NewRoutingTable(myContact)
	net := Network{table}
	kadem := Kademlia{net}
	return &kadem

}
func JoinNetwork(knownIP string, myip string, port int) (kademlia *Kademlia) {
	/*	This Node is about to join a existing network.
		Create new bucket
		Create contact for known node
		Create contact for self
		Check so known node is alive
		if alive:
			Add to bucket
	*/

	knownContact := NewContact(nil, knownIP+":"+strconv.Itoa(port))
	myContact := NewContact(NewRandomKademliaID(), (myip + ":" + strconv.Itoa(port)))

	table := NewRoutingTable(myContact)
	net := Network{table}
	knownID, err := net.SendPingMessage(&knownContact)

	fmt.Println("ID received: ", knownID, err)
	bootstrapContact := NewContact(NewKademliaID(knownID), knownIP+":"+strconv.Itoa(port))
	if err == nil {
		net.table.AddContact(bootstrapContact)
	}
	fmt.Println("Known contact node: ", bootstrapContact)

	net.SendFindContactMessage(&myContact, &knownContact)

	//Lookup

	kadem := Kademlia{net}
	return &kadem
}

func (kademlia *Kademlia) HandleMessage(msgChan chan InternalMessage) {
	for {
		var m = <-msgChan
		fmt.Println("Internal message recieved:", m.msg.Type)
		switch m.msg.Type {
		case "ping":
			go kademlia.Net.SendPingAckMessage(&m.conn, &m.remoteAddr)
		case "LookUpNode":
			go kademlia.LookupContact(&m.msg.TargetContact, &m.conn, &m.remoteAddr)
		case "LookUpData":
			fmt.Println("LookUpData RECIEVED, TODO IMPLEMENTATION")
		case "StoreData":
			fmt.Println("StoreData RECIEVED, TODO IMPLEMENTATION")
		case "put":
			fmt.Println("Store data")
		case "get":
			fmt.Println("Get data")
		case "exit":
			//Kill network object
			fmt.Println("Quitting node...")
		}

	}
}
