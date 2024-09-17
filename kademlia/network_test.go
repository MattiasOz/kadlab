package kademlia

import (
	"testing"
)

func TestInit(t *testing.T) {
	localId := NewRandomKademliaID()
	localIp := GetLocalIP()
	localContact := NewContact(localId, localIp)
	localRouting := NewRoutingTable(localContact)
	net := NetworkInit(&localContact, localRouting)
	data := CommData{localIp, localIp, *localId, ":3000", ":3000", "com", "computer", false}
	net.sendCh <- data
	if data != <-net.receiveCh {
		t.Error("The data did not match")
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
