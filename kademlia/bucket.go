package kademlia

import (
	"container/list"
	"time"
)

const heartbeatTimeout = 3 * time.Second

// bucket definition
// contains a List
type bucket struct {
	list *list.List
}

// newBucket returns a new instance of a bucket
func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *bucket) AddContact(contact Contact) (bool, *Contact) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

    // if contact is not already in bucket
	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
		} else {
			// Check that the last node is alive and add the current if it's not
			oldContact := bucket.list.Back().Value.(Contact)
			go bucket.heartbeat(oldContact, contact)
			return true, &oldContact // indicating that we want to ping
		}
	} else {
		bucket.list.MoveToFront(element)
	}
	return false, nil
}

func (bucket *bucket) heartbeat(oldContact Contact, newContact Contact) {

	time.Sleep(heartbeatTimeout)
	if bucket.list.Back().Value != oldContact {
		// oldest node returned ping in time
		return
	}
	// oldest node did not return the ping
	bucket.list.Remove(bucket.list.Back())
	bucket.list.PushFront(newContact)

}

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

// Len return the size of the bucket
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
