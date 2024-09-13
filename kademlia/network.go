package kademlia

import (
	"encoding/json"
	"fmt"
	"net"
)

const port = ":3000"

type Network struct {
    receiveCh chan CommData
    sendCh chan CommData
    localContact *Contact
}

type CommData struct {
    SenderIP string
    ReceiverIP string
    SenderID KademliaID
    SenderPort string
    Port string
    RPCCommand string
    Data string
    Response bool
}

// func (network *Network) Init() (<-chan CommData, chan<- CommData) {
func Init(localContact *Contact) (*Network) {
    receive := make(chan CommData, 10)
    send := make(chan CommData, 10)
    go Listen(receive)
    go Broadcast(send)
    network := Network{receive, send, localContact}
    go network.Interpreter(receive)
    return &network
    // return receive, send
}

func Listen(commReceive chan CommData) {
    localAddress, err := net.ResolveUDPAddr("udp", port)
    if err != nil {
        fmt.Println("cannot resolve UDPAddr: ", err)
    }
    connection, err := net.ListenUDP("udp", localAddress)
    if err != nil {
        fmt.Println("cannot listenUDP: ", err)
    }
    defer connection.Close()
    for {
        var message CommData
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
        // if (message.Identifier == com_id) {
        //     fmt.Println("Msg from: ", message.SenderIP)
        // }
        fmt.Println("Msg: ", message)
        commReceive <- message
    }
}

func Broadcast(commSend chan CommData) {
    for{
        message := <- commSend
        ip := message.ReceiverIP

        fmt.Println("COMM: Broadcasting message to: " + ip + port)
        broadcastAddress, err := net.ResolveUDPAddr("udp", ip + port)
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
        fmt.Println("COMM: Message sent successfully!")

        connection.Close()
    }
}

func (network *Network) Interpreter(commReceive chan CommData) {
    for {
        receivedCommData := <- commReceive
        senderID := receivedCommData.SenderID
        senderIP := receivedCommData.SenderIP
        senderContact := NewContact(&senderID, senderIP)
        switch receivedCommData.RPCCommand {
        case PING:
            if !receivedCommData.Response {
                network.SendPingMessage(&senderContact, true)
            }
        default:
            fmt.Println("Error. In the default case of Interpreter")
            continue
        }
    }
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

func (network *Network) SendPingMessage(contact *Contact, response bool) {
    pingData := CommData{
        network.localContact.Address,
        contact.Address,
        *(network.localContact.ID),
        port,
        port,
        PING,
        "",
        response,
    }
    
    network.sendCh <- pingData
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
