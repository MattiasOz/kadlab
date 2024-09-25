package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

const port = ":3000"

type Network struct {
	receiveCh    chan CommData
	sendCh       chan CommData
	localContact *Contact
	lookupChs    map[KademliaID]chan Contact

	dataStore map[KademliaID]string
}

type CommData struct {
	SenderIP   string
	ReceiverIP string
	SenderID   KademliaID
	SenderPort string
	Port       string
	RPCCommand string
	Data       string
	Response   bool
	QueryID    KademliaID
}

// func (network *Network) Init() (<-chan CommData, chan<- CommData) {
func NetworkInit(localContact *Contact, routingTable *RoutingTable) *Network {
	receive := make(chan CommData, 10)
	send := make(chan CommData, 10)
	go Listen(receive)
	go Broadcast(send)
	contacts := make(map[KademliaID]chan Contact)
	dataStore := make(map[KademliaID]string)
	network := Network{receive, send, localContact, contacts, dataStore}
	go network.Interpreter(receive, routingTable)
	return &network
	// return receive, send
}

func Listen(commReceive chan CommData) {
	localAddress, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		fmt.Println("cannot resolve UDPAddr: ", err)
	}
	connection, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		fmt.Println("cannot listenUDP: ", err)
	}
	defer connection.Close()
	for {
		var message CommData
		buffer := make([]byte, 4096)
		length, _, err := connection.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("cannot ReadFromUDP: ", err)
		}
		buffer = buffer[:length]
		err = json.Unmarshal(buffer, &message)
		if err != nil {
			fmt.Println("cannot unmarshal: ", err)
		}
		// if (message.Identifier == com_id) {
		//     fmt.Println("Msg from: ", message.SenderIP)
		// }
		fmt.Println("Msg: ", message)
		commReceive <- message
	}
}

func Broadcast(commSend chan CommData) {
	for {
		message := <-commSend
		ip := message.ReceiverIP

		fmt.Println("COMM: Broadcasting message to: " + ip + port)
		broadcastAddress, err := net.ResolveUDPAddr("udp", ip+port)
		if err != nil {
			fmt.Println("ERROR: ", err)
		}

		localAddress, err := net.ResolveUDPAddr("udp", GetLocalIP())
		connection, err := net.DialUDP("udp", localAddress, broadcastAddress)
		if err != nil {
			fmt.Println("ERROR: ", err)
		}

		convMsg, err := json.Marshal(message)
		if err != nil {
			fmt.Println("ERROR: ", err)
		}
		connection.Write(convMsg)

		connection.Close()
	}
}

func (network *Network) Interpreter(commReceive chan CommData, routingTable *RoutingTable) {
	for {
		receivedCommData := <-commReceive
		senderID := receivedCommData.SenderID
		senderIP := receivedCommData.SenderIP
		senderContact := NewContact(&senderID, senderIP)
		switch receivedCommData.RPCCommand {
		case PING:
			if !receivedCommData.Response {
				network.SendPingMessage(&senderContact, true)
			}
		case FIND_CONTACT:
			if !receivedCommData.Response {
				network.SendFindContactResponse(&senderContact, routingTable, receivedCommData.Data, receivedCommData.QueryID)
			} else {
				kClosestContactsToContactedNode := ConvertDataToContactlist(receivedCommData.Data, *network.localContact, receivedCommData.QueryID)
				for _, contact := range kClosestContactsToContactedNode {
					if _, ok := network.lookupChs[receivedCommData.QueryID]; ok { // Check that the channel exists
						network.lookupChs[receivedCommData.QueryID] <- contact
					}
					shouldPing, pingTo := routingTable.AddContact(contact)
					if shouldPing { // Bucket is full, heartbeat ping oldest contact in the bucket
						network.SendPingMessage(pingTo, false)
					}
				}
			}
		case FIND_DATA:
			//TODO
		case STORE:
			network.StoreData(receivedCommData.Data, receivedCommData.QueryID)
		default:
			fmt.Println("Error. In the default case of Interpreter")
			continue
		}
		shouldPing, pingTo := routingTable.AddContact(senderContact) //
		if shouldPing {                                              // Bucket is full, heartbeat ping oldest contact in the bucket
			network.SendPingMessage(pingTo, false)
		}
	}
}

func GetLocalIP() string {
	var localIP string
	addr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Printf("GetLocalIP in communication failed")
		return "localhost"
	}

	for _, val := range addr {

		if ip, ok := val.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				localIP = ip.IP.String()
			}
		}
	}

	return localIP
}

func (network *Network) SendPingMessage(contact *Contact, response bool) {
	pingData := CommData{
		network.localContact.Address,
		contact.Address,
		*(network.localContact.ID),
		port,
		port,
		PING,
		"",
		response,
		*NewKademliaID("0000000000000000000000000000000000000000"),
	}

	network.sendCh <- pingData
}

func (network *Network) SendFindContactMessage(contact *Contact, findContactID KademliaID) {
	findContactMessage := CommData{
		network.localContact.Address,
		contact.Address,
		*(network.localContact.ID),
		port,
		port,
		FIND_CONTACT,
		"",
		false,
		findContactID,
	}
	network.sendCh <- findContactMessage
}

func (network *Network) SendFindContactResponse(contact *Contact, routingTable *RoutingTable, receivedCommData string, queryID KademliaID) {
	kademliaIDtoLookup := NewKademliaID(receivedCommData)
	kClosestContacts := routingTable.FindClosestContacts(kademliaIDtoLookup, bucketSize)

	kClosestContactData := ""
	for _, contact := range kClosestContacts {
		kClosestContactData = kClosestContactData + fmt.Sprintf("%s,%s;", contact.ID, contact.Address)
	}

	findContactResponse := CommData{
		network.localContact.Address,
		contact.Address,
		*(network.localContact.ID),
		port,
		port,
		FIND_CONTACT,
		kClosestContactData,
		true,
		queryID,
	}
	network.sendCh <- findContactResponse
}

func (network *Network) CreateNewFindContactChannel(queryID KademliaID) {
	lookupChannel := make(chan Contact, 10)
	network.lookupChs[queryID] = lookupChannel
}

func (network *Network) RemoveFindContactChannel(queryID KademliaID) {
	delete(network.lookupChs, queryID)
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte, contact Contact, queryID KademliaID) {
	storeMessage := CommData{
		network.localContact.Address,
		contact.Address,
		*(network.localContact.ID),
		port,
		port,
		STORE,
		string(data),
		false,
		queryID,
	}

	network.sendCh <- storeMessage
}

func (network *Network) StoreData(data string, queryID KademliaID) {
	network.dataStore[queryID] = data
}

func ConvertDataToContactlist(recievedCommData string, localContact Contact, queryID KademliaID) (contacts []Contact) {
	contactStrings := strings.Split(recievedCommData, ";")
	contactStrings = contactStrings[:len(contactStrings)-1]

	for _, contactString := range contactStrings {
		contactElements := strings.Split(contactString, ",")
		kademliaID := NewKademliaID(contactElements[0])
		fmt.Println("The kademliaId we found was ", kademliaID)
		address := contactElements[1]
		distance := kademliaID.CalcDistance(queryID)

		contact := Contact{kademliaID, address, distance}

		contacts = append(contacts, contact)
	}

	return contacts
}
