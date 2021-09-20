package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"sort"
	"strconv"
)

type Kademlia struct {
	Net         Network
	HashStorage []string
}

var numberOfParallelRequests int = 3
var timeoutDur int = 1

func chooseNContacts(shortlist, visited []Contact, n int) []Contact {
	var returnArr []Contact

	for i := 0; i < n; i++ {
		if i >= len(shortlist) {
			break
		}
		_, contains := Find(visited, shortlist[i])
		if !contains {
			returnArr = append(returnArr, shortlist[i])
		}
	}

	return returnArr
}

func (kademlia *Kademlia) LookupData(hash string, retContactChan chan []Contact, retDataChan chan string) {
	shortlist := kademlia.Net.table.FindClosestContacts(NewKademliaID(hash), bucketSize)
	var visitedNodes []Contact
	var closestNode Contact = shortlist[0]

	var alpha1 []Contact = chooseNContacts(shortlist, visitedNodes, numberOfParallelRequests)
	fmt.Println("FIND_DATA Shortlist: ", len(shortlist), "nContacts: ", alpha1)

	contactChan := make(chan []Contact, len(alpha1))
	dataChan := make(chan string, 1)
	removeChan := make(chan Contact, 1)

	for i, node := range alpha1 { //Alpha 1
		fmt.Println("FIND_DATA Looping neighbour: ", i, node.Address)
		go kademlia.Net.SendFindDataMessage(hash, &node, contactChan, removeChan, dataChan)
		visitedNodes = append(visitedNodes, node)
	}
	for i, _ := range alpha1 {
		fmt.Println("FIND_DATA in loop alpha1")
		select {
		case recievedContacts := <-contactChan: //Recieved responses from findContactMessage
			for _, contact := range recievedContacts {
				shortlist = append(shortlist, contact)
				fmt.Println("FIND_DATA Recieved contact: ", i, contact)
			}
		case recievedData := <-dataChan:
			retDataChan <- recievedData
			return
		case removeContact := <-removeChan:
			fmt.Println("*********FIND_DATA TIMEOUT alpha 1********")
			RemoveContact(shortlist, removeContact)
			//Remove node from shortlist
			break
		}
	}
	var madeProgress bool = true

	for madeProgress {
		madeProgress = false
		var alpha2 []Contact = chooseNContacts(shortlist, visitedNodes, numberOfParallelRequests)
		alpha2Channel := make(chan []Contact, len(alpha2))
		alpha2DataChannel := make(chan string, 1)
		alpha2TimeoutChannel := make(chan Contact, 1)

		for j, node := range alpha2 { //Alpha 2
			fmt.Println("in loop alpha2", j)
			go kademlia.Net.SendFindDataMessage(hash, &node, alpha2Channel, alpha2TimeoutChannel, alpha2DataChannel)
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
			case recievedData := <-alpha2DataChannel:
				retDataChan <- recievedData
				return
			case removeContact := <-alpha2TimeoutChannel:
				fmt.Println("*********TIMEOUT alpha2********")
				RemoveContact(shortlist, removeContact)
				break loop
			}
		}

		if len(shortlist) >= bucketSize {
			fmt.Println("Shortlist size exceeded bucketsize")
			break
		}
	}
	fmt.Println("shortlist: ", shortlist, "closest node: ", closestNode)

	retContactChan <- shortlist
}

