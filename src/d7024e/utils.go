package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
	"os"
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
	err := os.WriteFile(filepath+filename, data, 0644)
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
func Add(slice []Contact, val Contact) []Contact {
	var returnArr []Contact
	_, check := Find(slice, val)
	if !check {
		returnArr = append(slice, val)
		return returnArr
	} else {
		return slice
	}
}
func Hash(data []byte) string {
	//Hash data to sha1 and return
	sh := sha1.Sum(data)
	return hex.EncodeToString(sh[:])
}
