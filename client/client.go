package client

import (
    "bufio"
    "bytes"
    "fmt"
    "net"
    "os"
    "strings"
)

var (
    CLOSE_CONNECTION = "CLOSECONNECTION"
)

type Client struct {
    socket net.Conn
}

func StartClient() {
    host := "localhost:12345"
    fmt.Println("Starting client...")
    connection, error := net.Dial("tcp", host)
    if error != nil {
        fmt.Println(error)
    }
    client := &Client{socket: connection}
    fmt.Printf("Connected to: %s\n", host)

    go client.receive()
    client.send()
}

func (client *Client) receive() {
    for {
        message := make([]byte, 4096)
        length, err := client.socket.Read(message)
        if err != nil {
            fmt.Printf("Error, terminating client - ERROR: %+v\n", err)
            client.socket.Close()
            break
        }
        if length > 0 {
            stringMessage := string(bytes.Trim(message, "\x00"))

            if stringMessage == CLOSE_CONNECTION {
                fmt.Println("Connection closed, terminating client")
                client.socket.Close()
                os.Exit(0)
            }

            fmt.Println(stringMessage)
        }
    }
}

func (client *Client) send() {
    for {
        reader := bufio.NewReader(os.Stdin)
        message, _ := reader.ReadString('\n')
        client.socket.Write([]byte(strings.TrimRight(message, "\n")))
    }
}
