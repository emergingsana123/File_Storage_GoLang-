package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

// Message struct for handshake and chat
type Message struct {
	Type     string `json:"type"`
	Username string `json:"username,omitempty"`
	Content  string `json:"content,omitempty"`
}

func main() {
	conn, err := net.Dial("tcp", "localhost:9000")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	input := bufio.NewReader(os.Stdin)

	// Step 1: Receive server greeting
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	var msg Message
	json.Unmarshal([]byte(line), &msg)
	fmt.Println(msg.Content)

	// Step 2: Send username
	fmt.Print("Enter username: ")
	username, _ := input.ReadString('\n')
	username = strings.TrimSpace(username)

	handshake := Message{Type: "handshake", Username: username}
	data, _ := json.Marshal(handshake)
	conn.Write(append(data, '\n'))

	// Step 3: Receive handshake confirmation
	line, _ = reader.ReadString('\n')
	line = strings.TrimSpace(line)
	json.Unmarshal([]byte(line), &msg)
	fmt.Println(msg.Content)

	// Step 4: Chat loop
	for {
		fmt.Print("Message: ")
		text, _ := input.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		chat := Message{Type: "chat", Username: username, Content: text}
		data, _ := json.Marshal(chat)
		conn.Write(append(data, '\n'))
	}
}