func (kademlia *Kademlia) LookupContact(target *Contact, retChan chan []Contact) {
	//Locate k closest nodes
	shortlist := kademlia.Net.table.FindClosestContacts(target.ID, bucketSize) //3 närmsta grannarna
	var visitedNodes []Contact
	var closestNode Contact = shortlist[0]

	var alpha1 []Contact = chooseNContacts(shortlist, visitedNodes, numberOfParallelRequests)
	fmt.Println("FIND_NODE Shortlist: ", len(shortlist), "nContacts: ", alpha1)

	contactChan := make(chan []Contact, len(alpha1))
	removeChan := make(chan Contact, 1)

	for i, node := range alpha1 { //Alpha 1
		fmt.Println("FIND_NODE Looping neighbour: ", i, node.Address)
		go kademlia.Net.SendFindContactMessage(target, &node, contactChan, removeChan)
		visitedNodes = append(visitedNodes, node)
	}

	for i, _ := range alpha1 {
		fmt.Println("in loop alpha1")
		select {
		case recievedContacts := <-contactChan: //Recieved responses from findContactMessage
			for _, contact := range recievedContacts {
				//TODO: only add unique contacts
				shortlist = append(shortlist, contact)
				fmt.Println("Recieved contact: ", i, contact)
			}
		case removeContact := <-removeChan:
			fmt.Println("*********TIMEOUT alpha 1********")
			//Remove node from shortlist
			RemoveContact(shortlist, removeContact)
			break
		}
	}
	var madeProgress bool = true
	fmt.Println("Made Progress print", visitedNodes, shortlist)
	for madeProgress {
		madeProgress = false
		var alpha2 []Contact = chooseNContacts(shortlist, visitedNodes, numberOfParallelRequests)
		alpha2Channel := make(chan []Contact, len(alpha2))
		alpha2TimeoutChannel := make(chan Contact, 1)

		for j, node := range alpha2 { //Alpha 2
			fmt.Println("in loop alpha2", j, node.Address)
			if !kademlia.Net.table.me.ID.Equals(node.ID) {
				go kademlia.Net.SendFindContactMessage(target, &node, alpha2Channel, alpha2TimeoutChannel)
			}
			visitedNodes = append(visitedNodes, node)
		}
		fmt.Println("HERE!!")
	loop:
		for {
			fmt.Println("Before select")

			select {

			case recievedContacts := <-alpha2Channel:
				for _, contact := range recievedContacts {
					contact.CalcDistance(kademlia.Net.table.me.ID)
					fmt.Println("Appending alpha 2: ", closestNode)
					//TODO Only add unique
					shortlist = append(shortlist, contact)
					if contact.Less(&closestNode) {
						closestNode = contact
						madeProgress = true
					}

				}
			case removeContact := <-alpha2TimeoutChannel:
				fmt.Println("*********TIMEOUT alpha2********")
				RemoveContact(shortlist, removeContact)
				break loop
			}
		}
		if len(shortlist) >= bucketSize {
			fmt.Println("Shortlist size exceeded bucketsize")
			break
		}
	}
	//fmt.Println("shortlist: ", shortlist, "closest node: ", closestNode)

	retChan <- shortlist
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
	kadem := Kademlia{net, nil}
	return &kadem

}
func (kademlia *Kademlia) FindNode(target *Contact) []Contact {
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
	kadem := Kademlia{net, nil}

	knownID, err := kadem.Net.SendPingMessage(&knownContact)
	//_, err := SendPingMessage(&knownContact)
	//knownID := ""

	bootstrapContact := NewContact(NewKademliaID(knownID), knownIP+":"+strconv.Itoa(port))
	if err == nil {
		net.table.AddContact(bootstrapContact)
	}
	//fmt.Println("Known contact node: ", bootstrapContact)
	//net.SendFindContactMessage(&myContact, &knownContact, contactChan)

	closeNodes := kadem.FindNode(&kadem.Net.table.me)
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
		case "LookUpData":
			go kademlia.HandleFindData(m)
		case "StoreData":
			fmt.Println("StoreData RECIEVED:", m.msg)
			go kademlia.HandleStoreData(m.msg.Data, m.conn, m.remoteAddr)
		case "put":
			fmt.Println("Store data", m.msg)
			go kademlia.Store(m.msg.Data)
		case "get":
			fmt.Println("Get data")
		case "exit":
			//Kill network object
			fmt.Println("Quitting node...")
		}

	}
}
func (kademlia *Kademlia) Store(data []byte) {
	// TODO MAX Store
	returnChan := make(chan Message)

	storeContact := NewContact(NewKademliaID(Hash(data)), "")
	neighbours := kademlia.FindNode(&storeContact)

	fmt.Println("Neighbours for store", neighbours, "hash of data: ", Hash(data))
	for _, node := range neighbours {
		go kademlia.Net.SendStoreMessage(&node, data, returnChan)
	}
	returnMsg := <-returnChan
	fmt.Println("Store message: ", returnMsg)
}

func (kademlia *Kademlia) HandleStoreData(data []byte, conn net.UDPConn, retAddr net.UDPAddr) {
	//Spara data till fil med hash som namn.
	//Skicka tillbaka om det gick bra att spara. Terminerar om det ej går att spara.
	WriteToFile(data, Hash(data))
	kademlia.HashStorage = append(kademlia.HashStorage, Hash(data))
	conn.WriteToUDP([]byte("OK: Message stored"+kademlia.Net.table.me.Address), &retAddr)
}
func (kademlia *Kademlia) HandleFindNode(m InternalMessage) {
	resp := Message{
		Type:           "find-node-resp",
		ReturnContacts: kademlia.Net.table.FindClosestContacts(m.msg.TargetContact.ID, numberOfParallelRequests),
	}
	jsonMsg, err := json.Marshal(resp)
	fmt.Println("Handle FindNode")
	handleErr(err)
	m.conn.WriteToUDP(jsonMsg, &m.remoteAddr)
	fmt.Println("Returned closest contacts to target: ", m.msg.TargetContact.Address, "neighbours: ", resp.ReturnContacts)
}
func (kademlia *Kademlia) HandleFindData(m InternalMessage) {
	//DO I HAVE DATA
	files, _ := ioutil.ReadDir("../../D7024E/DATA")
	for _, file := range files {
		fmt.Println("FILE IN DATA: ", file.Name())
		if file.Name() == m.msg.TargetHash {
			//Send back data!!
			fmt.Println("MATCH!!!!!!!!!!!!!! TODO SEND BACK DATA")
			return
		}
	}
	resp := Message{
		Type:           "find-data-resp",
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
