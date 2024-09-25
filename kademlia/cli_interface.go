package kademlia

import (
    "net"
    "fmt"
    "encoding/json"
)

const local_port = ":80"

type Cli_command struct {
    RPC_command string
    Content string
}

func Cli_init() {
    receive := make(chan Cli_command, 10)
    go cli_listener(receive)
}

func cli_listener(resCh chan Cli_command){
	localAddress, err := net.ResolveUDPAddr("udp", local_port)
	if err != nil {
		fmt.Println("cannot resolve UDPAddr: ", err)
	}
	connection, err := net.ListenUDP("udp", localAddress)
	if err != nil {
		fmt.Println("cannot listenUDP: ", err)
	}
	defer connection.Close()
	for {
		var message Cli_command
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
		fmt.Println("THIS MESSGE IS Msg: ", message.RPC_command, message.Content)
	}
}
