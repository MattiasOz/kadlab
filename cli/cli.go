package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
    _"flag"
    "d7024e/kademlia"
)

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
    conn, err := net.Dial("unix", kademlia.SOCKET_PATH)
    if err != nil {
        fmt.Println("Couldn't dial socker", err)
    }
    defer conn.Close()

    msg, err := json.Marshal(message)
    if err != nil {
        fmt.Println("Error marshalling:", err)
        return
    }
    conn.Write(msg)

    buf := make([]byte, 1024)
    length, err := conn.Read(buf)
    if err != nil {
        fmt.Println("Error reading:", err)
    }
    buf = buf[:length]

    var resp string
    err = json.Unmarshal(buf, &resp)
    if err != nil {
        fmt.Println("Error:", err)
    }

    fmt.Println("Response was", resp)
    fmt.Println("Response was", string(buf[:length]))
}
