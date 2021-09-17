package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net"
	"sort"
	"strconv"
	"time"
)

type Kademlia struct {
	Net Network
}

var numberOfParrallellRequests int = 3
var timeoutDuration int = 1

func (kademlia *Kademlia) LookupContact(target *Contact, conn *net.UDPConn, addr *net.UDPAddr) []Contact {
	//Locate k closest nodes
	var returnedContacts []Contact

	neighbours := kademlia.Net.table.FindClosestContacts(target.ID, numberOfParrallellRequests) //3 närmsta grannarna
	contactChan := make(chan []Contact, len(neighbours))
	var closestNeighbours []Contact

	kademlia.Net.table.me.CalcDistance(target.ID)

	fmt.Println("Neighbours: ", len(neighbours))
	var emptySlice []Contact
	meWrappedInSLice := append(emptySlice, kademlia.Net.table.me)
	if len(neighbours) == 0 {
		go kademlia.Net.SendContactNode(conn, addr, meWrappedInSLice)
		close(contactChan)
		fmt.Println("NO NEIGHBOURS!!")
		if target.Address != "" { //bootstrap-target.distance > me-target.distance
			fmt.Println("Adding contact: ", *target)
			kademlia.Net.table.AddContact(*target) //ADD AFTER EACH LOOKUP
		}
		return nil
	}
	for i, node := range neighbours {
		//If neighbour is further from target then self, -> return self
		//else return neighbour
		fmt.Println("Looping neighbour: ", i, node.Address)

		if kademlia.Net.table.me.distance.Less(node.distance) {
			fmt.Println("WE ARE CLOSER TO TARGET THEN: ", node.Address)
			//closestNeighbours = append(closestNeighbours, kademlia.Net.table.me)
			go kademlia.Net.SendContactNode(conn, addr, meWrappedInSLice)
			//close(contactChan)
			// if target.Address != "" {
			// 	fmt.Println("Adding contact: ", *target)
			// 	kademlia.Net.table.AddContact(*target) //ADD AFTER EACH LOOKUP
			// }
			break
		} else {
			fmt.Println("go routine: ", i)
			go kademlia.Net.SendFindContactMessage(target, &node, contactChan)
			closestNeighbours = append(closestNeighbours, node)

		}
	}
	fmt.Println("Closest neighbours: ", closestNeighbours)
	for range closestNeighbours {
		fmt.Println("in loop closest neighbour")
		select {
		case returnContact := <-contactChan: //Recieved responses from findContactMessage
			for _, contact := range returnContact {
				returnedContacts = append(returnedContacts, contact)
				fmt.Println("Return contact: ", contact)
			}
		case <-time.After(1 * time.Second):
			fmt.Println("*********TIMEOUT********")
			break
		}
	}
	//Sort numberofParrallell stycken closest contacts
	//fmt.Println("Sorting")
	// sortSliceByDistance(returnedContacts)
	// fmt.Println("Done sorting")
	if len(returnedContacts) != 0 {
		go kademlia.Net.SendContactNode(conn, addr, returnedContacts)
	}

	if target.Address != "" { //bootstrap-target.distance > me-target.distance
		fmt.Println("Adding contact: ", *target)
		kademlia.Net.table.AddContact(*target) //ADD AFTER EACH LOOKUP
	}
	fmt.Println("LookUpNode Complete")
	return nil
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO TOM LookUpData
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO MAX Store
	contactChan := make(chan []Contact)
	//<key,value>
	storeContact := NewContact(NewKademliaID(Hash(data)), "")
	neighbours := kademlia.Net.table.FindClosestContacts(kademlia.Net.table.me.ID, numberOfParrallellRequests)
	for _, node := range neighbours {
		go kademlia.Net.SendFindContactMessage(&storeContact, &node, contactChan)
	}

	returnContact := <-contactChan
	for _, contact := range returnContact {
		kademlia.Net.SendStoreMessage(&contact, data)
	}
}

func sortSliceByDistance(slice []Contact) {
	fmt.Println("SORTING INPUT: ", slice)
	sort.Slice(slice[:], func(i, j int) bool {
		return slice[i].distance.Less(slice[j].distance)
	})
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
	//fmt.Println("Known ID: ", knownID)

	bootstrapContact := NewContact(NewKademliaID(knownID), knownIP+":"+strconv.Itoa(port))
	if err == nil {
		net.table.AddContact(bootstrapContact)
	}
	//fmt.Println("Known contact node: ", bootstrapContact)
	contactChan := make(chan []Contact, 1)
	net.SendFindContactMessage(&myContact, &knownContact, contactChan)
	returnContacts := <-contactChan //{[contact1,contact2,contact3]}
	fmt.Println("############## ------ Join network contact chan ------- ############", returnContacts)
	for _, contact := range returnContacts {
		net.table.AddContact(contact) //ADD CONTACT AFTER FIRST LOOKUP ON SELF
	}

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
			fmt.Println("StoreData RECIEVED:", m.msg)
		case "put":
			fmt.Println("Store data", m.msg)
			kademlia.Store(m.msg.Data)
		case "get":
			fmt.Println("Get data")
		case "exit":
			//Kill network object
			fmt.Println("Quitting node...")
		}

	}
}
func Hash(data []byte) string {
	//Hash data to sha1 and return
	sh := sha1.Sum(data)
	return hex.EncodeToString(sh[:])
}
