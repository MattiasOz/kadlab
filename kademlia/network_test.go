package kademlia

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	localId := NewRandomKademliaID()
	localIp := GetLocalIP()
	localContact := NewContact(localId, localIp)
	localRouting := NewRoutingTable(localContact)
	net := NetworkInit(&localContact, localRouting)
	data := CommData{localIp, localIp, *localId, ":3000", ":3000", "com", "computer", false, *localId}
	net.sendCh <- data
	if data != <-net.receiveCh {
		t.Error("The data did not match")
	}
}

func TestConvertDataToContactlist(t *testing.T) {
	localContact := Contact{NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "192.128.0.1", NewKademliaID("0000000000000000000000000000000000000000")}
	testID1 := NewKademliaID("FFFFFFFF0000000000000000000000000000000F")
	testID2 := NewKademliaID("FFFFFFFF000000000000000000000000000000F0")
	testID3 := NewKademliaID("FFFFFFFF00000000000000000000000000000F00")

	testIP1 := "192.128.0.2:3000"
	testIP2 := "192.128.0.3:3000"
	testIP3 := "192.128.0.4:3000"

	testString := fmt.Sprintf("%s,%s;%s,%s;%s,%s;", testID1, testIP1, testID2, testIP2, testID3, testIP3)
	fmt.Println("The test string is: ", testString)
	testContacts := ConvertDataToContactlist(testString, localContact, *localContact.ID)

	contact1 := Contact{testID1, testIP1, testID1.CalcDistance(*localContact.ID)}
	contact2 := Contact{testID2, testIP2, testID2.CalcDistance(*localContact.ID)}
	contact3 := Contact{testID3, testIP3, testID3.CalcDistance(*localContact.ID)}
	correctContacts := [...]Contact{contact1, contact2, contact3}

	for i := 0; i < 3; i++ {
		if *testContacts[i].ID != *correctContacts[i].ID {
			t.Errorf("TestConvertDataToContactlist failed, element %v should've had ID %v but had ID %v", i, correctContacts[i].ID, testContacts[i].ID)
		}
		if testContacts[i].Address != correctContacts[i].Address {
			t.Errorf("TestConvertDataToContactlist failed, element %v should've had address %v but had address %v", i, correctContacts[i].Address, testContacts[i].Address)
		}
		if *testContacts[i].distance != *correctContacts[i].distance {
			t.Errorf("TestConvertDataToContactlist failed, element %v should've had distance %v but had distance %v", i, correctContacts[i].distance, testContacts[i].distance)
		}
	}
}

func TestSendPingMessage(t *testing.T) {
	id := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	localContact := NewContact(id, "localhost:8000")
	targetContact := NewContact(
		NewRandomKademliaID(),
		"172.18.0.3",
	)

	receive := make(chan CommData, 10)
	send := make(chan CommData, 10)
	contacts := make(map[KademliaID]chan Contact)
	dataChs := make(map[KademliaID]chan string)
	network := Network{receive, send, &localContact, contacts, map[KademliaID]string{}, dataChs}
	network.SendPingMessage(&targetContact, false)
	message := <-send
	testCommData := CommData{network.localContact.Address, targetContact.Address, *(network.localContact.ID), ":3000", ":3000", PING, "", false, *NewKademliaID("0000000000000000000000000000000000000000")}
	if message != testCommData {
		t.Errorf("TestSendPingMessage failed, got %v, expected %v", message, testCommData)
	}
}

func TestSendFindContactMessage(t *testing.T) {
	id := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	localContact := NewContact(id, "localhost:8000")
	targetContact := NewContact(
		NewRandomKademliaID(),
		"172.18.0.3",
	)

	receive := make(chan CommData, 10)
	send := make(chan CommData, 10)
	contacts := make(map[KademliaID]chan Contact)
	dataStore := make(map[KademliaID]string)
	dataChs := make(map[KademliaID]chan string)
	network := Network{receive, send, &localContact, contacts, dataStore, dataChs}
	network.SendFindContactMessage(&targetContact, *network.localContact.ID)
	message := <-send
	testCommData := CommData{network.localContact.Address, targetContact.Address, *(network.localContact.ID), ":3000", ":3000", FIND_CONTACT, "", false, *(network.localContact.ID)}
	if message != testCommData {
		t.Errorf("TestSendFindContactMessage failed, got %v, expected %v", message, testCommData)
	}
}

