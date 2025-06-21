package main

import (
	"bufio"
	"fmt"
	"net"
	"slices"
	"strings"
	"sync"
)

var STORAGE sync.Map
var COMMANDS = []string{"SET", "GET", "DEL"}

func handleSet(key string, value string) {
	STORAGE.Store(key, value)
}
func handleGet(key string) (string, bool) {
	value, ok := STORAGE.Load(key)
	if !ok {
		return "", ok
	}
	return value.(string), ok

}
func handleDel(key string) (string, bool) {

	old, ok := STORAGE.LoadAndDelete(key)
	if !ok {
		return "", false
	}
	return old.(string), ok

}

func splitThree(s string) (a, b, c string) {
	parts := strings.SplitN(s, " ", 3) // split **at most** 2 times → 3 parts
	switch len(parts) {
	case 3:
		return parts[0], parts[1], parts[2]
	case 2: // missing a third part
		return parts[0], parts[1], ""
	case 1: // nothing to split
		return parts[0], "", ""
	default:
		return "", "", ""
	}
}

func parser(input string) (string, error) {
	//divide input into three parts
	command, key, value := splitThree(input)
	//identify command , if invalid return error
	if command == "" {
		return "ERROR", fmt.Errorf("command not found ❌")
	}
	if !slices.Contains(COMMANDS, command) {
		return "ERROR", fmt.Errorf("%s is not a valid command ❌ ", command)
	}
	switch command {
	case "SET":
		handleSet(key, value)
		return "OK ✅", nil
	case "GET":
		value, ok := handleGet(key)
		if !ok {
			return "value not found ❌", nil
		}
		return value, nil
	case "DEL":
		old, ok := handleDel(key)
		if !ok {
			return "key not found ❌", nil
		}
		return old, nil
	}
	return "function end", nil

	//identify key , if not found return error
	//identify value , if not found return error
}

func handleConnection(conn net.Conn) {
	//defer conn.close
	defer conn.Close()
	conn.Write([]byte("go-redis initiated 🚀\n"))
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
			fmt.Println("user session terminated 🔚")
			//break the loop
			break
		}
		output, err := parser(line)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		//write back to user
		conn.Write([]byte(output + "\n"))
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
