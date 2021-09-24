package d7024e

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"sort"
	"strconv"
	"time"
)

type Kademlia struct {
	Net Network
}

var numberOfParallelRequests int = 3

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
		knownContact := node
		go kademlia.Net.SendFindDataMessage(hash, &knownContact, contactChan, removeChan, dataChan)
		visitedNodes = append(visitedNodes, knownContact)
	}
	for i, _ := range alpha1 {
		fmt.Println("FIND_DATA in loop alpha1")
		select {
		case receivedContacts := <-contactChan: //received responses from findContactMessage
			for _, contact := range receivedContacts {
				shortlist = Add(shortlist, contact)
				fmt.Println("FIND_DATA received contact: ", i, contact)
			}
		case receivedData := <-dataChan:
			fmt.Println("========== Received in lookup data", receivedData)
			retDataChan <- receivedData
			return
		case removeContact := <-removeChan:
			fmt.Println("========= FIND_DATA TIMEOUT alpha 1 ========")
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
			knownContact := node
			go kademlia.Net.SendFindDataMessage(hash, &knownContact, alpha2Channel, alpha2TimeoutChannel, alpha2DataChannel)
			visitedNodes = append(visitedNodes, knownContact)
		}
	loop:
		for {
			select {
			case receivedContacts := <-alpha2Channel:
				for _, contact := range receivedContacts {
					contact.CalcDistance(kademlia.Net.table.me.ID)
					fmt.Println("Appending alpha 2: ", closestNode)
					shortlist = Add(shortlist, contact)
					if contact.Less(&closestNode) {
						closestNode = contact
						madeProgress = true
					}
				}
			case receivedData := <-alpha2DataChannel:
				retDataChan <- receivedData
				return
			case removeContact := <-alpha2TimeoutChannel:
				fmt.Println("========= TIMEOUT alpha2 ========")
				RemoveContact(shortlist, removeContact)
				break loop
			case <-time.After(timeoutDur):
				fmt.Println("*********TIMEOUT duration alpha 2********")
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

	contactChan := make(chan []Contact, len(alpha1))
	removeChan := make(chan Contact, 1)

	for _, node := range alpha1 { //Alpha 1

		knownContact := node
		go kademlia.Net.SendFindContactMessage(target, &knownContact, contactChan, removeChan)
		visitedNodes = append(visitedNodes, knownContact)
	}

	for range alpha1 {
		select {
		case receivedContacts := <-contactChan: //received responses from findContactMessage
			for _, contact := range receivedContacts {
				//TODO: only add unique contacts
				shortlist = Add(shortlist, contact)
			}
		case removeContact := <-removeChan:
			//Remove node from shortlist
			RemoveContact(shortlist, removeContact)
			break
		}
	}
	var madeProgress bool = true
	for madeProgress {
		madeProgress = false
		var alpha2 []Contact = chooseNContacts(shortlist, visitedNodes, numberOfParallelRequests)
		alpha2Channel := make(chan []Contact, len(alpha2))
		alpha2TimeoutChannel := make(chan Contact, 1)

		for _, node := range alpha2 { //Alpha 2
			knownContact := node
			go kademlia.Net.SendFindContactMessage(target, &knownContact, alpha2Channel, alpha2TimeoutChannel)
			visitedNodes = append(visitedNodes, knownContact)
		}
	loop:
		for {
			if len(alpha2) == 0 {
				break loop
			}
			select {

			case receivedContacts := <-alpha2Channel:
				for _, contact := range receivedContacts {
					contact.CalcDistance(kademlia.Net.table.me.ID)
					shortlist = Add(shortlist, contact)
					if contact.Less(&closestNode) {
						closestNode = contact
						madeProgress = true
					}
				}
			case removeContact := <-alpha2TimeoutChannel:
				RemoveContact(shortlist, removeContact)
				break loop

			case <-time.After(timeoutDur):
				break loop
			}
		}
		if len(shortlist) >= bucketSize {
			break
		}
	}
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
	kadem := Kademlia{net}
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
	kadem := Kademlia{net}

	knownMsg, err := kadem.Net.SendPingMessage(&knownContact)
	knownID := knownMsg.SenderContact.ID
	fmt.Println("Boostrap Contact: ", knownMsg.SenderContact.Address)
	bootstrapContact := NewContact(knownID, knownIP+":"+strconv.Itoa(port))
	if err == nil {
		net.table.AddContact(bootstrapContact)
	}
	closeNodes := kadem.FindNode(&kadem.Net.table.me)
	for _, node := range closeNodes {
		kadem.Net.table.AddContact(node)
	}
	fmt.Println("############## ------ Join network contact chan ------- ############", closeNodes)
	//Update buckets further away

	return &kadem
}
func (kademlia *Kademlia) updateBuckets() {
	/*
	* After this, the joining node refreshes all k-buckets further away than the k-bucket the bootstrap node falls in.
	* This refresh is just a lookup of a random key that is within that k-bucket range
	 */
}
func (kademlia *Kademlia) updateContact(senderContact *Contact) {
	//Use channels
	bucketIndex := kademlia.Net.table.getBucketIndex(senderContact.ID)
	fmt.Println("bucketIndex", bucketIndex)
	if kademlia.Net.table.buckets[bucketIndex].Len() < bucketSize {
		//Bucket is not full
		fmt.Println("Bucket not full")
		kademlia.Net.table.AddContact(*senderContact)
	} else {
		//Bucket is full
		//ping node at the head of bucket
		//if that fails to respons -> remove from list
		head := kademlia.Net.table.buckets[bucketIndex].list.Back()
		headContact := head.Value.(Contact)
		fmt.Println("headContact: ", headContact)
		pingResp, err := kademlia.Net.SendPingMessage(&headContact)
		fmt.Println("ping response: ", pingResp, err)
		if pingResp.Type == "ping" {
			//Ignore new contact

		} else {
			kademlia.Net.table.buckets[bucketIndex].list.Remove(head)          //Drop headContact
			kademlia.Net.table.buckets[bucketIndex].AddContact(*senderContact) //Add senderContact at tail

			//Drop headContact
			//Add senderContact at tail
		}

	}
}
func (kademlia *Kademlia) HandleMessage(msgChan chan InternalMessage) {
	for {
		var m = <-msgChan
		fmt.Println("Internal message received:", m.msg.Type)
		if !(m.msg.Type == "put" || m.msg.Type == "get" || m.msg.Type == "exit" || m.msg.Type == "test") {
			go kademlia.updateContact(&m.msg.SenderContact)
		}
		switch m.msg.Type {
		case "ping":
			go kademlia.Net.SendPingAckMessage(&m.conn, &m.remoteAddr)
		case "LookUpNode":
			go kademlia.HandleFindNode(m)
		case "LookUpData":
			go kademlia.HandleFindData(m)
		case "StoreData":
			fmt.Println("StoreData received:", m.msg)
			go kademlia.HandleStoreData(m.msg.Data, m.conn, m.remoteAddr)
		case "put":
			fmt.Println("Store data", m.msg)
			go kademlia.Store(m.msg.Data, m.conn, m.remoteAddr)
		case "get":
			fmt.Println("Get data")
			//go kademlia.HandleFindData(m)
			go kademlia.Locate(m.msg.TargetHash, m.conn, m.remoteAddr)
		case "exit":
			//Kill network object
			fmt.Println("Quitting node...")
			//os.Exit(66)
		case "test":
			cont := kademlia.Net.table.FindClosestContacts(NewKademliaID(Hash(m.msg.Data)), 6)
			msg, _ := json.Marshal(cont)
			fmt.Println("Hash: ", Hash(m.msg.Data))
			m.conn.WriteToUDP(msg, &m.remoteAddr)
		}

	}
}
func (kademlia *Kademlia) Store(data []byte, conn net.UDPConn, remoteAddr net.UDPAddr) {
	// TODO MAX Store
	returnChan := make(chan Message)

	storeContact := NewContact(NewKademliaID(Hash(data)), "")
	neighbours := kademlia.FindNode(&storeContact)

	fmt.Println("#######################Neighbours for store", neighbours, "hash of data: ", Hash(data))
	for _, node := range neighbours {
		knownContact := node
		go kademlia.Net.SendStoreMessage(&knownContact, data, returnChan)
	}
	returnMsg := <-returnChan
	fmt.Println("Store message: ", returnMsg)
	conn.WriteToUDP([]byte(returnMsg.TargetHash+"\n"), &remoteAddr)
}
func (kademlia *Kademlia) Locate(targetHash string, conn net.UDPConn, remoteAddr net.UDPAddr) {

	//kanaler
	returnData := make(chan string, 3)
	returnContact := make(chan []Contact, 3)
	kademlia.LookupData(targetHash, returnContact, returnData)
	select {
	case receivedContact := <-returnContact:
		fmt.Println("Received contact: ", receivedContact)
	case receivedData := <-returnData:
		fmt.Println("Received data: ", receivedData)
		conn.WriteToUDP([]byte(receivedData+"\n"), &remoteAddr)
	}
}

func (kademlia *Kademlia) HandleStoreData(data []byte, conn net.UDPConn, retAddr net.UDPAddr) {
	//Spara data till fil med hash som namn.
	//Skicka tillbaka om det gick bra att spara. Terminerar om det ej går att spara.
	WriteToFile(data, Hash(data))
	//kademlia.HashStorage = append(kademlia.HashStorage, Hash(data))
	m := Message{
		Type:          "store-ack",
		SenderContact: kademlia.Net.table.me,
		TargetHash:    Hash(data),
	}
	msg, _ := json.Marshal(m)

	conn.WriteToUDP(msg, &retAddr)
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
func (kademlia *Kademlia) HandleFindData(m InternalMessage) {
	//DO I HAVE DATA
	//fmt.Printf("GET MESSAGE: %+v\n", m, "\n")
	storagePath := "/app/filestorage"
	files, _ := ioutil.ReadDir(storagePath)
	fmt.Println("FILES DONE,", files)
	for _, file := range files {
		fmt.Println("FILE IN DATA: ", file.Name())
		if file.Name() == m.msg.TargetHash {
			content, err := ioutil.ReadFile(storagePath + "/" + m.msg.TargetHash)
			handleErr(err)
			resp := Message{
				Type: "find-data-resp",
				Data: content,
			}
			fmt.Println("MATCH!!!!!!!!!!!!!! ", string(resp.Data))
			jsonMsg, err := json.Marshal(resp)
			handleErr(err)
			m.conn.WriteToUDP(jsonMsg, &m.remoteAddr)
			return
		}
	}
	//fmt.Println("WTF GET MESSAGE,", hash)
	returnContacts := kademlia.Net.table.FindClosestContacts(NewKademliaID(m.msg.TargetHash), numberOfParallelRequests)
	//fmt.Println("RETURNCONTACTS IN GET MESSAGE ", returnContacts)
	resp := Message{
		Type:           "find-data-resp",
		ReturnContacts: returnContacts,
	}
	//fmt.Println("RESPONSE IN GET MESSAGE,", resp)
	jsonMsg, err := json.Marshal(resp)
	handleErr(err)
	m.conn.WriteToUDP(jsonMsg, &m.remoteAddr)
	fmt.Println("Returned closest contacts to target: ", m.msg.TargetContact.Address, "neighbours: ", resp.ReturnContacts)
}
