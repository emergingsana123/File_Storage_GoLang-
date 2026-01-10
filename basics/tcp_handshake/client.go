package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "strings"
)

func main() {
	//connects to server running on 9000 , dial turns TCP connection object conn
    conn, err := net.Dial("tcp", "localhost:9000")
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    reader := bufio.NewReader(conn)

    //Receive server greeting
    greeting, _ := reader.ReadString('\n')
    fmt.Print(greeting)

    // Send username to server over tcp
    fmt.Print("Enter username: ")
    input := bufio.NewReader(os.Stdin)
    username, _ := input.ReadString('\n')
    username = strings.TrimSpace(username)
    conn.Write([]byte(username + "\n"))

    // Receive handshake confirmation and prints it
    confirmation, _ := reader.ReadString('\n')
    fmt.Print(confirmation)

    //send chat messages/looping of message
    for {
        fmt.Print("Message: ")
        msg, _ := input.ReadString('\n')
        conn.Write([]byte(msg))
    }
}
