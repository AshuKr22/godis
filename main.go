package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
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
		return command + " is not a valid command ❌", fmt.Errorf("user entered an invalid command : %s ❌ ", command)
	}
	switch command {
	case "SET":
		//write commands to log before updating memory
		writeToLogs(input)

		handleSet(key, value)
		return "OK ✅", nil
	case "GET":
		value, ok := handleGet(key)
		if !ok {
			return "value not found ❌", nil
		}
		return value, nil
	case "DEL":
		//write commands to log before updating memory
		writeToLogs(input)
		old, ok := handleDel(key)
		if !ok {
			return "key not found ❌", nil
		}
		return old, nil
	}
	return "function end", nil
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
		}

		//write back to user
		conn.Write([]byte(output + "\n"))
	}
	//check for scanner error
	if scanner.Err() != nil {
		fmt.Print(scanner.Err())
	}

}
func skipCommand(command string) bool {
	// if empty command then skip
	if command == "" {
		return true
	}
	return false
	//later add more conditions to check
	//like end of line marker to make sure data/value is not malformed

}
func readAndExecuteCommands() error {
	file, err := os.Open("log.txt")

	if err != nil {
		return err
	}
	//close file
	defer file.Close()
	scan_file := bufio.NewScanner(file)
	//print log entries
	fmt.Println("📊 printing log entries ... ")
	for scan_file.Scan() {
		command := scan_file.Text()
		command = strings.TrimSpace(command)
		//check valid command & execute
		if !skipCommand(command) {
			parser(command)
		}
		fmt.Println(command + " executed successfully")

	}
	return nil

}
func writeToLogs(command string) error {
	file, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	//close file
	defer file.Close()
	buf_writer := bufio.NewWriter(file)
	//write all the writes together at once in file once function ends
	defer buf_writer.Flush()

	fmt.Fprintln(buf_writer, command)

	return nil

}
func main() {
	//read from the file and execute all commands before starting listening on tcp
	err := readAndExecuteCommands()
	if err != nil {
		fmt.Printf("error occurred while recovering data : %s", err)
	}

	//server listening on 6379
	listener, err := net.Listen("tcp", "localhost:6379")
	fmt.Println("⏳ session starting up")
	if err != nil {
		// handle error
		fmt.Printf("internal error occurred")
		return
	}
	defer listener.Close()
	for {
		//server accepting connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		//added goroutine for multiple connections
		go handleConnection(conn)

	}

}
