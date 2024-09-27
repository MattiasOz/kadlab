package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Cli_command struct {
    RPC_command string
    Content string
}

func Cli_init() {
    receive := make(chan Cli_command, 10)
    go cli_listener(receive)
}

func cli_listener(resCh chan Cli_command){
    fmt.Println("STARTED cli_listener")
    if err := os.RemoveAll(SOCKET_PATH); err != nil {
        fmt.Println("Error, clearing socket", err)
    }   

    ln, err := net.Listen("unix", SOCKET_PATH)
    if err != nil {
        fmt.Println("Error listening to socker", err)
    }
    defer ln.Close()

    fmt.Println("Listening to", SOCKET_PATH)

    for {
        conn, err := ln.Accept()
        fmt.Println("Received a package")
        if err != nil {
            fmt.Println("Error accepting package", err)
            continue
        }
        defer conn.Close()

        buf := make([]byte, 1024)
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
        go handCliInput(command, conn)
    }
}


func handCliInput(command Cli_command, conn net.Conn) {
    defer conn.Close()
    fmt.Println("Command was", command)
    conn.Write([]byte("M1911a1"))
}
