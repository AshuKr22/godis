package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	//defer conn.close
	conn.Close()
	conn.Write([]byte("welcome to the serve\n"))
	//reading user input wrapped in buffer on connection ? can i use anything else
	scanner := bufio.NewScanner(conn)
	//while input exists , we keep readings
	for scanner.Scan() {
		line := scanner.Text() // Get the current input as string
		// Process the received line and print on our app server

		fmt.Println("Received:", line)
		line = strings.TrimSpace(line)
		if line == "exit" {
			conn.Close()
			fmt.Println("user terminated")
			//break the loop
			break
		}
		//write back to user
		conn.Write([]byte(line + " to you to sir" + "\n"))
	}
	//check for scanner error
	if scanner.Err() != nil {
		fmt.Print(scanner.Err())
	}

}
func main() {
	//server listrning on 6379
	net, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		// handle error
		fmt.Printf("internal error occurred")
		return
	}
	defer net.Close()
	for {
		//server accepting conncetions
		conn, _ := net.Accept()

		//added goroutine for multiple connections
		go handleConnection(conn)

	}

}
