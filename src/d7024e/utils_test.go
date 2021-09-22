package d7024e

import (
	"testing"
	"time"

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

func slowVoidOp(ms int64) func() {
	return func() {
		time.Sleep(MillisecondsDuration(ms))
	}
}

func slowAnyOp(value int, ms int64) func() Any {
	return func() Any {
		time.Sleep(MillisecondsDuration(ms))
		return value
	}
}

func Test_VoidPromise_should_return_immediately(t *testing.T) {
	op := slowVoidOp(9000) // op takes 9 seconds
	start := time.Now()
	GoVoid(op) // should return immediately
	go1.AssertTrue(t, time.Since(start).Milliseconds() < 10)
}

func Test_VoidPromise_WaitUntil_ok(t *testing.T) {
	op := slowVoidOp(500) // op takes 0.5 seconds
	start := time.Now()
	promise := GoVoid(op)
	// Op should finish before the 1 second timeout.
	ok := promise.WaitFor(MillisecondsDuration(1000))
	end := time.Since(start).Milliseconds()
	go1.AssertTrue(t, ok)
	// And it should take roughly 500 ms.
	go1.AssertTrue(t, end > 400)
	go1.AssertTrue(t, end < 600)
}

func Test_VoidPromise_WaitUntil_timeout(t *testing.T) {
	op := slowVoidOp(500) // op takes 0.5 seconds
	start := time.Now()
	promise := GoVoid(op)
	// Op should time out after 0.2 s (before op finishes).
	ok := promise.WaitFor(MillisecondsDuration(200))
	end := time.Since(start).Milliseconds()
	go1.AssertFalse(t, ok)
	// And it should take roughly 200 ms.
	go1.AssertTrue(t, end > 100)
	go1.AssertTrue(t, end < 300)
}

func Test_AnyPromise_should_return_immediately(t *testing.T) {
	op := slowAnyOp(7, 9000) // op takes 9 seconds
	start := time.Now()
	GoAny(op) // should return immediately
	go1.AssertTrue(t, time.Since(start).Milliseconds() < 10)
}

func Test_AnyPromise_WaitUntil_ok(t *testing.T) {
	op := slowAnyOp(7, 500) // op takes 0.5 seconds
	start := time.Now()
	promise := GoAny(op)
	// Op should finish before the 1 second timeout.
	val, ok := promise.WaitFor(MillisecondsDuration(1000))
	end := time.Since(start).Milliseconds()
	go1.AssertTrue(t, ok)
	go1.AssertEquals(t, 7, val)
	// And it should take roughly 500 ms.
	go1.AssertTrue(t, end > 400)
	go1.AssertTrue(t, end < 600)
}

func Test_AnyPromise_WaitUntil_timeout(t *testing.T) {
	op := slowAnyOp(7, 500) // op takes 0.5 seconds
	start := time.Now()
	promise := GoAny(op)
	// Op should time out after 0.2 s (before op finishes).
	val, ok := promise.WaitFor(MillisecondsDuration(200))
	end := time.Since(start).Milliseconds()
	go1.AssertFalse(t, ok)
	go1.AssertEquals(t, nil, val)
	// And it should take roughly 200 ms.
	go1.AssertTrue(t, end > 100)
	go1.AssertTrue(t, end < 300)
}

func Test_SecondsDuration(t *testing.T) {
	dur := SecondsDuration(3)
	go1.AssertEquals(t, time.Duration(3e9), dur)
	go1.AssertEquals(t, time.Duration(3e9).Seconds(), dur.Seconds())
	go1.AssertEquals(t, 3.0, float64(dur.Seconds()))
}

func Test_MillisecondsDuration(t *testing.T) {
	dur := MillisecondsDuration(500)
	go1.AssertEquals(t, time.Duration(500e6), dur)
	go1.AssertEquals(t, time.Duration(500e6).Seconds(), dur.Seconds())
	// 0.5 is an exact power of 2, so the comparison is safe.
	go1.AssertEquals(t, 0.5, float64(dur.Seconds()))
}
