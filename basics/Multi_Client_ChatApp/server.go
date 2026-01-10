package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"sync"
)

// Message struct for handshake and chat
type Message struct {
	Type     string `json:"type"`              // "handshake" or "chat"
	Username string `json:"username,omitempty"` // optional
	Content  string `json:"content,omitempty"`  // message content
}

var (
	clients     = make(map[net.Conn]string) // active clients
	clientsLock = sync.Mutex{}              // to protect clients map
)

// Broadcasts a message to all clients except sender
func broadcast(sender net.Conn, msg Message) {
	clientsLock.Lock()
	defer clientsLock.Unlock()
	data, _ := json.Marshal(msg)
	for conn := range clients {
		if conn != sender { // don't send back to sender
			conn.Write(append(data, '\n'))
		}
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		clientsLock.Lock()
		username := clients[conn]
		delete(clients, conn)
		clientsLock.Unlock()
		conn.Close()
		fmt.Printf("%s disconnected\n", username)
	}()

	reader := bufio.NewReader(conn)

	// Step 1: Send greeting
	greet := Message{Type: "handshake", Content: "WELCOME! Please send your username:"}
	data, _ := json.Marshal(greet)
	conn.Write(append(data, '\n'))

	// Step 2: Receive handshake (username)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	var msg Message
	err := json.Unmarshal([]byte(line), &msg)
	if err != nil || msg.Type != "handshake" || msg.Username == "" {
		fmt.Println("Invalid handshake, closing connection")
		return
	}
	username := msg.Username

	// Add client to map
	clientsLock.Lock()
	clients[conn] = username
	clientsLock.Unlock()
	fmt.Printf("%s connected!\n", username)

	// Step 3: Send handshake confirmation
	confirm := Message{Type: "handshake", Username: username, Content: "Handshake complete."}
	data, _ = json.Marshal(confirm)
	conn.Write(append(data, '\n'))

	// Step 4: Chat loop
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		var chatMsg Message
		err = json.Unmarshal([]byte(line), &chatMsg)
		if err != nil || chatMsg.Type != "chat" {
			fmt.Println("Invalid chat message received")
			continue
		}
		fmt.Printf("[%s]: %s\n", chatMsg.Username, chatMsg.Content)
		// Broadcast to other clients
		broadcast(conn, chatMsg)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server listening on port 9000...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleConnection(conn)
	}
}
