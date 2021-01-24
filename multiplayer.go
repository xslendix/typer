package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/ahmetalpbalkan/go-cursor"
)

var port_default int = 5600

func ConfigureMultiplayer() (port int, name string, err error) {
	for port == 0 {
		fmt.Printf("Port (%d): ", port_default)
		_, err = fmt.Scanf("%d", &port)
		if err != nil {
			return
		}
	}

	for name == "" {
		fmt.Print("Game name: ")
		_, err = fmt.Scanf("%s", &name)
		if err != nil {
			return
		}
	}

	return
}

func startServer(port int, name string) {
	fmt.Print(cursor.ClearEntireScreen(), cursor.MoveUpperLeft(1))
	fmt.Printf("Starting server on port %d with name \"%s\"\n", port, name)

	l, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Println("Server started!")
	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Printf("Client connected on ip %s\n", c.RemoteAddr().String())
		go handleConnection(c, name)
	}
}

func handleConnection(conn net.Conn, name string) {
	writer := bufio.NewWriter(conn)
	reader := bufio.NewReader(conn)

	writer.WriteString("n\n")
	buffer, err := reader.ReadBytes('\n')
	if err != nil {
		log.Print("Client %s left\n", conn.RemoteAddr().String())
		conn.Close()
		return
	}
	log.Printf("Client on ip %s has name \"%s\"",
		conn.RemoteAddr().String(), string(buffer[:len(buffer)-1]))
	handleConnection(conn, name)
}

func connect(ip string, port int) {}
