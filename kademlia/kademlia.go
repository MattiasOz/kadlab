package kademlia

import (
	"sort"
	"time"
)

const concurrencyParameter = 3

type Kademlia struct {
	routingTable *RoutingTable
	network      *Network
}

func Init() Kademlia {
	localID := NewRandomKademliaID()
	localContact := NewContact(localID, GetLocalIP())
	routingTable := NewRoutingTable(localContact)

	network := NetworkInit(&localContact, routingTable)

	res := Kademlia{
		routingTable: routingTable,
		network:      network,
	}
	return res
}

// TODO: this should probably be automated and removed
func (kademlia *Kademlia) Ping(target *Contact) {
	kademlia.network.SendPingMessage(target, false)
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) {
	nClosestContacts := kademlia.routingTable.FindClosestContacts(target, concurrencyParameter)

	for _, contact := range nClosestContacts {
		kademlia.network.SendFindContactMessage(&contact, target.String())
	}

	go kademlia.ProcessContactLookupReturns(target)
}

func (Kademlia *Kademlia) ProcessContactLookupReturns(target *KademliaID) {
	time.Sleep(3 * time.Second)
	var closestNodeSeen Contact
	var contactList []Contact
	for {
		if len(Kademlia.network.findContactCh) == 0 {
			break
		}
		returnedContact := <-Kademlia.network.findContactCh
		if !IsContactAlreadyInList(contactList, returnedContact) {
			contactList = append(contactList, returnedContact)
		}

	}
	sort.Slice(contactList, func(i, j int) bool { return contactList[i].distance.Less(contactList[j].distance) })
	closestNodeSeen = contactList[0]

	var newClosestNode Contact
	leftToSendTo := concurrencyParameter
	for {
		if newClosestNode.distance == closestNodeSeen.distance {
			break
		}
	}
}

func IsContactAlreadyInList(contactList []Contact, newContact Contact) bool {
	for _, contact := range contactList {
		if contact.ID == newContact.ID {
			return true
		}
	}
	return false
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
