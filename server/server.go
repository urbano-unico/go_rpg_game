package server

import (
    "fmt"
    "net"
    "os"
    "os/signal"
    "syscall"
)

func signalsListener(listener *net.Listener, manager *ClientManager) {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
    _ = <-sigs

    fmt.Println("Terminating server, closing connections...")
    for _, client := range manager.clients {
        manager.unregister <- client
    }
    os.Exit(0)
}

func StartServer() {
    fmt.Println("Starting server...")
    listener, err := net.Listen("tcp", ":12345")
    if err != nil {
        fmt.Printf("Failed to listen tcp connection: %+v\n", err)
    }

    // Start the manager
    manager := ClientManager{
        clients:    []*Client{},
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
    go signalsListener(&listener, &manager)
    go manager.start()

    // Listen for connections
    for {
        connection, err := listener.Accept()
        if err != nil {
            fmt.Printf("[ERROR] Failed to accept connection: %+v\n", err)
            continue
        }
        client := &Client{socket: connection, data: make(chan []byte)}

        // Register new connection
        manager.register <- client
    }
}
