package kademlia

import (
	"fmt"
	"testing"
	"time"
)


func TestNewBucket(t *testing.T) {
    bucket := newBucket()
    length := bucket.list.Len()
    if length != 0 {
        t.Errorf("Bucket was length %d instead of zero", length)
    }
}

func fillBucket(t *testing.T) (*bucket, [bucketSize]Contact, int) {
    bucket := newBucket()
    var contacts [bucketSize]Contact
    var it int
    for it = range len(contacts) {
        ip := "172.19.0." + fmt.Sprint(it)
        contact := NewContact(NewRandomKademliaID(), ip)
        contacts[it] = contact
        res, _ := bucket.AddContact(contact)
        if res != false {
            t.Errorf("Bucket filled up prematurely")
        }
    }
    
    return bucket, contacts, it
}


func TestAddContact(t *testing.T) {
    bucket, contacts, it := fillBucket(t)

    it++
    ip := "172.19.0." + fmt.Sprint(it)
    contact := NewContact(NewRandomKademliaID(), ip)
    res, oldContact := bucket.AddContact(contact)
    if res != true {
        t.Errorf("Bucket is not full")
    } else if *oldContact != contacts[0] {
        t.Error("Eldest contact was not returned")
    }


    // check that heartbeat works
    time.Sleep(heartbeatTimeout + 100*time.Millisecond)
    if bucket.list.Back().Value == *oldContact {
        t.Errorf("The old contact was not removed\n%v, %v", *(oldContact), bucket.list.Back().Value.(Contact))
    }
    if bucket.list.Front().Value != contact {
        t.Errorf("The new contact was not inserted after heartbeat\n%v, %v", bucket.list.Front().Value.(Contact), contact)
    }

    //check that there won't be duplicates
    contact = contacts[bucketSize-1]
    res, _ = bucket.AddContact(contact)
    if res != false {
        t.Error("Contact should already exist")
    }
    if bucket.list.Front().Value != contact {
        t.Error("The old contact wasn't moved to the top")
    }
}


func TestGetContactAndCalcDistance(t *testing.T) {
    bucket, _, _ := fillBucket(t)
    targetID := NewRandomKademliaID()
    res := bucket.GetContactAndCalcDistance(targetID)
    if bucket.list.Len() != len(res) {
        t.Errorf("The lists aren't the same length")
    }
    it := 0
    for contact := bucket.list.Front(); contact != nil; contact, it = contact.Next(), it+1 {
        cont1 := contact.Value.(Contact)
        cont2 := res[it]
        if cont1.ID != cont2.ID || cont1.Address != cont2.Address {
            t.Errorf("The lists aren't the same/in the same order\n%v\n%v", cont1, cont2)
        }
        if cont2.distance == nil {
            t.Error("The address wansn't properly added")
        }
    }
}
