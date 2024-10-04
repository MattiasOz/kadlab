package main

import (
	"d7024e/kademlia"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

const timeout = 60 * time.Second

func main() {
	command, data := handle_args()
	send_message(command, data)
}

func handle_args() (string, string) {
	command := strings.ToUpper(os.Args[1])
	os.Args = append(os.Args[0:1], os.Args[2:]...)
	data := flag.String("data", "", "Data to be sent")
	flag.Parse()
	return command, *data
}

func send_message(command string, data string) {
	message := kademlia.Cli_command{
		RPC_command: command,
		Content:     data,
	}

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
	_, err = conn.Write(msg)
	if err != nil {
		fmt.Println("Error sending:", err)
		return
	}
	fmt.Println("Message sent successfully")

	conn.SetReadDeadline(time.Now().Add(timeout))
	buf := make([]byte, 1024)
	length, err := conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			fmt.Println("No data was returned")
			return
		} else {
			panic(err)
		}
	}
	buf = buf[:length]

	fmt.Printf("Response was:\n%s\n", string(buf[:length]))
}
