package main

import (
    "./client"
    "./server"
    "flag"
)

var protocol, host, mode string
var port int

func main() {
    flag.StringVar(&protocol, "protocol", "localhost", "protocol")
    flag.StringVar(&host, "host", "localhost", "host")
    flag.IntVar(&port, "port", 12345, "port")
    flag.StringVar(&mode, "mode", "client", "game mode")
    flag.Parse()

    if mode == "client" {
        client.StartClient()
    } else {
        server.StartServer()
    }
}
