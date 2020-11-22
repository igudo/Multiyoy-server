package main

import (
	"Multiyoy/players"
	"fmt"
	"net"
	"os"
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

func playersSync() {
	numPlayersTemp := numPlayersNow
	needSendGameData := 0
	somebodyDisconnected := false
	for GameData["status"] == "waiting" {
		// Check if new connection
		if numPlayersNow > numPlayersTemp {
			fmt.Println("New connection! Sync...")
			GameData["num_players_now"] = numPlayersNow
			numPlayersTemp = numPlayersNow
			needSendGameData++
		}

		// Check if somebody offline
		for i, p := range allPlayers {
			if p != nil {
				if !p.Online {
					fmt.Println("Somebody disconnected! Sync...")
					numPlayersNow--
					numPlayersTemp--
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

	// If there are starting status let's send this to players
	for _, p := range allPlayers {
		if p != nil {
			if needSendGameData > 0 && !p.NeedSendGameData {
				p.NeedSendGameData = true
			}
		}
	}
}

func main() {
	numPlayers := 5 // Todo: get this from args

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
				}
				// Handle connection in a new goroutine
				go allPlayers[i].HandleConnection()
				break
			}
		}
		fmt.Println("Connected", conn.RemoteAddr().String())
	}

	GameData["status"] = "starting"
	for i, p := range allPlayers {
		if p != nil {
			fmt.Println(i, "is online")
		}
	}
	fmt.Println("WE ARE STARTING!!!!")
	for {
	}
}
