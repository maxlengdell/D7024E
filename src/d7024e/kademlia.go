package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
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
func (kademlia *Kademlia) LookupContact(target *Contact, retChan chan []Contact) {
	//Locate k closest nodes
	shortlist := kademlia.Net.table.FindClosestContacts(target.ID, bucketSize) //3 närmsta grannarna
	var visitedNodes []Contact
	var closestNode Contact = shortlist[0]

	var alpha1 []Contact = chooseNContacts(shortlist, visitedNodes, numberOfParallelRequests)
	fmt.Println("Shortlist: ", len(shortlist), "nContacts: ", alpha1)

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
					contact.CalcDistance(kademlia.Net.table.me.ID)
					fmt.Println("Appending alpha 2: ", closestNode)
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
	fmt.Println("shortlist: ", shortlist, "closest node: ", closestNode)

	retChan <- shortlist
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
func (kademlia *Kademlia) findNode(target *Contact) []Contact {
	contactChan := make(chan []Contact, 1)
	kademlia.LookupContact(target, contactChan)
	nodes := <-contactChan
	fmt.Println("Lookup done, returned contacts: ", nodes)
	return nodes
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
	kadem := Kademlia{net}

	knownID, err := kadem.Net.SendPingMessage(&knownContact)
	//_, err := SendPingMessage(&knownContact)
	//knownID := ""

	bootstrapContact := NewContact(NewKademliaID(knownID), knownIP+":"+strconv.Itoa(port))
	if err == nil {
		net.table.AddContact(bootstrapContact)
	}
	//fmt.Println("Known contact node: ", bootstrapContact)
	//net.SendFindContactMessage(&myContact, &knownContact, contactChan)

	closeNodes := kadem.findNode(&kadem.Net.table.me)
	for _, node := range closeNodes {
		kadem.Net.table.AddContact(node)
	}
	fmt.Println("############## ------ Join network contact chan ------- ############", closeNodes)

	return &kadem
}

func (kademlia *Kademlia) HandleMessage(msgChan chan InternalMessage) {
	for {
		var m = <-msgChan
		fmt.Println("Internal message recieved:", m.msg.Type)
		switch m.msg.Type {
		case "ping":
			fmt.Println("Adding ping sender to contacts", m.msg.SenderContact.Address)
			kademlia.Net.table.AddContact(m.msg.SenderContact)
			go kademlia.Net.SendPingAckMessage(&m.conn, &m.remoteAddr)
		case "LookUpNode":
			go kademlia.HandleFindNode(m)
			//go kademlia.LookupContact(&m.msg.TargetContact, &m.conn, &m.remoteAddr)
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
func (kademlia *Kademlia) HandleFindNode(m InternalMessage) {
	resp := Message{
		Type:           "find-node-resp",
		ReturnContacts: kademlia.Net.table.FindClosestContacts(m.msg.TargetContact.ID, numberOfParallelRequests),
	}
	jsonMsg, err := json.Marshal(resp)
	handleErr(err)
	m.conn.WriteToUDP(jsonMsg, &m.remoteAddr)
	fmt.Println("Returned closest contacts to target: ", m.msg.TargetContact.Address, "neighbours: ", resp.ReturnContacts)

}
func Hash(data []byte) string {
	//Hash data to sha1 and return
	sh := sha1.Sum(data)
	return hex.EncodeToString(sh[:])
}
