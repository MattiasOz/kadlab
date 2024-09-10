package kademlia

import (
	"testing"
)

func TestInit(t *testing.T) {
    Init()
    t.Log("Init succeeded")
}

// This test sucks
func TestListen(t *testing.T) {
    receive, _ := Init()
    Listen(receive)
    t.Log("Listen didn't error")
}

func TestBroadcast(t *testing.T) {
    // check if we can ball from broadcast to listen. Otherwise this won't amount to much
}
