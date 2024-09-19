package kademlia

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

func (kademlia *Kademlia) LookupContact(target *Contact) {
	nClosestContacts := kademlia.routingTable.FindClosestContacts(target.ID, concurrencyParameter)

	for _, contact := range nClosestContacts {
		kademlia.network.SendFindContactMessage(&contact)
	}

	go kademlia.ProcessContactLookupReturns(target)
}

func (Kademlia *Kademlia) ProcessContactLookupReturns(target *Contact) {
	for {

	}
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
