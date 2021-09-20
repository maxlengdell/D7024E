package d7024e

import (
	"container/list"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

// contactExp is used for exporting Contacts - it mirrors Contact but doesn't
// use pointers.
type contactExp struct {
	ID       KademliaID
	Address  string
	Distance KademliaID
}

// toContact convers a contactExp to a Contact.
func (c contactExp) toContact() Contact {
	var distance *KademliaID
	zero := NewKademliaID("")
	if c.Distance != *zero {
		distance = &c.Distance
	}
	return Contact{
		ID:       &c.ID,
		Address:  c.Address,
		distance: distance,
	}
}

// fromContact convers a Contact to a contactExp.
func fromContact(c Contact) contactExp {
	var distance KademliaID
	if c.distance != nil {
		distance = *c.distance
	}
	return contactExp{
		ID:       *c.ID,
		Address:  c.Address,
		Distance: distance,
	}
}

// NodeState represents the state of the node.
// It holds the ID and address (ip:port) of the node as a Contact and
// the routing table as a map. The indices of non-empty buckets are
// mapped to a slice of contacts.
type NodeState struct {
	Me contactExp
	RT map[int][]contactExp
}

// ExportNodeState extracts the node state from a RoutingTable.
func ExportNodeState(rt *RoutingTable) NodeState {
	buckets := make(map[int][]contactExp)
	for i, b := range rt.buckets {
		if b.list.Len() <= 0 {
			//empty := make([]contactExp, 0)
			//buckets[i] = empty
			continue
		}
		buckets[i] = bucketToContacts(b)
	}
	me := fromContact(rt.me)
	return NodeState{Me: me, RT: buckets}
}

func bucketToContacts(b *bucket) []contactExp {
	var contacts []contactExp
	for e := b.list.Front(); e != nil; e = e.Next() {
		contact := e.Value.(Contact)
		exp := fromContact(contact)
		contacts = append(contacts, exp)
	}
	return contacts
}

// ImportNodeState recreates a RoutingTable from a NodeState.
func ImportNodeState(state NodeState) RoutingTable {
	var buckets [IDLength * 8]*bucket
	for i, contacts := range state.RT {
		buckets[i] = contactsToBucket(contacts)
	}
	for i := 0; i < len(buckets); i++ {
		if buckets[i] == nil {
			buckets[i] = &bucket{list.New()}
		}
	}
	return RoutingTable{state.Me.toContact(), buckets}
}

func contactsToBucket(contacts []contactExp) *bucket {
	list := list.New()
	for _, c := range contacts {
		list.PushBack(c.toContact())
	}
	return &bucket{list}
}

// SaveRoutingTable converts the routing table to a NodeState and writes its
// JSON-encoding to the given Writer.
func SaveRoutingTable(rt *RoutingTable, writer io.Writer) error {
	state := ExportNodeState(rt)
	jsonBytes, err := json.Marshal(state)
	if err != nil {
		return err
	}
	_, err = writer.Write(jsonBytes)
	return err
}

// SaveRoutingTableToFile works like SaveRoutingTable and uses a file
// to write to.
func SaveRoutingTableToFile(rt *RoutingTable, filename string) error {
	state := ExportNodeState(rt)
	jsonBytes, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, jsonBytes, 0666)
}

// LoadRoutingTable reads a JSON-encoded NodeState from the given Reader and
// converts it to a routing table.
func LoadRoutingTable(reader io.Reader) (*RoutingTable, error) {
	jsonBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var state NodeState
	err = json.Unmarshal(jsonBytes, &state)
	if err != nil {
		return nil, err
	}
	rt := ImportNodeState(state)
	return &rt, err
}

// LoadRoutingTableFromFile works like LoadRoutingTable and uses a file
// to read from.
func LoadRoutingTableFromFile(filename string) (*RoutingTable, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return LoadRoutingTable(file)
}
