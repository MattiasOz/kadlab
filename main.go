// TODO: Add package documentation for `main`, like this:
// Package main something something...
package main

import (
	"d7024e/kademlia"
	"fmt"
    "time"
)

func main() {
	fmt.Println("Pretending to run the kademlia app...")
	// Using stuff from the kademlia package here. Something like...
	id := kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	contact := kademlia.NewContact(id, "localhost:8000")
	fmt.Println(contact.String())
	fmt.Printf("%v\n", contact)

    // init
    // _, sendCh := kademlia.Init()
    localContact := kademlia.NewContact(
        kademlia.NewRandomKademliaID(),
        kademlia.GetLocalIP(),
    )
    targetContact := kademlia.NewContact(
        kademlia.NewRandomKademliaID(),
        "172.18.0.3",
    )
    network := kademlia.Init(&localContact)
    for {
        time.Sleep(10*time.Second)
        network.SendPingMessage(&targetContact, false)
    }
}
