package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	//server listrning on 6379
	net, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		// handle error
		fmt.Printf("internal error occured")
		return
	}
	defer net.Close()
	for {
		//server accepting conncetions
		conn, _ := net.Accept()
		//welcoming user
		conn.Write([]byte("welcome to the serve\n"))
		//reading user input wrapped in buffer on connection ? can i use anything else
		scanner := bufio.NewScanner(conn)
		//while input exists , we keep readings
		for scanner.Scan() {
			line := scanner.Text() // Get the current input as string
			// Process the received line and print on our app server
			fmt.Println("Received:", line)
			//write back to user
			conn.Write([]byte(line + " to you to sir" + "\n"))
		}

	}

}
