// TODO: Add package documentation for `main`, like this:
// Package main something something...
package main

import (
	"d7024e/kademlia"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("Running the kademlia app...")

	kadlab := kademlia.Init()

	if kademlia.GetLocalIP() != "172.18.0.3" {
		time.Sleep((time.Duration(5 + rand.Intn(10))) * time.Second)
		bootstrapNode := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "172.18.0.3")
		kadlab.Ping(&bootstrapNode)
		time.Sleep((time.Duration(5 + rand.Intn(10))) * time.Second)
		kadlab.LookupSelf()
	} else {
		bootstrapNode := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "172.18.0.3")
		kadlab.Ping(&bootstrapNode)
	}
	fmt.Println("\033[32mStartup is complete\033[0m")

	exit_ch := make(chan bool)
	kademlia.Cli_Start(kadlab, exit_ch)

	for {
		select {
		case <-exit_ch:
			fmt.Println("Shutting down")
			return
		case <-time.After(1 * time.Second):
			continue
		}
	}
}
