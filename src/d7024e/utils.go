package d7024e

func RemoveContact(shortlist []Contact, contact Contact) []Contact {
	var result []Contact
	for _, node := range shortlist {
		if !node.ID.Equals(contact.ID) {
			result = append(result, node)
		}
	}
	return result
}
func StoreData(data []byte) {

}
