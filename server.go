package main

import (
	"Multiyoy/players"
	"fmt"
	"net"
	"os"
	"time"
)

// Const values. Todo: get this from arguments
const (
	host = "localhost"
	port = "7456"
)

// GameData - Game variables for sending to players
var GameData = map[string]interface{}{
	"num_players":     8,
	"num_players_now": 0,
	"status":          "waiting",
}

// Players game data synchronization with the server's game data
var allPlayers [8]*players.Player
var numPlayersNow = 0
var gameStatus = "waiting"

func playersSync() {
	needSendGameData := 0
	somebodyDisconnected := false

	// While status == waiting lets check connections and disconnections
	for GameData["status"] == gameStatus {
		// Check if new connection
		if numPlayersNow > GameData["num_players_now"].(int) {
			fmt.Println("New connection! Sync...")
			GameData["num_players_now"] = numPlayersNow
			needSendGameData++
		}

		// Check if somebody offline
		for i, p := range allPlayers {
			if p != nil {
				if !p.Online {
					fmt.Println("Somebody disconnected! Sync...")
					numPlayersNow--
					GameData["num_players_now"] = numPlayersNow
					allPlayers[i] = nil
					somebodyDisconnected = true
				} else if needSendGameData > 0 && !p.NeedSendGameData {
					p.NeedSendGameData = true
				}
			}
		}

		if needSendGameData > 0 {
			needSendGameData--
		}
		if somebodyDisconnected {
			needSendGameData++
			somebodyDisconnected = false
		}
	}

	GameData["status"] = gameStatus
	// There are starting status so let's send this to players
	for _, p := range allPlayers {
		if p != nil {
			p.NeedSendGameData = true
		}
	}
}

func main() {
	numPlayers := 3 // Todo: get this from args

	// Listen for incoming connections
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener on application close
	defer l.Close()

	fmt.Println("Listening on " + host + ":" + port)

	GameData["num_players"] = numPlayers
	go playersSync()

	// Wait for every client to connect
	for numPlayersNow < numPlayers {
		for ; numPlayersNow < numPlayers; numPlayersNow++ {
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				os.Exit(1)
			}

			// Set new player to a free position
			for i, p := range allPlayers {
				if p == nil {
					allPlayers[i] = &players.Player{
						Connection:       conn,
						GameData:         &GameData,
						NeedSendGameData: false,
						Online:           true,
						Trusted:          false,
					}
					// Handle connection in a new goroutine
					go allPlayers[i].HandleConnection()
					break
				}
			}
			fmt.Println("Connected", conn.RemoteAddr().String())
		}

		// Wait until everyone is trusted
		if numPlayersNow == numPlayers {
			fmt.Println("Everybody connected! Waiting while trusted")
			for i, p := range allPlayers {
				for {
					if allPlayers[i] == nil || p.Trusted || numPlayersNow != numPlayers {
						break
					}
				}
			}
		}
	}

	// Here we starting
	gameStatus = "starting"
	time.Sleep(5 * time.Second)
	fmt.Println("WE ARE STARTING!!!!")
	time.Sleep(3 * time.Second)
}
