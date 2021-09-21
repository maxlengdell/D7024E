package d7024e

import (
	"testing"

	"github.com/maxlengdell/D7024E/go1"
)

var contactZero Contact = NewContact(zero, "localhost:port")
var contactOne Contact = NewContact(one, "localhost:port")
var contactTwo Contact = NewContact(two, "localhost:port")
var contactThree Contact = NewContact(three, "localhost:port")

func Test_Find_empty(t *testing.T) {
	var contacts []Contact
	i, ok := Find(contacts, contactOne)
	go1.AssertEquals(t, -1, i)
	go1.AssertFalse(t, ok)
}

func Test_Find_single_match(t *testing.T) {
	contacts := []Contact{contactOne}
	i, ok := Find(contacts, contactOne)
	go1.AssertEquals(t, 0, i)
	go1.AssertTrue(t, ok)
}

func Test_Find_single_no_match(t *testing.T) {
	contacts := []Contact{contactOne}
	i, ok := Find(contacts, contactTwo)
	go1.AssertEquals(t, -1, i)
	go1.AssertFalse(t, ok)
}

func Test_Find_multi_match(t *testing.T) {
	contacts := []Contact{contactOne, contactTwo, contactThree, contactTwo}
	i, ok := Find(contacts, contactTwo)
	go1.AssertEquals(t, 1, i)
	go1.AssertTrue(t, ok)
}

func Test_Add_empty(t *testing.T) {
	contacts := []Contact{}
	contacts = Add(contacts, contactZero)
	expected := []Contact{contactZero}
	go1.AssertEquals(t, expected, contacts)
}

func Test_Add_new(t *testing.T) {
	contacts := []Contact{contactOne, contactTwo, contactThree}
	contacts = Add(contacts, contactZero)
	expected := []Contact{contactOne, contactTwo, contactThree, contactZero}
	go1.AssertEquals(t, expected, contacts)
}

func Test_Add_existing(t *testing.T) {
	contacts := []Contact{contactOne, contactTwo, contactThree}
	contacts = Add(contacts, contactTwo)
	expected := []Contact{contactOne, contactTwo, contactThree}
	go1.AssertEquals(t, expected, contacts)
}
