package kademlia

import (
	"testing"
)

func TestNewKademliaID(t *testing.T) {
	testID := NewKademliaID("FFFFFFFF0000000000000000000000000000000F")
	correctID := "ffffffff0000000000000000000000000000000f"
	// Converting to kademliaID and then back to string makes letters lowercase, but that is expected behav.
	if testID.String() != correctID {
		t.Errorf("TestNewKademliaID failed, got id %v, expected id %v", testID.String(), correctID)
	}
}

func TestLess(t *testing.T) {
	testID1 := NewKademliaID("FFFFFFFF0000000000000000000000000000000F")
	testID2 := NewKademliaID("FFFFFFFF000000000000000000000000000000FF")

	lessTrue := testID1.Less(testID2)   // This should be true
	lessFalse := testID2.Less(testID1)  // This should be false
	lessFalse2 := testID2.Less(testID2) // This should be false

	if !lessTrue || lessFalse || lessFalse2 { //lesstrue is false and lessfalse is true
		t.Errorf("TestLess failed, got %v, %v, %v - expected True, False, False", lessTrue, lessFalse, lessFalse2)
	}
}

func TestEquals(t *testing.T) {
	testID1 := NewKademliaID("FFFFFFFF0000000000000000000000000000000F")
	testID2 := NewKademliaID("FFFFFFFF000000000000000000000000000000FF")

	equalTrue := testID1.Equals(testID1)  // This should be true
	equalFalse := testID2.Equals(testID1) // This should be false

	if !equalTrue || equalFalse { //lesstrue is false and lessfalse is true
		t.Errorf("TestEquals failed, got %v, %v - expected True, False", equalTrue, equalFalse)
	}
}
