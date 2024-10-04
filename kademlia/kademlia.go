package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
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

func (kademlia *Kademlia) LookupContact(target *KademliaID) ([]Contact, error) {
	kademlia.network.CreateNewLookupChannel(*target)

	nClosestContacts := kademlia.routingTable.FindClosestContacts(target, concurrencyParameter)

	for _, contact := range nClosestContacts { //Kontakta alfa närmsta kontakterna till målet
		kademlia.network.SendFindContactMessage(&contact, *target)
	}

	return kademlia.ProcessContactLookupReturns(target)
}

func (Kademlia *Kademlia) ProcessContactLookupReturns(target *KademliaID) ([]Contact, error) {

	time.Sleep(3 * time.Second)
	var closestNodeSeen Contact
	var contactList []Contact
	for { //Ta emot svaren (från de som hunnit svara)
		if len(Kademlia.network.lookupChs[*target]) == 0 {
			break
		}
		returnedContact := <-Kademlia.network.lookupChs[*target]
		if !IsContactAlreadyInList(contactList, returnedContact) {
			contactList = append(contactList, returnedContact)
		}

	}
	sort.Slice(contactList, func(i, j int) bool {
		return contactList[i].distance.Less(contactList[j].distance)
	})
	if len(contactList) >= 1 {
		closestNodeSeen = contactList[0]
	} else {
		return nil, errors.New("Couldn't find any contacts")
	}

	var newClosestNode Contact
	for {
		if newClosestNode.distance == closestNodeSeen.distance { // Om vi inte hittar någon närmare nod går vi vidare
			break
		}
		tmp := int(math.Min(concurrencyParameter, float64(len(contactList))))
		closestNodeSeen = newClosestNode
		for i := 0; i < tmp; i++ { // Kontakta alfa nya närmsta
			Kademlia.network.SendFindContactMessage(&contactList[i], *target)
		}
		time.Sleep(2 * time.Second)
		for { // Ta emot svaren från alfa nya närmsta
			if len(Kademlia.network.lookupChs[*target]) == 0 {
				break
			}
			returnedContact := <-Kademlia.network.lookupChs[*target]
			if !IsContactAlreadyInList(contactList, returnedContact) {
				contactList = append(contactList, returnedContact)
			}
		}
		sort.Slice(contactList, func(i, j int) bool {
			return contactList[i].distance.Less(contactList[j].distance)
		})
		newClosestNode = contactList[0]
	}

	// Skicka och ta emot svaren från k närmsta till målet som är okontaktade
	if len(contactList) >= bucketSize {
		for i := concurrencyParameter; i < bucketSize; i++ {
			Kademlia.network.SendFindContactMessage(&contactList[i], *target)
		}
	} else {
		for i := concurrencyParameter; i < len(contactList); i++ {
			Kademlia.network.SendFindContactMessage(&contactList[i], *target)
		}
	}

	time.Sleep(2 * time.Second)
	for {
		if len(Kademlia.network.lookupChs[*target]) == 0 {
			break
		}
		returnedContact := <-Kademlia.network.lookupChs[*target]
		if !IsContactAlreadyInList(contactList, returnedContact) {
			contactList = append(contactList, returnedContact)
		}
	}
	sort.Slice(contactList, func(i, j int) bool {
		return contactList[i].distance.Less(contactList[j].distance)
	})

	Kademlia.network.RemoveLookupChannel(*target)

	if len(contactList) > bucketSize {
		return contactList[:bucketSize], nil
	} else {
		return contactList, nil
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

func (Kademlia *Kademlia) LookupData(hash string) (string, error) {
	target := NewKademliaID(hash)
	Kademlia.network.CreateNewLookupChannel(*target)
	Kademlia.network.CreateNewDataChannel(*target)

	nClosestContacts := Kademlia.routingTable.FindClosestContacts(target, concurrencyParameter)

	for _, contact := range nClosestContacts { //Kontakta alfa närmsta kontakterna till målet
		Kademlia.network.SendFindDataMessage(contact, *target)
	}

	returnedData, err := Kademlia.ProcessDataLookupReturns(target)

	Kademlia.network.RemoveLookupChannel(*target)
	Kademlia.network.RemoveDataChannel(*target)

	return returnedData, err

}

func (Kademlia *Kademlia) ProcessDataLookupReturns(target *KademliaID) (string, error) {
	time.Sleep(3 * time.Second)
	var closestNodeSeen Contact
	var contactList []Contact
	for { //Ta emot svaren (från de som hunnit svara)
		if len(Kademlia.network.dataChs[*target]) > 0 { // När datan hittas returnerar vi direkt och sparar i närmsta nod som inte hade datan
			foundData := <-Kademlia.network.dataChs[*target]
			if len(contactList) > 0 {
				sort.Slice(contactList, func(i, j int) bool {
					return contactList[i].distance.Less(contactList[j].distance)
				})
				Kademlia.network.SendStoreMessage([]byte(foundData), contactList[0], *target)
			}
			return foundData, nil
		}
		if len(Kademlia.network.lookupChs[*target]) == 0 {
			break
		}
		returnedContact := <-Kademlia.network.lookupChs[*target]
		if !IsContactAlreadyInList(contactList, returnedContact) {
			contactList = append(contactList, returnedContact)
		}

	}
	sort.Slice(contactList, func(i, j int) bool {
		return contactList[i].distance.Less(contactList[j].distance)
	})
	if len(contactList) >= 1 {
		closestNodeSeen = contactList[0]
	} else {
		return "", errors.New("PrDataLookup: Couldn't find any contacts")
	}

	var newClosestNode Contact
	for {
		if newClosestNode.distance == closestNodeSeen.distance { // Om vi inte hittar någon närmare nod går vi vidare
			break
		}
		closestNodeSeen = newClosestNode
		for i := 0; i < concurrencyParameter; i++ { // Kontakta alfa nya närmsta
			Kademlia.network.SendFindDataMessage(contactList[i], *target)
		}
		time.Sleep(2 * time.Second)
		for { // Ta emot svaren från alfa nya närmsta
			if len(Kademlia.network.dataChs[*target]) > 0 { // När datan hittas returnerar vi direkt
				foundData := <-Kademlia.network.dataChs[*target]
				sort.Slice(contactList, func(i, j int) bool {
					return contactList[i].distance.Less(contactList[j].distance)
				})
				Kademlia.network.SendStoreMessage([]byte(foundData), contactList[0], *target)
				return foundData, nil
			}
			if len(Kademlia.network.lookupChs[*target]) == 0 {
				break
			}
			returnedContact := <-Kademlia.network.lookupChs[*target]
			if !IsContactAlreadyInList(contactList, returnedContact) {
				contactList = append(contactList, returnedContact)
			}
		}
		sort.Slice(contactList, func(i, j int) bool {
			return contactList[i].distance.Less(contactList[j].distance)
		})
		newClosestNode = contactList[0]
	}

	// Skicka och ta emot svaren från k närmsta till målet som är okontaktade, datan kan ju råka finnas där
	if len(contactList) >= bucketSize {
		for i := concurrencyParameter; i < bucketSize; i++ {
			Kademlia.network.SendFindDataMessage(contactList[i], *target)
		}
	} else {
		for i := concurrencyParameter; i < len(contactList); i++ {
			Kademlia.network.SendFindDataMessage(contactList[i], *target)
		}
	}

	time.Sleep(2 * time.Second)
	if len(Kademlia.network.dataChs[*target]) > 0 { // Om datan nu har hittats returnerar vi den, annars finns den nog inte
		foundData := <-Kademlia.network.dataChs[*target]
		sort.Slice(contactList, func(i, j int) bool {
			return contactList[i].distance.Less(contactList[j].distance)
		})
		Kademlia.network.SendStoreMessage([]byte(foundData), contactList[0], *target)
		return foundData, nil
	}

	return "", nil
}

func (kademlia *Kademlia) Store(data []byte) (string, error) {
	storedDataHash := sha1.Sum(data)

	kademliaIDofStoredDataHash := NewKademliaID(hex.EncodeToString(storedDataHash[:]))
	contactsToStoreDataIn, err := kademlia.LookupContact(kademliaIDofStoredDataHash)

	if err != nil {
		errf := fmt.Sprintf("Error: Failed to store, reason: %v", err)
		return "", errors.New(errf)
	}

	for _, contact := range contactsToStoreDataIn {
		kademlia.network.SendStoreMessage(data, contact, *kademliaIDofStoredDataHash)
	}

	return kademliaIDofStoredDataHash.String(), nil
}
