package kademlia

import (
	"testing"
	"time"
)

func TestProcessContactLookupReturns(t *testing.T) {
	targetID := NewKademliaID("FFFFFFFF0000000000000000000000000000000F")
	targetAdr := "172.0.0.10"
	targetContact := NewContact(targetID, targetAdr)

	kadTest := Init()
	kadTest.network.CreateNewLookupChannel(*targetContact.ID)
	dataChan := kadTest.network.lookupChs[*targetContact.ID]
	go dataspoofer(dataChan)

	returnedList := kadTest.ProcessContactLookupReturns(targetContact.ID)

	// These are the contacts we'll be placing in the reference list
	id1 := NewKademliaID("FFFFFFFF00000000000000000000FF0000000000")
	adr1 := "172.0.0.4"
	cont1 := NewContact(id1, adr1)
	id2 := NewKademliaID("FFFFFFFF0000000000000000000000FF00000000")
	adr2 := "172.0.0.5"
	cont2 := NewContact(id2, adr2)
	id3 := NewKademliaID("FFFFFFFF000000000000000000000000FF000000")
	adr3 := "172.0.0.6"
	cont3 := NewContact(id3, adr3)
	id4 := NewKademliaID("FFFFFFFF00000000000000000000000000FF0000")
	adr4 := "172.0.0.7"
	cont4 := NewContact(id4, adr4)
	id5 := NewKademliaID("FFFFFFFF0000000000000000000000000000FF00")
	adr5 := "172.0.0.8"
	cont5 := NewContact(id5, adr5)
	id6 := NewKademliaID("FFFFFFFF000000000000000000000000000000FF")
	adr6 := "172.0.0.9"
	cont6 := NewContact(id6, adr6)

	var referenceList []Contact
	referenceList = append(referenceList, cont6, cont5, cont4, cont3, cont2, cont1)

	for i := 0; i < len(returnedList); i++ {
		if returnedList[i].ID.String() != referenceList[i].ID.String() || returnedList[i].Address != referenceList[i].Address {
			// If address or ID doesn't match the reference list something is wrong
			t.Errorf("TestProcessContactLookupReturns failed, got %v, %v, expected %v, %v", returnedList[i].ID, returnedList[i].Address, referenceList[i].ID, referenceList[i].Address)
		}
	}
}

func dataspoofer(datachannel chan Contact) {
	id1 := NewKademliaID("FFFFFFFF00000000000000000000FF0000000000")
	adr1 := "172.0.0.4"
	cont1 := Contact{id1, adr1, id1}
	id2 := NewKademliaID("FFFFFFFF0000000000000000000000FF00000000")
	adr2 := "172.0.0.5"
	cont2 := Contact{id2, adr2, id2}
	id3 := NewKademliaID("FFFFFFFF000000000000000000000000FF000000")
	adr3 := "172.0.0.6"
	cont3 := Contact{id3, adr3, id3}
	datachannel <- cont1
	datachannel <- cont2
	datachannel <- cont3
	//^ Initial Contacts provided to the method
	time.Sleep(4 * time.Second)
	id4 := NewKademliaID("FFFFFFFF00000000000000000000000000FF0000")
	adr4 := "172.0.0.7"
	cont4 := Contact{id4, adr4, id4}
	id5 := NewKademliaID("FFFFFFFF0000000000000000000000000000FF00")
	adr5 := "172.0.0.8"
	cont5 := Contact{id5, adr5, id5}
	id6 := NewKademliaID("FFFFFFFF000000000000000000000000000000FF")
	adr6 := "172.0.0.9"
	cont6 := Contact{id6, adr6, id6}
	datachannel <- cont4
	datachannel <- cont5
	datachannel <- cont6
	//^ Second round of contacts provided to the method
	time.Sleep(2 * time.Second)
	datachannel <- cont4
	datachannel <- cont5
	datachannel <- cont6
	// resending the same messages again
}

func TestProcessDataLookupReturns(t *testing.T) {
	targetID := NewKademliaID("FFFFFFFF0000000000000000000000000000000F")
	targetAdr := "172.0.0.10"
	targetContact := NewContact(targetID, targetAdr)

	kadTest := Init()
	kadTest.network.CreateNewDataChannel(*targetContact.ID)
	dataChan := kadTest.network.dataChs[*targetContact.ID]
	kadTest.network.CreateNewLookupChannel(*targetContact.ID)
	lookupChan := kadTest.network.lookupChs[*targetContact.ID]
	go dataspoofer2(dataChan, lookupChan)

	returnedstring := kadTest.ProcessDataLookupReturns(targetContact.ID)
	expectedstring := "thisisateststring"

	if returnedstring != expectedstring {
		t.Errorf("TestProcessDataLookupReturns, got %v, expected %v", returnedstring, expectedstring)
	}
}

func dataspoofer2(datachannel chan string, lookupchannel chan Contact) {
	id1 := NewKademliaID("FFFFFFFF00000000000000000000FF0000000000")
	adr1 := "172.0.0.4"
	cont1 := Contact{id1, adr1, id1}
	id2 := NewKademliaID("FFFFFFFF0000000000000000000000FF00000000")
	adr2 := "172.0.0.5"
	cont2 := Contact{id2, adr2, id2}
	id3 := NewKademliaID("FFFFFFFF000000000000000000000000FF000000")
	adr3 := "172.0.0.6"
	cont3 := Contact{id3, adr3, id3}
	lookupchannel <- cont1
	lookupchannel <- cont2
	lookupchannel <- cont3
	//^ Initial Contacts provided to the method
	time.Sleep(4 * time.Second)
	id4 := NewKademliaID("FFFFFFFF00000000000000000000000000FF0000")
	adr4 := "172.0.0.7"
	cont4 := Contact{id4, adr4, id4}
	id5 := NewKademliaID("FFFFFFFF0000000000000000000000000000FF00")
	adr5 := "172.0.0.8"
	cont5 := Contact{id5, adr5, id5}
	id6 := NewKademliaID("FFFFFFFF000000000000000000000000000000FF")
	adr6 := "172.0.0.9"
	cont6 := Contact{id6, adr6, id6}
	lookupchannel <- cont4
	lookupchannel <- cont5
	lookupchannel <- cont6
	//^ Second round of contacts provided to the method
	time.Sleep(2 * time.Second)
	datachannel <- "thisisateststring"
	// resending the same messages again
}
