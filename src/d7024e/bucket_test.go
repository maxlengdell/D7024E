package d7024e

import (
	"container/list"
	"testing"

	"github.com/maxlengdell/D7024E/go1"
)

func indexOf(list *list.List, contact Contact) int {
	i := 0
	for e := list.Front(); e != nil; e = e.Next() {
		if e.Value == contact {
			return i
		}
		i++
	}
	return -1
}

func Test_Bucket_AddContact(t *testing.T) {
	contactOne := NewContact(one, "localhost:8001")
	contactTwo := NewContact(two, "localhost:8002")
	b := newBucket()
	go1.AssertEquals(t, 0, b.Len())
	b.AddContact(contactOne)
	go1.AssertEquals(t, 1, b.Len())
	// Adding again should be a NOP (when it's the only contact).
	b.AddContact(contactOne)
	go1.AssertEquals(t, 1, b.Len())

	// Adding a new contact should put it at the front.
	b.AddContact(contactTwo)
	go1.AssertEquals(t, 2, b.Len())
	go1.AssertEquals(t, 1, indexOf(b.list, contactOne))
	go1.AssertEquals(t, 0, indexOf(b.list, contactTwo))

	// Adding the old contact should MOVE it to the front.
	b.AddContact(contactOne)
	go1.AssertEquals(t, 2, b.Len())
	go1.AssertEquals(t, 0, indexOf(b.list, contactOne))
	go1.AssertEquals(t, 1, indexOf(b.list, contactTwo))
}

func Test_Bucket_GetContactAndCalcDistance(t *testing.T) {
	contactZero := NewContact(zero, "localhost:8000")
	contactOne := NewContact(one, "localhost:8001")
	contactTwo := NewContact(two, "localhost:8002")
	contactThree := NewContact(three, "localhost:8003")

	b := newBucket()
	b.AddContact(contactThree)
	b.AddContact(contactTwo)
	b.AddContact(contactOne)
	b.AddContact(contactZero)

	contacts := b.GetContactAndCalcDistance(three)
	go1.AssertEquals(t, three, contacts[0].distance)
	go1.AssertEquals(t, two, contacts[1].distance)
	go1.AssertEquals(t, one, contacts[2].distance)
	go1.AssertEquals(t, zero, contacts[3].distance)
}

func Test_BriefString_with_1_contact(t *testing.T) {
	bucket := newBucket()
	bucket.AddContact(NewContact(NewKademliaID("1a2b3c4d"), "1.2.3.4:123"))
	expected := "1a2b..0000"
	go1.AssertEquals(t, expected, bucket.BriefString())
}

func Test_BriefString_with_2_contacts(t *testing.T) {
	bucket := newBucket()
	bucket.AddContact(NewContact(NewKademliaID("1a2b3c4d"), "1.2.3.4:123"))
	bucket.AddContact(NewContact(NewKademliaID("2b3c4d5e"), "2.3.4.5:234"))
	expected := "2b3c..0000, 1a2b..0000"
	go1.AssertEquals(t, expected, bucket.BriefString())
}

func Test_String_with_1_contact(t *testing.T) {
	bucket := newBucket()
	bucket.AddContact(NewContact(NewKademliaID("1a2b3c4d"), "1.2.3.4:123"))
	expected := `contact("1a2b3c4d00000000000000000000000000000000", "1.2.3.4:123")`
	go1.AssertEquals(t, expected, bucket.String())
}

func Test_String_with_2_contact(t *testing.T) {
	bucket := newBucket()
	bucket.AddContact(NewContact(NewKademliaID("1a2b3c4d"), "1.2.3.4:123"))
	bucket.AddContact(NewContact(NewKademliaID("2b3c4d5e"), "2.3.4.5:234"))
	expected := `contact("2b3c4d5e00000000000000000000000000000000", "2.3.4.5:234"), contact("1a2b3c4d00000000000000000000000000000000", "1.2.3.4:123")`
	go1.AssertEquals(t, expected, bucket.String())
}
