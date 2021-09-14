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
	contactChan := make(chan Contact, numberOfParrallellRequests)

	neighbours := kademlia.Net.table.FindClosestContacts(target.ID, numberOfParrallellRequests) //3 närmsta grannarna
	kademlia.Net.table.me.CalcDistance(target.ID)

	fmt.Println("Neighbours: ", neighbours)
	if len(neighbours) == 0 {
		kademlia.Net.table.AddContact(*target)
		go kademlia.Net.SendContactNode(conn, addr, kademlia.Net.table.me)
		fmt.Println("NO NEIGHBOURS!!")
	}
	for i, node := range neighbours {
		//If neighbour is further from target then self, -> return self
		//else return neighbour

		if kademlia.Net.table.me.distance.Less(node.distance) {
			fmt.Println("WE ARE CLOSER TO TARGET THEN: ", node)

			kademlia.Net.table.AddContact(*target)
			go kademlia.Net.SendContactNode(conn, addr, kademlia.Net.table.me)
			break
		} else {
			fmt.Println("go routine: ", i)
			go kademlia.Net.SendFindContactMessage(target, &node, contactChan)
			returnContact := <-contactChan
			fmt.Println("Return contact: ", returnContact)
			go kademlia.Net.SendContactNode(conn, addr, returnContact)

			fmt.Println("##########################################Contact channel: ", returnContact, conn, cap(kademlia.Net.table.buckets))
		}
	}
	fmt.Println("NUMBER OF CONTACTS IN RT: ", len(kademlia.Net.table.buckets))
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
	fmt.Println("My id: ", myContact.ID.String())
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
	fmt.Println("My id: ", myContact.ID.String())
	table := NewRoutingTable(myContact)
	net := Network{table}
	knownID, err := net.SendPingMessage(&knownContact)
	//_, err := SendPingMessage(&knownContact)
	//knownID := ""
	fmt.Println("Known ID: ", knownID)

	bootstrapContact := NewContact(NewKademliaID(knownID), knownIP+":"+strconv.Itoa(port))
	if err == nil {
		net.table.AddContact(bootstrapContact)
	}
	//fmt.Println("Known contact node: ", bootstrapContact)
	contactChan := make(chan Contact, numberOfParrallellRequests)
	net.SendFindContactMessage(&myContact, &knownContact, contactChan)
	fmt.Println("Join network contact chan", <-contactChan)

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
