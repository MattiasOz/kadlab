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
    targetContact := kademlia.NewContact(
        kademlia.NewRandomKademliaID(),
        "172.18.0.3",
    )
    kadlab := kademlia.Init()
    for {
        kadlab.Pring(&targetContact)
        time.Sleep(10*time.Second)
    }
}
