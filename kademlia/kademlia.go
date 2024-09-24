package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
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
	if GetLocalIP() == "172.18.0.3" { // En hårdkodad bootstrap-nod
		localID = NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	}
	localContact := NewContact(localID, GetLocalIP())
	routingTable := NewRoutingTable(localContact)

	network := NetworkInit(&localContact, routingTable)

	res := Kademlia{
		routingTable: routingTable,
		network:      network,
	}
	return res
}

func (kademlia *Kademlia) Ping(target *Contact) {
	kademlia.network.SendPingMessage(target, false)
}

func (kademlia *Kademlia) LookupContact(target *KademliaID) []Contact {
	nClosestContacts := kademlia.routingTable.FindClosestContacts(target, concurrencyParameter)

	for _, contact := range nClosestContacts { //Kontakta alfa närmsta kontakterna till målet
		kademlia.network.SendFindContactMessage(&contact, target.String())
	}

	return kademlia.ProcessContactLookupReturns(target)
}

func (Kademlia *Kademlia) ProcessContactLookupReturns(target *KademliaID) []Contact {
	time.Sleep(3 * time.Second)
	var closestNodeSeen Contact
	var contactList []Contact
	for { //Ta emot svaren (från de som hunnit svara)
		if len(Kademlia.network.findContactCh) == 0 {
			break
		}
		returnedContact := <-Kademlia.network.findContactCh
		if !IsContactAlreadyInList(contactList, returnedContact) {
			contactList = append(contactList, returnedContact)
		}

	}
	sort.Slice(contactList, func(i, j int) bool {
		return contactList[i].distance.CalcDistance(*target).Less(contactList[j].distance.CalcDistance(*target))
	})
	closestNodeSeen = contactList[0]

	var newClosestNode Contact
	for {
		if newClosestNode.distance == closestNodeSeen.distance { // Om vi inte hittar någon närmare nod går vi vidare
			break
		}
		closestNodeSeen = newClosestNode
		for i := 0; i < concurrencyParameter; i++ { // Kontakta alfa nya närmsta
			Kademlia.network.SendFindContactMessage(&contactList[i], target.String())
		}
		time.Sleep(2 * time.Second)
		for { // Ta emot svaren från alfa nya närmsta
			if len(Kademlia.network.findContactCh) == 0 {
				break
			}
			returnedContact := <-Kademlia.network.findContactCh
			if !IsContactAlreadyInList(contactList, returnedContact) {
				contactList = append(contactList, returnedContact)
			}
		}
		sort.Slice(contactList, func(i, j int) bool {
			return contactList[i].distance.CalcDistance(*target).Less(contactList[j].distance.CalcDistance(*target))
		})
		newClosestNode = contactList[0]
	}

	// Skicka och ta emot svaren från k närmsta till målet som är okontaktade
	for i := concurrencyParameter; i < bucketSize; i++ {
		Kademlia.network.SendFindContactMessage(&contactList[i], target.String())
	}
	time.Sleep(2 * time.Second)
	for {
		if len(Kademlia.network.findContactCh) == 0 {
			break
		}
		returnedContact := <-Kademlia.network.findContactCh
		if !IsContactAlreadyInList(contactList, returnedContact) {
			contactList = append(contactList, returnedContact)
		}
	}
	sort.Slice(contactList, func(i, j int) bool {
		return contactList[i].distance.CalcDistance(*target).Less(contactList[j].distance.CalcDistance(*target))
	})

	if len(contactList) > bucketSize {
		return contactList[:bucketSize]
	} else {
		return contactList
	}
}

func (Kademlia *Kademlia) LookupSelf() { // Used for bootstrapping
	Kademlia.LookupContact(Kademlia.routingTable.me.ID)
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
	storedDataHash := sha1.Sum(data)

	kademliaIDofStoredDataHash := NewKademliaID(hex.EncodeToString(storedDataHash[:]))
	contactsToStoreDataIn := kademlia.LookupContact(kademliaIDofStoredDataHash)

	for _, contact := range contactsToStoreDataIn {
		kademlia.network.SendStoreMessage(data, contact)
	}
}
