package d7024e

import (
	"testing"

	"github.com/maxlengdell/D7024E/go1"
)

func Test_Contact_CalcDistance(t *testing.T) {
	contact := NewContact(three, "localhost:port")
	go1.AssertTrue(t, contact.distance == nil)
	contact.CalcDistance(one)
	go1.AssertEquals(t, two, contact.distance)
}

func Test_Contact_Less(t *testing.T) {
	contactOne := NewContact(one, "localhost:port")
	contactTwo := NewContact(two, "localhost:port")
	// Set their distances to "0...0"
	contactOne.CalcDistance(zero)
	contactTwo.CalcDistance(zero)
	// Then check who is closer.
	go1.AssertTrue(t, contactOne.Less(&contactTwo))
	go1.AssertFalse(t, contactTwo.Less(&contactOne))
}
