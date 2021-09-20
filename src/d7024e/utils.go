package d7024e

import "os"

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
