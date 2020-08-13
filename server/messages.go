package server

import (
    "fmt"
)

var (
    COMMANDS     = []string{"attack", "magic", "quit"}
    STATS_HEADER = "------------- Player Stats -------------"

    WAITING_PLAYER    = []byte("Waiting for another player")
    INITIALIZING_GAME = []byte("\n\nInitializing game")
    ENDING_GAME       = []byte("Game ended")

    WINNER = []byte("\n\n********** You win the game!!! **********\n")
    LOSER  = []byte("\n\n********** You lose the game!!! **********\n")

    CLIENT_WAITING = []byte("Wait for your turn")
    CLIENT_TURN    = []byte(fmt.Sprintf("Is your turn, please choose a command: %+v", COMMANDS))

    ERROR_NOT_IN_GAME        = []byte("[ERROR] Not in game yet")
    ERROR_NOT_YOUR_TURN      = []byte("[ERROR] Not your turn, please wait")
    ERROR_CONNECTION_REFUSED = []byte("[ERROR] No more new connections")

    CLIENT_EXIT      = []byte("Client exit")
    CLOSE_CONNECTION = []byte("CLOSECONNECTION")
)
