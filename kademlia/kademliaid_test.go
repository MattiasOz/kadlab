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

func TestDistance(t *testing.T) {
	testID1 := NewKademliaID("FFFFFFFF0000000000000000000000000000000F")
	testID2 := NewKademliaID("FFFFFFFF000000000000000000000000000000FF")
	testID3 := NewKademliaID("FFFFFFFF000000000000000000000000000000F0")

	dist1 := testID1.CalcDistance(*testID1) // This should be zero
	dist2 := testID1.CalcDistance(*testID2) // This should be 0x000...0f0
	dist3 := testID1.CalcDistance(*testID3) // This should be 0x000...0ff

	correctDist1 := NewKademliaID("0000000000000000000000000000000000000000")
	correctDist2 := NewKademliaID("00000000000000000000000000000000000000f0")
	correctDist3 := NewKademliaID("00000000000000000000000000000000000000ff")

	if dist1.String() != correctDist1.String() {
		t.Errorf("TestEquals failed, got %v, expected %v", dist1, correctDist1)
	}
	if dist2.String() != correctDist2.String() {
		t.Errorf("TestEquals failed, got %v, expected %v", dist2, correctDist2)
	}
	if dist3.String() != correctDist3.String() {
		t.Errorf("TestEquals failed, got %v, expected %v", dist3, correctDist3)
	}
}
