package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

const (
	CONN_PORT      = ":8888"
	CONN_TYPE      = "tcp"
	MSG_DISCONNECT = "Disconnected from the server.\n"
)

var (
	clients    = make(map[net.Conn]bool)
	broadcast  = make(chan string)
	register   = make(chan net.Conn)
	unregister = make(chan net.Conn)
	numUsers   int
)

func main() {
	listener, err := net.Listen(CONN_TYPE, CONN_PORT)
	if err != nil {
		log.Println("Error: ", err)
		os.Exit(1)
	}
	defer listener.Close()
	log.Println("Listening on " + CONN_PORT)

	go broadcaster()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error: ", err)
			continue
		}

		clients[conn] = true
		register <- conn
		numUsers++

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer func() {
		unregister <- conn
		conn.Close()
		numUsers--
	}()

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error reading:", err)
			return
		}

		message = strings.TrimSpace(message)

		if message == "exit" {
			return
		} else if message == "number of users" {
			// Send the number of connected users to the client
			fmt.Fprintf(conn, "Number of connected users: %d\n", numUsers)
			continue
		}

		broadcast <- fmt.Sprintf("Client at %s says: %s", conn.RemoteAddr(), message)
	}
}

func broadcaster() {
	for {
		select {
		case message := <-broadcast:
			for client := range clients {
				_, err := client.Write([]byte(message + "\n"))
				if err != nil {
					log.Println("Error writing:", err)
					delete(clients, client)
					client.Close()
				}
			}
		case conn := <-register:
			log.Println("New client connected:", conn.RemoteAddr())
		case conn := <-unregister:
			log.Println("Client disconnected:", conn.RemoteAddr())
			delete(clients, conn)
		}
	}
}
