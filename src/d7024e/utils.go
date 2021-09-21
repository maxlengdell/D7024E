package d7024e

import (
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
