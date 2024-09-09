package kademlia

import(
    "net"
    "fmt"
    "encoding/json"
)

const port = ":3000"

type Network struct {
}

type CommData struct {
    SenderIP string
    ReceiverIP string
    SenderID string
    SenderPort string
    Port string
    RPCCommand string
    Data string
}

func Init(something string) {
    // TODO
}

func Listen(ip string, port string) {
    localAddress, err := net.ResolveUDPAddr("udp", port)
    if err != nil {
        fmt.Println("ERROR: ", err)
    }
    connection, err := net.ListenUDP("udp", localAddress)
    if err != nil {
        fmt.Println("ERROR: ", err)
    }
    defer connection.Close()
    for {
        var message CommData
        buffer := make([]byte, 4096)
        length, _, err := connection.ReadFromUDP(buffer)
        if err != nil {
            fmt.Println("ERROR: ", err)
        }
        buffer = buffer[:length]
        err = json.Unmarshal(buffer, &message)
        if err != nil {
            fmt.Println("ERROR: ", err)
        }
        // if (message.Identifier == com_id) {
        //     fmt.Println("Msg from: ", message.SenderIP)
        // }
        fmt.Println("Msg: ", message)
    }
}

func (network *Network) SendPingMessage(contact *Contact) {
   // TODO 
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
