package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

type message struct {
	Name      string
	Content   string
	Timestamp time.Time
}

func main() {
	//initialize server
	/*
		todo:
		load config file for settings like
		encryption
		names / anonymous
		password
		persistent storage of chat log
		security preferences (rate limits etc)
	*/
	server, err := net.Listen("tcp", ":8080")
	defer server.Close()
	if err != nil {
		panic(err)
	}

	for {
		conn, err := server.Accept()
		if err != nil {
			//some issue with returning next connection to listener
			//TBD how to handle
			log.Println("failed to accept connection, ", err.Error())
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	//we received a connection, now we handle it
	/*
		todo:
		check banned IPs?
		take initialization package (name, etc)?
		do handshake / send response
		establish encryption key?
		send contents of chat channel
	*/
	fmt.Println("connection received from: ", conn.LocalAddr())

	buffer := make([]byte, 1024)
	l, err := conn.Read(buffer)
	if err != nil {
		log.Println(err.Error())
		return
	}
	fmt.Printf("received: %s\n", string(buffer[:l]))
}
