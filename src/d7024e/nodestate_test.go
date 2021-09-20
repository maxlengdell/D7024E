package d7024e

import (
	"fmt"
)

var id = NewKademliaID("abc")

// This is the same test case as in routingtable_test, except that the routing
// table is first exported to NodeState and then imported again.
// The point is that it still produces the same output as the other test.
func Example() {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	err := SaveRoutingTableToFile(rt, "rt.json")
	if err != nil {
		fmt.Printf("error while saving: %v\n", err)
	}
	rt2, err := LoadRoutingTableFromFile("rt.json")
	if err != nil {
		fmt.Printf("error while loading: %v\n", err)
	}

	// Note: uses the imported routing table rt2.
	contacts := rt2.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
	// Output:
	// contact("2111111400000000000000000000000000000000", "localhost:8002")
	// contact("1111111400000000000000000000000000000000", "localhost:8002")
	// contact("1111111100000000000000000000000000000000", "localhost:8002")
	// contact("1111111200000000000000000000000000000000", "localhost:8002")
	// contact("1111111300000000000000000000000000000000", "localhost:8002")
	// contact("ffffffff00000000000000000000000000000000", "localhost:8001")
}
