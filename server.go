package main

import (
	"Multiyoy/players"
	"fmt"
	"net"
	"os"
)

// Const values. Todo: get this from arguments
const (
	host       = "localhost"
	port       = "7456"
	numPlayers = 2
)

// Global list of clients
var allPlayers [numPlayers]*players.Player

var gameStatus = "initializing"

func main() {
	// Listen for incoming connections
	l, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	// Close the listener when the application closes
	defer l.Close()

	fmt.Println("Listening on " + host + ":" + port)

	// Wait for every client to connect
	gameStatus = "waiting"
	for i := 0; i < numPlayers; i++ {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		allPlayers[i] = &players.Player{Connection: conn}
		fmt.Println("Connected", conn.RemoteAddr().String())

		// Handle connection in a new goroutine
		go allPlayers[i].HandleConnection()
	}

	for {

	}
}
