package d7024e

import (
	"testing"

	"github.com/maxlengdell/D7024E/go1"
)

var zero *KademliaID = NewKademliaID("0000000000000000000000000000000000000000")
var one *KademliaID = NewKademliaID("0000000000000000000000000000000000000001")
var two *KademliaID = NewKademliaID("0000000000000000000000000000000000000002")
var three *KademliaID = NewKademliaID("0000000000000000000000000000000000000003")
var oneL *KademliaID = NewKademliaID("1000000000000000000000000000000000000000")
var twoL *KademliaID = NewKademliaID("2000000000000000000000000000000000000000")
var threeL *KademliaID = NewKademliaID("3000000000000000000000000000000000000000")
var ones *KademliaID = NewKademliaID("ffffffffffffffffffffffffffffffffffffffff")

func Test_NewKademliaID_with_correct_length(t *testing.T) {
	str := "0123456789012345678901234567890123456789"
	id := NewKademliaID(str)
	go1.AssertEquals(t, str, id.String())
}

// When the input is too short, the rest is 0 padded.
func Test_NewKademliaID_with_short_input(t *testing.T) {
	str := "0123456789"
	expected := "0123456789000000000000000000000000000000"
	id := NewKademliaID(str)
	go1.AssertEquals(t, expected, id.String())
}

// When the input is too long, it is truncated.
func Test_NewKademliaID_with_long_input(t *testing.T) {
	str := "0123456789012345678901234567890123456789abc"
	expected := "0123456789012345678901234567890123456789"
	id := NewKademliaID(str)
	go1.AssertEquals(t, expected, id.String())
}

// When the input is not a hex string, it is effectively
// empty and results in all 0.
func Test_NewKademliaID_with_invalid_data(t *testing.T) {
	str := "hello" // Not hex.
	expected := "0000000000000000000000000000000000000000"
	id := NewKademliaID(str)
	go1.AssertEquals(t, expected, id.String())
}

// If we generate 3 random IDs, they should all be different.
func Test_NewRandomKademliaID(t *testing.T) {
	id1 := NewRandomKademliaID()
	id2 := NewRandomKademliaID()
	id3 := NewRandomKademliaID()
	go1.AssertNotEquals(t, id1, id2)
	go1.AssertNotEquals(t, id2, id3)
	go1.AssertNotEquals(t, id1, id3)
}

func Test_KademliaID_CalcDistance(t *testing.T) {
	distance := zero.CalcDistance(zero)
	go1.AssertEquals(t, zero, distance)
	distance = one.CalcDistance(one)
	go1.AssertEquals(t, zero, distance)
	distance = two.CalcDistance(two)
	go1.AssertEquals(t, zero, distance)
	distance = three.CalcDistance(three)
	go1.AssertEquals(t, zero, distance)
	distance = oneL.CalcDistance(oneL)
	go1.AssertEquals(t, zero, distance)
	distance = twoL.CalcDistance(twoL)
	go1.AssertEquals(t, zero, distance)
	distance = threeL.CalcDistance(threeL)
	go1.AssertEquals(t, zero, distance)
	distance = ones.CalcDistance(ones)
	go1.AssertEquals(t, zero, distance)

	distance = zero.CalcDistance(one)
	go1.AssertEquals(t, one, distance)

	distance = one.CalcDistance(three)
	go1.AssertEquals(t, two, distance)

	distance = two.CalcDistance(one)
	go1.AssertEquals(t, three, distance)

	distance = oneL.CalcDistance(threeL)
	go1.AssertEquals(t, twoL, distance)
}

func Test_KademliaID_Equals(t *testing.T) {
	go1.AssertEquals(t, zero, zero)
	go1.AssertEquals(t, one, one)
	go1.AssertEquals(t, two, two)
	go1.AssertEquals(t, three, three)
	go1.AssertEquals(t, oneL, oneL)
	go1.AssertEquals(t, twoL, twoL)
	go1.AssertEquals(t, threeL, threeL)
	go1.AssertEquals(t, ones, ones)
	go1.AssertNotEquals(t, zero, one)
	go1.AssertNotEquals(t, zero, two)
	go1.AssertNotEquals(t, zero, three)
	go1.AssertNotEquals(t, zero, oneL)
	go1.AssertNotEquals(t, zero, twoL)
	go1.AssertNotEquals(t, zero, threeL)
	go1.AssertNotEquals(t, oneL, twoL)
	go1.AssertNotEquals(t, oneL, threeL)
	go1.AssertNotEquals(t, twoL, threeL)
	go1.AssertNotEquals(t, ones, one)
	go1.AssertNotEquals(t, ones, two)
	go1.AssertNotEquals(t, ones, three)
	go1.AssertNotEquals(t, ones, oneL)
	go1.AssertNotEquals(t, ones, twoL)
	go1.AssertNotEquals(t, ones, threeL)
}

func TestLess(t *testing.T) {
	go1.AssertTrue(t, zero.Less(one))
	go1.AssertTrue(t, zero.Less(two))
	go1.AssertTrue(t, zero.Less(three))
	go1.AssertTrue(t, zero.Less(oneL))
	go1.AssertTrue(t, zero.Less(twoL))
	go1.AssertTrue(t, zero.Less(threeL))
	go1.AssertTrue(t, zero.Less(ones))

	go1.AssertTrue(t, one.Less(two))
	go1.AssertTrue(t, one.Less(three))
	go1.AssertTrue(t, one.Less(oneL))
	go1.AssertTrue(t, one.Less(twoL))
	go1.AssertTrue(t, one.Less(threeL))
	go1.AssertTrue(t, one.Less(ones))

	go1.AssertTrue(t, two.Less(three))
	go1.AssertTrue(t, two.Less(oneL))
	go1.AssertTrue(t, two.Less(twoL))
	go1.AssertTrue(t, two.Less(threeL))
	go1.AssertTrue(t, two.Less(ones))

	go1.AssertTrue(t, oneL.Less(twoL))
	go1.AssertTrue(t, oneL.Less(threeL))
	go1.AssertTrue(t, oneL.Less(ones))

	go1.AssertTrue(t, twoL.Less(threeL))
	go1.AssertTrue(t, twoL.Less(ones))

	go1.AssertTrue(t, threeL.Less(ones))

	go1.AssertFalse(t, one.Less(zero))
	go1.AssertFalse(t, two.Less(zero))
	go1.AssertFalse(t, three.Less(zero))
	go1.AssertFalse(t, oneL.Less(zero))
	go1.AssertFalse(t, twoL.Less(zero))
	go1.AssertFalse(t, threeL.Less(zero))
	go1.AssertFalse(t, ones.Less(zero))

	go1.AssertFalse(t, two.Less(one))
	go1.AssertFalse(t, three.Less(two))
	go1.AssertFalse(t, twoL.Less(oneL))
	go1.AssertFalse(t, threeL.Less(twoL))
}

func Test_KademliaID_String(t *testing.T) {
	str := "0123456789abcdefedcba9876543210123456789"
	id := NewKademliaID(str)
	go1.AssertEquals(t, str, id.String())
}

func Test_Store(t *testing.T) {

}
