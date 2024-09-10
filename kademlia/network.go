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

func Init() (<-chan CommData, chan<- CommData) {
    receive := make(chan CommData, 10)
    send := make(chan CommData, 10)
    go Listen(receive)
    go Broadcast(send)
    return receive, send
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
