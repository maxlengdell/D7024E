package d7024e

import (
	"testing"
)

func Test_NewKademliaID_with_correct_length(t *testing.T) {
	str := "0123456789012345678901234567890123456789"
	id := NewKademliaID(str)
	AssertEquals(t, str, id.String())
}

// When the input is too short, the rest is 0 padded.
func Test_NewKademliaID_with_short_input(t *testing.T) {
	str := "0123456789"
	expected := "0123456789000000000000000000000000000000"
	id := NewKademliaID(str)
	AssertEquals(t, expected, id.String())
}

// When the input is too long, it is truncated.
func Test_NewKademliaID_with_long_input(t *testing.T) {
	str := "0123456789012345678901234567890123456789abc"
	expected := "0123456789012345678901234567890123456789"
	id := NewKademliaID(str)
	AssertEquals(t, expected, id.String())
}

// When the input is not a hex string, it is effectively
// empty and results in all 0.
func Test_NewKademliaID_with_invalid_data(t *testing.T) {
	str := "hello"	// Not hex.
	expected := "0000000000000000000000000000000000000000"
	id := NewKademliaID(str)
	AssertEquals(t, expected, id.String())
}

// If we generate 3 random IDs, they should all be different.
func Test_NewRandomKademliaID(t *testing.T) {
	id1 := NewRandomKademliaID()
	id2 := NewRandomKademliaID()
	id3 := NewRandomKademliaID()
	AssertNotEquals(t, id1, id2)
	AssertNotEquals(t, id2, id3)
	AssertNotEquals(t, id1, id3)
}