func TestSendStoreMessage(t *testing.T) {
	id := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	localContact := NewContact(id, "localhost:8000")
	targetContact := NewContact(
		NewRandomKademliaID(),
		"172.18.0.3",
	)

	receive := make(chan CommData, 10)
	send := make(chan CommData, 10)
	contacts := make(map[KademliaID]chan Contact)
	dataStore := make(map[KademliaID]string)
	dataChs := make(map[KademliaID]chan string)
	network := Network{receive, send, &localContact, contacts, dataStore, dataChs}
	data := []byte("thisstringisatest")
	network.SendStoreMessage(data, targetContact, *network.localContact.ID)
	message := <-send
	testCommData := CommData{network.localContact.Address, targetContact.Address, *(network.localContact.ID), ":3000", ":3000", STORE, string(data), false, *(network.localContact.ID)}
	if message != testCommData {
		t.Errorf("TestSendStoreMessage failed, got %v, expected %v", message, testCommData)
	}
}

func TestSendFindContactResponse(t *testing.T) {
	id := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	localContact := NewContact(id, "localhost:8000")
	targetContact1 := NewContact(
		NewKademliaID("FFFFFFFF0000000000000000000000000000000F"),
		"172.18.0.3",
	)
	targetContact2 := NewContact(
		NewKademliaID("FFFFFFFF000000000000000000000000000000FF"),
		"172.18.0.4",
	)
	targetContact3 := NewContact(
		NewKademliaID("FFFFFFFF000000000000000000000000000000F0"),
		"172.18.0.5",
	)

	routingTable := NewRoutingTable(localContact)
	routingTable.AddContact(targetContact1)
	routingTable.AddContact(targetContact2)
	routingTable.AddContact(targetContact3)

	receive := make(chan CommData, 10)
	send := make(chan CommData, 10)
	contacts := make(map[KademliaID]chan Contact)
	dataStore := make(map[KademliaID]string)
	dataChs := make(map[KademliaID]chan string)
	network := Network{receive, send, &localContact, contacts, dataStore, dataChs}
	network.SendFindContactResponse(&targetContact1, routingTable, targetContact1.ID.String(), *(network.localContact.ID))
	message := <-send
	orderedContacts := ""
	orderedContacts = orderedContacts + fmt.Sprintf("%s,%s;", targetContact1.ID, targetContact1.Address)
	orderedContacts = orderedContacts + fmt.Sprintf("%s,%s;", targetContact3.ID, targetContact3.Address)
	orderedContacts = orderedContacts + fmt.Sprintf("%s,%s;", targetContact2.ID, targetContact2.Address)

	testCommData := CommData{network.localContact.Address, targetContact1.Address, *(network.localContact.ID), ":3000", ":3000", FIND_CONTACT, orderedContacts, true, *(network.localContact.ID)}
	if message != testCommData {
		t.Errorf("TestSendFindContactResponse failed, got %v, expected %v", message, testCommData)
	}
}

// // This test sucks
// func TestListen(t *testing.T) {
//     Init()
//     t.Log("Listen didn't error")
// }
//
// func TestBroadcast(t *testing.T) {
//     // check if we can ball from broadcast to listen. Otherwise this won't amount to much
//     _, sendCh := Init()
//     reader := bufio.NewReader(os.Stdin)
//     sendCh <- CommData{"172.18.0.2", "172.18.0.3", "1", ":3000", ":3000", "com", "computer"}
//     time.Sleep(1*time.Second)
//     text, _ := reader.ReadString('\n')
//     fmt.Println("the text was", text)
// }
