package d7024e

import (
	"fmt"
	"testing"

	"github.com/maxlengdell/D7024E/go1"
)

// TODO: remove? Doesn't really test anything?
func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
}
func TestRoutingTable2(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))

	contacts := rt.FindClosestContacts(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println("MAX: ", contacts[i].String())
	}
}

func Test_RoutingTable_String(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("1a2b3c4d"), "1.2.3.4:123"))

	rt.AddContact(NewContact(NewKademliaID("000000000000000000000000000000002b3c4d5e"), "5.4.3.2:234"))
	rt.AddContact(NewContact(NewKademliaID("2b3c4d5e00000000000000000000000000000000"), "2.3.4.5:234"))
	rt.AddContact(NewContact(NewKademliaID("3b3c4d5e00000000000000000000000000000000"), "3.4.5.6:345"))
	expected := `  bucket[  2]: contact("3b3c4d5e00000000000000000000000000000000", "3.4.5.6:345"), contact("2b3c4d5e00000000000000000000000000000000", "2.3.4.5:234")  bucket[  3]: contact("000000000000000000000000000000002b3c4d5e", "5.4.3.2:234")`
	go1.AssertEquals(t, expected, rt.String())

}

func Example_RoutingTable_ShowBucketSizes() {
	rt := NewRoutingTable(NewContact(NewKademliaID("1a2b3c4d"), "1.2.3.4:123"))

	rt.AddContact(NewContact(NewKademliaID("000000000000000000000000000000002b3c4d5e"), "5.4.3.2:234"))
	rt.AddContact(NewContact(NewKademliaID("2b3c4d5e00000000000000000000000000000000"), "2.3.4.5:234"))
	rt.AddContact(NewContact(NewKademliaID("3b3c4d5e00000000000000000000000000000000"), "3.4.5.6:345"))
	rt.ShowBucketSizes()
	// Output: bucket[  2]: 2  bucket[  3]: 1
}

func Example_RoutingTable_ShowFullBucketContents() {
	rt := NewRoutingTable(NewContact(NewKademliaID("1a2b3c4d"), "1.2.3.4:123"))

	rt.AddContact(NewContact(NewKademliaID("000000000000000000000000000000002b3c4d5e"), "5.4.3.2:234"))
	rt.AddContact(NewContact(NewKademliaID("2b3c4d5e00000000000000000000000000000000"), "2.3.4.5:234"))
	rt.AddContact(NewContact(NewKademliaID("3b3c4d5e00000000000000000000000000000000"), "3.4.5.6:345"))
	rt.ShowFullBucketContents()
	// Output: bucket[  2]: contact("3b3c4d5e00000000000000000000000000000000", "3.4.5.6:345"), contact("2b3c4d5e00000000000000000000000000000000", "2.3.4.5:234")  bucket[  3]: contact("000000000000000000000000000000002b3c4d5e", "5.4.3.2:234")
}

func Example_RoutingTable_ShowBriefBucketContents() {
	rt := NewRoutingTable(NewContact(NewKademliaID("1a2b3c4d"), "1.2.3.4:123"))

	rt.AddContact(NewContact(NewKademliaID("000000000000000000000000000000002b3c4d5e"), "5.4.3.2:234"))
	rt.AddContact(NewContact(NewKademliaID("2b3c4d5e00000000000000000000000000000000"), "2.3.4.5:234"))
	rt.AddContact(NewContact(NewKademliaID("3b3c4d5e00000000000000000000000000000000"), "3.4.5.6:345"))
	rt.ShowBriefBucketContents()
	// Output: bucket[  2]: 3b3c..0000, 2b3c..0000  bucket[  3]: 0000..4d5e
}
