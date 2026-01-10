// package tcphandshake
package main

import (
    "bufio"
    "fmt"
    "net"
    "strings"
)
//send a greeting to the client asking for a username 
func handleConnection(conn net.Conn) {

    defer conn.Close()
    
    // writes bytes to the TCP connection 
    conn.Write([]byte("WELCOME! Please send your username:\n"))

    // Step 2: Read username from client
    reader := bufio.NewReader(conn) //allows reading lines 
    username, _ := reader.ReadString('\n') //reads until client presses enter
    username = strings.TrimSpace(username) 

    // Step 3: Acknowledge handshake/ sens back a handshake confirmation
    conn.Write([]byte(fmt.Sprintf("Hello %s! Handshake complete.\n", username)))

    // Keep connection alive for chat/ prints each message to server console
    for {
        msg, err := reader.ReadString('\n')
        if err != nil { // err triggers when client closes the connection 
            fmt.Printf("%s disconnected.\n", username)
            break
        }
        fmt.Printf("[%s]: %s", username, msg)
    }
}

func main() {
	//this line starts the port of 9000  
    listener, err := net.Listen("tcp", ":9000")
    if err != nil {
        panic(err)
    }
    fmt.Println("Server listening on port 9000...")
    // listener waits for the incoming client
    for {
		//listener = blocks until a client connects, return connection object conn
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Connection error:", err)
            continue
        }
        go handleConnection(conn)
		//lightweigth threas for each client, spawns a new goroutine
    }
}
