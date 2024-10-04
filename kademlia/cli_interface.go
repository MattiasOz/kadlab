package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Cli_command struct {
	RPC_command string
	Content     string
}

type Cli_response struct {
	Content string
}

func Cli_Start(kadlab Kademlia, exit_ch chan bool) {
	go cli_listener(kadlab, exit_ch)
}

func cli_listener(kadlab Kademlia, exit_ch chan bool) {
	if err := os.RemoveAll(SOCKET_PATH); err != nil {
		fmt.Println("Error, clearing socket")
		panic(err)
	}

	ln, err := net.Listen("unix", SOCKET_PATH)
	if err != nil {
		fmt.Println("Error listening to socker")
		panic(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting package", err)
			continue
		}
		defer conn.Close()

		buf := make([]byte, 1024) //arbitrary number
		length, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading from package", err)
			continue
		}
		buf = buf[:length]
		var command Cli_command
		err = json.Unmarshal(buf, &command)
		if err != nil {
			fmt.Println("Error unmarshaling", err)
			continue
		}
		go handCliInput(command, conn, kadlab, exit_ch)
	}
}

func handCliInput(command Cli_command, conn net.Conn, kadlab Kademlia, exit_ch chan bool) {
	defer conn.Close()
	var res []byte
	switch command.RPC_command {
	case PUT:
		hash, err := kadlab.Store([]byte(command.Content))
		if err != nil {
			fmt.Println("Error in PUT:", err)
			return
		}
		res = []byte(hash)
	case GET:
		data, err := kadlab.LookupData(command.Content)
		if err != nil {
			fmt.Println("Error in GET:", err)
			return
		}
		res = []byte(data)
	case EXIT:
		exit_ch <- true
		res = []byte("Exiting")
	default:
		fmt.Println("Invalid CLI command")
		res = []byte("Invalid command")
	}
	_, err := conn.Write(res)
	if err != nil {
		fmt.Println("Failed to send response to CLI:", err)
	}
}
