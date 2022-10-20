package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

func initConnection(server string, username string) {
	var wg sync.WaitGroup

	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connection successful on " + conn.RemoteAddr().String())

	// Get confirmation
	buffer := make([]byte, 1024)
	size, err := conn.Read(buffer)
	response := buffer[0:size]
	if string(response) == "KO" {
		fmt.Println("Max clients reached !")
		os.Exit(1)
	}

	// Clear buffer
	for i, _ := range buffer {
		buffer[i] = 0
	}

	// Send username
	conn.Write([]byte(username))

	// Get confirmation
	_, err = conn.Read(buffer)
	if string(buffer) == "OK" {
		fmt.Println("Username already used !")
		os.Exit(1)
	}

	wg.Add(2)

	// Write message
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print(username + ": ") // PROMPT
			scanner.Scan()

			content := username + ": " + scanner.Text()
			//fmt.Println(string(content)) //DEBUG
			conn.Write([]byte(content))
		}
	}()

	// Receive message
	go func() {
		defer wg.Done()
		buffer := make([]byte, 1024)
		for {
			_, err := conn.Read(buffer)
			if err != nil {
				fmt.Println("Server disconnected !")
				os.Exit(1)
			}

			// Print message and the prompt
			fmt.Print("\r" + string(buffer))
			fmt.Print("\n" + username + ": ")

			// Clear buffer
			for i, _ := range buffer {
				buffer[i] = 0
			}
		}
	}()

	wg.Wait()

}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Helper : 127.0.0.1:8081 username")
		os.Exit(1)
	}

	server := os.Args[1]
	username := os.Args[2]

	if server == "" || !strings.Contains(server, ":") {
		fmt.Println("Helper : 127.0.0.1:8081")
		os.Exit(1)
	} else {
		initConnection(server, username)
	}
}
