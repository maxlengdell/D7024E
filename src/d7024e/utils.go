package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

var filepath string = "/app/filestorage/"

func RemoveContact(shortlist []Contact, contact Contact) []Contact {
	var result []Contact
	for _, node := range shortlist {
		if !node.ID.Equals(contact.ID) {
			result = append(result, node)
		}
	}
	return result
}

func WriteToFile(data []byte, filename string) {
	err := ioutil.WriteFile(filepath+filename, data, 0644)
	handleErr(err)
}

// TODO: Find should probably be renamed to IndexOf.

// Find returns the first index where the given contact is stored or -1
// if it cannot be found.
func Find(slice []Contact, val Contact) (int, bool) {
	for i, item := range slice {
		if item.ID.Equals(val.ID) {
			return i, true
		}
	}
	return -1, false
}

func Trace(fname string, args ...interface{}) {
	defer trace(fname, args...)()
}

func trace(fname string, args ...interface{}) func() {
	start := time.Now()
	s := fmt.Sprint(args...)
	log.Printf("Entering %s(%s)", fname, s)
	return func() {
		end := time.Since(start)
		log.Printf("Exited %s(%s) after (%.3fs)\n", fname, s, end.Seconds())
	}
}

// Add a contact to the slice of contacts if one with matching ID does not
// already exist. Note that only the contacts' Kademlia IDs are used for
// comparison.
func Add(contacts []Contact, val Contact) []Contact {
	for _, item := range contacts {
		if item.ID.Equals(val.ID) {
			return contacts
		}
	}
	return append(contacts, val)
}

func Hash(data []byte) string {
	//Hash data to sha1 and return
	sh := sha1.Sum(data)
	return hex.EncodeToString(sh[:])
}

// Nothing represents a void value.
type Nothing struct{}

// Any represents any value.
type Any interface{}

// VoidPromise represents an async computation that does not return a result
// and that we can wait on to complete later.
type VoidPromise <-chan Nothing

// AnyPromise represents an async computation that returns some result
// and that we can wait on to complete later and get the result.
type AnyPromise <-chan Any

// GoVoid runs a (void) function in a new goroutine and (immediately) returns
// a promise.
func GoVoid(f func()) VoidPromise {
	ch := make(chan Nothing)
	go func() {
		defer close(ch)
		f()
		ch <- Nothing{}
	}()
	return ch
}

// GoVoid runs a function returning some result in a new  goroutine and
// (immediately) returns a promise.
func GoAny(f func() Any) AnyPromise {
	ch := make(chan Any)
	go func() {
		defer close(ch)
		ch <- f()
	}()
	return ch
}

// WaitFor waits for the promise to complete or until the timeout is reached.
func (promise VoidPromise) WaitFor(timeout time.Duration) (ok bool) {
	select {
	case <-promise:
		return true
	case <-time.After(timeout):
		return false
	}
}

// WaitFor waits for the promise to complete or until the timeout is reached.
func (promise AnyPromise) WaitFor(timeout time.Duration) (val Any, ok bool) {
	select {
	case val := <-promise:
		return val, true
	case <-time.After(timeout):
		return nil, false
	}
}

// MillisecondsDuration converts a duration in milliseconds to a Duration.
func MillisecondsDuration(ms int64) time.Duration {
	return time.Duration(ms * 1e6)
}

// SecondsDuration converts a duration in seconds to a Duration.
func SecondsDuration(seconds int64) time.Duration {
	return time.Duration(seconds * 1e9)
}
