package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
    "d7024e/kademlia"
)

const port = ":80"

type cli_command struct {
    rpc_command string
    content string
}

func main() {
    args := os.Args[1:]
    fmt.Println(args)
    handle_input(args)
}

func handle_input(args []string) {
    if !is_enough_args(args, 2) {
        return
    }

    command1 := strings.ToLower(args[0])
    command2 := strings.ToLower(args[1])
    
    switch command1 {
    case "ping":
        ping(args[1])
    case "find":
        if !is_enough_args(args, 3) {
            return
        }
        content := args[2:]
        switch command2 {
        case "node":
            find_node(content)
        case "data":
            find_data(content)
        }
    
    case "store":
        if !is_enough_args(args, 3) {
            return
        }
        content := args[1:]
        store(content)
    }
}


func is_enough_args(args []string, length int) bool {
    if len(args) < length {
        fmt.Println("cli <command> <target-ip>")
        fmt.Println("commands: [ping]/[find node]/[find data]/[store]")
        return false
    }
    return true
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

func ping(ip string) {
    message := kademlia.Cli_command {
        RPC_command: kademlia.PING,
        Content: "",
    }
    send_message(ip, message)
}

func find_node(args []string) {
    ip := args[0]
    message := kademlia.Cli_command {
        RPC_command: kademlia.FIND_CONTACT,
        Content: strings.Join(args[1:], " "),
    }
    send_message(ip, message)
}

func find_data(args []string) {
    ip := args[0]
    message := kademlia.Cli_command {
        RPC_command: kademlia.FIND_DATA,
        Content: strings.Join(args[1:], " "),
    }
    send_message(ip, message)
}

func store(args []string) {
    ip := args[0]
    message := kademlia.Cli_command {
        RPC_command: kademlia.STORE,
        Content: strings.Join(args[1:], " "),
    }
    send_message(ip, message)
}

func send_message(ip string, message kademlia.Cli_command) {
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
