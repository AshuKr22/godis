package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

var DEFAULT_EXPIRATION_TIME time.Duration = time.Hour * 6

type TimedValue struct {
	expiry     time.Duration
	value      string
	time_stamp time.Time
}

func NewTimedValue(expiry time.Duration, value string) *TimedValue {
	return &TimedValue{expiry: expiry, value: value, time_stamp: time.Now()}
}

var STORAGE sync.Map
var COMMANDS = []string{"SET", "GET", "DEL", "SETEX"}

func handleSet(key string, value string) {
	timed_value := NewTimedValue(DEFAULT_EXPIRATION_TIME, value)
	STORAGE.Store(key, timed_value)
}
func handleGet(key string) (string, bool) {
	timed_value, ok := STORAGE.Load(key)
	if !ok {
		return "", ok
	}
	//check if timed value
	tv, ok := timed_value.(*TimedValue)
	if !ok {
		return "", ok
	}
	//subtract current time from time_stamp om object
	ttl := time.Since(tv.time_stamp)
	fmt.Printf("time since inception : %s", ttl)
	if ttl > tv.expiry {
		return "", false

	}
	//return actual value
	return tv.value, ok

}
func handleSetEx(key string, value string, timeInSeconds string) {
	//convert time string to seconds
	timeInt, err := strconv.Atoi(timeInSeconds)
	if err != nil {
		fmt.Println("error converting string to int")
		return
	}
	time := time.Duration(timeInt) * time.Second
	timed_value := NewTimedValue(time, value)
	STORAGE.Store(key, timed_value)

}
func handleDel(key string) (string, bool) {

	oldTimedValue, ok := STORAGE.LoadAndDelete(key)

	if !ok {
		return "", false
	}
	old, ok := oldTimedValue.(*TimedValue)
	if !ok {
		return "", false
	}

	return old.value, ok

}

func splitThree(s string) (a, b, c, d string) {
	parts := strings.SplitN(s, " ", 4) // split **at most** 2 times → 3 parts
	fmt.Println(parts)
	switch len(parts) {
	case 4:
		return parts[0], parts[1], parts[2], parts[3]
	case 3:
		return parts[0], parts[1], parts[2], ""
	case 2:
		return parts[0], parts[1], "", ""
	case 1:
		return parts[0], "", "", ""
	default:
		return "", "", "", ""
	}
}

func parser(input string, is_recovering bool) (string, error) {
	//divide input into three parts
	command, key, value, time := splitThree(input)
	//identify command , if invalid return error
	if command == "" {
		return "ERROR", fmt.Errorf("command not found ❌")
	}
	if !slices.Contains(COMMANDS, command) {
		return command + " is not a valid command ❌", fmt.Errorf("user entered an invalid command : %s ❌ ", command)
	}
	switch command {
	case "SET":
		if !is_recovering {
			writeToLogs(input)
		}
		handleSet(key, value)
		return "OK ✅", nil

	case "SETEX":
		if !is_recovering {
			writeToLogs(input)
		}
		handleSetEx(key, value, time)
		return "OK ✅", nil
	case "GET":
		value, ok := handleGet(key)
		if !ok {
			return "value not found or expired ❌", nil
		}
		return value, nil
	case "DEL":
		if !is_recovering {
			writeToLogs(input)
		}
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
	conn.Write([]byte("go-redis initiated 🚀 from " + os.Getenv("INSTANCE_ID") + "\n"))
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
		output, err := parser(line, false)
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
			parser(command, true)
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
	instanceID := os.Getenv("INSTANCE_ID")
	fmt.Printf("🚀 starting instance %s\n", instanceID)
	//read from the file and execute all commands before starting listening on tcp
	err := readAndExecuteCommands()
	if err != nil {
		fmt.Printf("error occurred while recovering data : %s", err)
	}

	//server listening on 6379
	listener, err := net.Listen("tcp", ":6379")
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
