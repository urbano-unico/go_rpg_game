package server

import (
    "bytes"
    "fmt"
    "math/rand"
    "net"
    "time"
)

const HEALTH = 10
const MAX_CLIENTS = 2

type ClientManager struct {
    clients    []*Client
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client

    inGame bool
    turn   *Client
}

type Client struct {
    socket net.Conn
    data   chan []byte
    player Player
}

type Player struct {
    force    int
    intelect int
    health   int
}

func newPlayer() Player {
    return Player{
        force:    randAttr(),
        intelect: randAttr(),
        health:   HEALTH,
    }
}

func randAttr() (attr int) {
    rand.Seed(time.Now().UnixNano())
    for attr == 0 {
        attr = rand.Intn(6)
    }
    return
}

func (manager *ClientManager) start() {
    go manager.handleBroadcast()
    for {
        select {
        case connection := <-manager.register:
            manager.registerNewClient(connection)
        case connection := <-manager.unregister:
            manager.unregisterClient(connection)
        }
    }
}

func (manager *ClientManager) handleBroadcast() {
    for {
        message := <-manager.broadcast
        for _, client := range manager.clients {
            client.data <- message
        }
    }
}

func (manager *ClientManager) registerNewClient(client *Client) {
    // Handle only two connections
    if len(manager.clients) == MAX_CLIENTS {
        fmt.Println("Recusing new connection!")
        client.data <- ERROR_CONNECTION_REFUSED
        close(client.data)
        return
    }

    manager.clients = append(manager.clients, client)
    go manager.receiver(client)
    go manager.sender(client)
    fmt.Printf("Added new Client%+v\n", client)

    if len(manager.clients) < MAX_CLIENTS {
        fmt.Println(string(WAITING_PLAYER))
        client.data <- WAITING_PLAYER
        return
    }

    fmt.Println("Two clients connected")
    manager.startGame()
}

func (manager *ClientManager) unregisterClient(client *Client) {
    clients := []*Client{}
    for _, managerClient := range manager.clients {
        if client != managerClient {
            clients = append(clients, managerClient)
            continue
        }

        fmt.Println("A connection has terminated! - Client", client)
        client.data <- CLOSE_CONNECTION
        manager.inGame = false
    }
    manager.clients = clients
}

func (manager *ClientManager) receiver(client *Client) {
    for {
        message := make([]byte, 2048)
        length, err := client.socket.Read(message)
        if err != nil {
            manager.unregister <- client
            break
        }
        if length > 0 {
            stringMessage := string(bytes.Trim(message, "\x00"))
            fmt.Printf("MESSAGE RECEIVED FROM Player%+v: %s\n", client.player, stringMessage)

            if stringMessage == "quit" {
                manager.unregister <- client
                continue
            }

            if !manager.inGame {
                client.data <- ERROR_NOT_IN_GAME
                continue
            }
            if client != manager.turn {
                client.data <- ERROR_NOT_YOUR_TURN
                continue
            }

            switch stringMessage {
            case "attack", "magic":
                manager.processAction(client, stringMessage)
                if manager.inGame {
                    manager.setClientTurn()
                    manager.sendTurn()
                } else {
                    manager.startGame()
                }
            default:
                fmt.Println("WRONG COMMAND: " + stringMessage)
            }
        }
    }
}

func (manager *ClientManager) sender(client *Client) {
    for {
        message, ok := <-client.data
        if !ok {
            return
        }
        client.socket.Write(message)
    }
}

func (manager *ClientManager) startGame() {
    manager.inGame = true

    manager.broadcast <- INITIALIZING_GAME
    for _, client := range manager.clients {
        client.player = newPlayer()
    }
    manager.setClientTurn()
    manager.sendTurn()

    fmt.Println("Starting the game!")
}

func (manager *ClientManager) endGame() {
    manager.inGame = false
    manager.turn = nil

    for _, client := range manager.clients {
        if client.player.health > 0 {
            client.data <- WINNER
        } else {
            client.data <- LOSER
        }
    }
    manager.sendTurn()
    manager.broadcast <- ENDING_GAME
    fmt.Println("Game ended!")
}

func (manager *ClientManager) processAction(clientTurn *Client, action string) {
    for _, client := range manager.clients {
        if client != manager.turn {
            switch action {
            case "attack":
                client.player.health = client.player.health - clientTurn.player.force
            case "magic":
                client.player.health = client.player.health - clientTurn.player.intelect
            }

            if client.player.health <= 0 {
                client.player.health = 0
                manager.endGame()
                break
            }
        }
    }
}

func (manager *ClientManager) setClientTurn() {
    rand.Seed(time.Now().UnixNano())
    selectedIndex := rand.Intn(len(manager.clients))

    if manager.inGame && manager.turn != manager.clients[selectedIndex] {
        manager.turn = manager.clients[selectedIndex]
        return
    }
    manager.setClientTurn()
}

func (manager *ClientManager) sendTurn() {
    for _, clientToSend := range manager.clients {
        message := ""
        for _, client := range manager.clients {
            playerMessage := "Player"
            if clientToSend == client {
                playerMessage = "[YOU] Player"
            }
            if manager.turn == client {
                playerMessage = fmt.Sprintf("* %s", playerMessage)
            }
            message = fmt.Sprintf("%s%s%+v\n", message, playerMessage, client.player)
        }
        clientToSend.data <- []byte(fmt.Sprintf("\n%s\n%s%s\n", STATS_HEADER, message, STATS_HEADER))

        turnMessage := CLIENT_WAITING
        if manager.inGame {
            if manager.turn == clientToSend {
                turnMessage = CLIENT_TURN
            }
            clientToSend.data <- turnMessage
        }
    }
}
