package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/websocket"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter your username: ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Print("enter server websocket url (default ws://localhost:8080/ws): ")
	serverURL, _ := reader.ReadString('\n')
	serverURL = strings.TrimSpace(serverURL)
	if serverURL == "" {
		serverURL = "ws://localhost:8080/ws"
	}

	c, _, err := websocket.DefaultDialer.Dial(serverURL, nil)
	if err != nil {
		log.Fatal("Connection failed:", err)
	}
	defer c.Close()

	fmt.Printf("\n--- Connected to Chat Server as [%s] ---\n", username)

	// first goroutine that reads messages from the server and prints them to the console
	go func() {
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println("\n[Disconnected from server]")
				return
			}
			// print the message to console or smth
			fmt.Printf("\r%s\n> ", string(message))
		}
	}()

	// i think this is the second goroutine that reads user input and sends it to the server i didnt put comments while coding this
	for {
		fmt.Print("> ")
		text, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		// format message with username or smth
		fullMessage := fmt.Sprintf("[%s]: %s", username, text)

		// send over websocket
		err = c.WriteMessage(websocket.TextMessage, []byte(fullMessage))
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}
}
