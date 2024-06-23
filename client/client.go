package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const (
	SERVER_ADDR = "localhost:8888"
)

func main() {
	conn, err := net.Dial("tcp", SERVER_ADDR)
	if err != nil {
		fmt.Println("Error connecting:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connected to server.")

	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if input == "/exit" {
			fmt.Println("Exiting...")
			fmt.Fprintf(conn, "%s\n", input)
			break
		} else if input == "/users" {
			fmt.Fprintf(conn, "%s\n", input)
		} else {
			fmt.Fprintf(conn, "%s\n", input)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error:", err)
	}
}
