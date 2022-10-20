package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

// handle message from clients
func handleMessage(conn net.Conn, username string, clients map[string]net.Conn) {
	buffer := make([]byte, 1024)

	for {
		size, err := conn.Read(buffer)
		if err != nil {
			// Disconnect and delete client
			msg := fmt.Sprintf("Client %s disconnected !", username)
			for client, c := range clients {
				if client == username {
					continue
				}
				c.Write([]byte(msg))
			}
			delete(clients, username)
			break
		}

		//fmt.Println(string(buffer)) // DEBUG
		message := buffer[0:size]

		// Forward message
		for client, c := range clients {
			if client == username {
				continue
			}
			c.Write(message)
		}

		// Clear buffer
		for i, _ := range buffer {
			buffer[i] = 0
		}
	}
}

// accept new connection from listener
func acceptConnection(listener net.Listener, clients map[string]net.Conn, maxClients int) {
	conn, err := listener.Accept()
	if err != nil {
		log.Print(err)
	}
	log.Println("New client: " + conn.RemoteAddr().String())

	// Check if max clients reached
	if len(clients) >= maxClients {
		conn.Write([]byte("KO"))
		time.Sleep(5 * time.Second)
		conn.Close()
		log.Println("Max clients reached !")
		return
	}
	// Send confirmation to client
	conn.Write([]byte("OK"))

	// Get client username
	buffer := make([]byte, 1024)
	size, err := conn.Read(buffer)
	if err != nil {
		log.Print(err)
	}

	// Add client
	usernameBytes := buffer[0:size]
	username := string(usernameBytes)
	clients[username] = conn

	// Send confirmation to client
	conn.Write([]byte("OK"))

	// Announce new client to every client
	msg := username + " joined the chat !"
	for client, c := range clients {
		if client == username {
			continue
		}
		c.Write([]byte(msg))
	}

	go func() {
		handleMessage(conn, username, clients)
	}()
}

func initServer(host string, maxClients int) {
	// Initiate listener
	fmt.Println("Starting listening on " + host)
	listener, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatal(err)
	}

	// List of clients
	clients := make(map[string]net.Conn)

	// Handle incoming connections
	for {
		acceptConnection(listener, clients, maxClients)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Example usage : 127.0.0.1:8081 2")
		os.Exit(1)
	}

	host := os.Args[1]
	maxClients, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	initServer(host, maxClients)
}
