package kademlia

import (
	"fmt"
	"testing"
)

// FIXME: This test doesn't actually test anything. There is only one assertion
// that is included as an example.

func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}

	// TODO: This is just an example. Make more meaningful assertions.
	if len(contacts) != 6 {
		t.Fatalf("Expected 6 contacts but instead got %d", len(contacts))
	}
}

func TestNewRoutingTable(t *testing.T) {
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	rt := NewRoutingTable(me)
	if me != rt.me {
		t.Errorf("TestNewRoutingTable failed, local contact is %v, expected %v", rt.me, me)
	}
}

func TestFindClosestContacts(t *testing.T) {
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	rt := NewRoutingTable(me)
	cont1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000001"), "localhost:8001")
	cont2 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000010"), "localhost:8002")
	cont3 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000100"), "localhost:8002")
	cont4 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000001000"), "localhost:8002")
	cont5 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000010000"), "localhost:8002")
	cont6 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000100000"), "localhost:8002")

	rt.AddContact(cont1)
	rt.AddContact(cont2)
	rt.AddContact(cont3)
	rt.AddContact(cont4)
	rt.AddContact(cont5)
	rt.AddContact(cont6)

	var referenceList1 []Contact
	referenceList1 = append(referenceList1, cont1, cont2, cont3, cont4, cont5, cont6)
	var referenceList2 []Contact
	referenceList2 = append(referenceList2, cont6, cont1, cont2, cont3, cont4, cont5)

	realList1 := rt.FindClosestContacts(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), 6)
	realList2 := rt.FindClosestContacts(NewKademliaID("FFFFFFFF00000000000000000000000000100000"), 6)

	for i := 0; i < 6; i++ {
		if referenceList1[i].ID != realList1[i].ID || referenceList2[i].ID != realList2[i].ID {
			t.Errorf("TestFindClosestContacts failed in step %v, expected1 is %v, got %v, expected2 is %v, got %v", i, realList1[i].ID, referenceList1[i].ID, realList2[i].ID, referenceList2[i].ID)
		}
	}
}
