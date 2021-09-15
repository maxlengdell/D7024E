package d7024e

import (
	"testing"
)

func Test_BriefString_with_1_contact(t *testing.T) {
	bucket := newBucket()
	bucket.AddContact(NewContact(NewKademliaID("1a2b3c4d"), "1.2.3.4:123"))
	expected := "1a2b..0000"
	AssertEquals(t, expected, bucket.BriefString())
}

func Test_BriefString_with_2_contacts(t *testing.T) {
	bucket := newBucket()
	bucket.AddContact(NewContact(NewKademliaID("1a2b3c4d"), "1.2.3.4:123"))
	bucket.AddContact(NewContact(NewKademliaID("2b3c4d5e"), "2.3.4.5:234"))
	expected := "2b3c..0000, 1a2b..0000"
	AssertEquals(t, expected, bucket.BriefString())
}

func Test_String_with_1_contact(t *testing.T) {
	bucket := newBucket()
	bucket.AddContact(NewContact(NewKademliaID("1a2b3c4d"), "1.2.3.4:123"))
	expected := `contact("1a2b3c4d00000000000000000000000000000000", "1.2.3.4:123")`
	AssertEquals(t, expected, bucket.String())
}

func Test_String_with_2_contact(t *testing.T) {
	bucket := newBucket()
	bucket.AddContact(NewContact(NewKademliaID("1a2b3c4d"), "1.2.3.4:123"))
	bucket.AddContact(NewContact(NewKademliaID("2b3c4d5e"), "2.3.4.5:234"))
	expected := `contact("2b3c4d5e00000000000000000000000000000000", "2.3.4.5:234"), contact("1a2b3c4d00000000000000000000000000000000", "1.2.3.4:123")`
	AssertEquals(t, expected, bucket.String())
}
