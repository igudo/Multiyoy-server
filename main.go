package main

import (
	"fmt"
	"net"
	"os"
)

// Const values. Todo: get this from arguments
const (
	host       = "localhost"
	port       = "7456"
	numClients = 2
)

// Global list of clients
var allClients [numClients]*Client

// Client is a player class
type Client struct {
	connection net.Conn // TCP
}

// Handles
func (c *Client) handleConnection() {
	buf := make([]byte, 2048)

	c.connection.Write([]byte("Write message: "))
	_, err := c.connection.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println(string(buf))
	c.connection.Write([]byte("Message received."))
	c.connection.Close()
	fmt.Println("Disconnected", c.connection.RemoteAddr().String())
}

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
	for i := 0; i < numClients; i++ {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		allClients[i] = &Client{connection: conn}
		fmt.Println("Connected", conn.RemoteAddr().String())

		// Handle connection in a new goroutine
		go allClients[i].handleConnection()
	}

	for {

	}
}
