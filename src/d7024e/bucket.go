package d7024e

import (
	"fmt"
	"container/list"
)

// bucket definition
// contains a List
type bucket struct {
	list *list.List
}

// String converts a bucket to a string by showing the IDs of its contacts.
func (b *bucket) String() string {
	var str string
	e := b.list.Front()
	contact := e.Value.(Contact)
	str += contact.String()
	for e = e.Next(); e != nil; e = e.Next() {
		contact := e.Value.(Contact)
		str += fmt.Sprintf(", %s", contact.String())
	}
	return str
}

func (b *bucket) BriefString() string {
	var str string
	e := b.list.Front()
	contact := e.Value.(Contact)
	idStr := contact.ID.String()
	str += fmt.Sprintf("%s..%s", idStr[0:4], idStr[len(idStr)-4:])
	for e = e.Next(); e != nil; e = e.Next() {
		contact := e.Value.(Contact)
		idStr := contact.ID.String()
		str += fmt.Sprintf(", %s..%s", idStr[0:4], idStr[len(idStr)-4:])
	}
	return str
}

// newBucket returns a new instance of a bucket
func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *bucket) AddContact(contact Contact) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
		}
	} else {
		bucket.list.MoveToFront(element)
	}
}

// GetContactAndCalcDistance returns an array of Contacts where 
// the distance has already been calculated
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

// Len return the size of the bucket
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
