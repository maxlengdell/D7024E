package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"sort"
	"strconv"
	"time"
)

type Kademlia struct {
	Net Network
}

var numberOfParallelRequests int = 3
var timeoutDur int = 1

func chooseNContacts(shortlist, visited []Contact, n int) []Contact {
	//TODO
	rand.Seed(time.Now().Unix())
	var returnArr []Contact
	if len(shortlist) < n || len(shortlist) == 0 {
		return shortlist
	}
	for i := 0; i < n; i++ {
		index := rand.Int() % len(shortlist)
		returnArr = append(returnArr, shortlist[index])
	}

	return returnArr
}
func (kademlia *Kademlia) LookupContact(target *Contact, conn *net.UDPConn, addr *net.UDPAddr) []Contact {
	//Locate k closest nodes

	shortlist := kademlia.Net.table.FindClosestContacts(target.ID, bucketSize) //3 närmsta grannarna
	var visitedNodes []Contact
	var closestNode Contact = kademlia.Net.table.me

	if len(shortlist) != 0 {
		closestNode = shortlist[0]
	}

	fmt.Println("Shortlist: ", len(shortlist))

	var alpha1 []Contact = chooseNContacts(shortlist, visitedNodes, numberOfParallelRequests)
	contactChan := make(chan []Contact, len(alpha1))

	for i, node := range alpha1 { //Alpha 1
		fmt.Println("Looping neighbour: ", i, node.Address)
		go kademlia.Net.SendFindContactMessage(target, &node, contactChan)
		visitedNodes = append(visitedNodes, node)
	}
	for i, _ := range alpha1 {
		fmt.Println("in loop alpha1")
		select {
		case recievedContacts := <-contactChan: //Recieved responses from findContactMessage
			for _, contact := range recievedContacts {
				shortlist = append(shortlist, contact)
				fmt.Println("Recieved contact: ", i, contact)
			}
		case <-time.After(time.Duration(timeoutDur) * time.Second):
			fmt.Println("*********TIMEOUT alpha 1********")
			//Remove node from shortlist
			break
		}
	}
	var madeProgress bool = true

	for madeProgress {
		madeProgress = false
		var alpha2 []Contact = chooseNContacts(shortlist, visitedNodes, numberOfParallelRequests)
		alpha2Channel := make(chan []Contact, len(alpha2))

		for j, node := range alpha2 { //Alpha 2
			fmt.Println("in loop alpha2", j)
			go kademlia.Net.SendFindContactMessage(target, &node, alpha2Channel)
			visitedNodes = append(visitedNodes, node)
		}
	loop:
		for {
			select {
			case recievedContacts := <-alpha2Channel:
				for _, contact := range recievedContacts {
					fmt.Println("Appending alpha 2")
					shortlist = append(shortlist, contact)
					if contact.Less(&closestNode) {
						closestNode = contact
						madeProgress = true
					}
				}
			case <-time.After(time.Duration(timeoutDur) * time.Second):
				fmt.Println("*********TIMEOUT alpha2********")
				break loop
			}
		}

		if len(shortlist) >= bucketSize {
			fmt.Println("Shortlist size exceeded bucketsize")
			break
		}
	}
	fmt.Println("shortlist: ", shortlist, closestNode)
	kademlia.Net.SendContactNode(conn, addr, shortlist)
	return shortlist
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO TOM LookUpData
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO MAX Store
	contactChan := make(chan []Contact)
	//<key,value>
	storeContact := NewContact(NewKademliaID(Hash(data)), "")
	neighbours := kademlia.Net.table.FindClosestContacts(kademlia.Net.table.me.ID, numberOfParallelRequests)
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
