package d7024e
import (
	"fmt"
	"strconv"
)

type Kademlia struct {
	Net Network
}
var numberOfParrallellRequests int = 3

func (kademlia *Kademlia) LookupContact(target *Contact) {
	fmt.Println("TODO LOOKUPCONTRACT")
	// TODO
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

	fmt.Println("ID",knownID,"ERR",err)
	//kademlia.Net.SendFindContactMessage(&knownContact)
	//Lookup
	if err != nil {
		bootstrapContact := NewContact(NewKademliaID(knownID), knownIP+":"+strconv.Itoa(port))
		table.AddContact(bootstrapContact)
	}
	kadem := Kademlia{net}
	return &kadem
}

func (kademlia *Kademlia) HandleMessage(msgChan chan InternelMessage) {
	for {
		var m = <-msgChan
		fmt.Println("Internel message recieved:", m.msg.Type)
		switch m.msg.Type {
		case "ping":
			go kademlia.Net.SendPingAckMessage(m.conn, m.remoteAddr)
		case "LookUpNode":
			go kademlia.LookupContact(&m.msg.TargetContact)
		case "LookUpData":
			fmt.Println("LookUpData RECIEVED, TODO IMPLEMENTATION")
		case "StoreData":
			fmt.Println("StoreData RECIEVED, TODO IMPLEMENTATION")
		}
	}
}