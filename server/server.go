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

// var messageChannel chan message
// for now will just handle strings until i implement proper message struct data sending
var messageChannel chan string

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

	messageChannel = make(chan string)
	server, err := net.Listen("tcp", ":8080")
	defer server.Close()
	if err != nil {
		panic(err)
	}

	//separate thread for printing out message channel
	go func() {
		for msg := range messageChannel {
			fmt.Println("Received: ", msg)
		}
	}()
	//main loop, handling connections
	for {
		conn, err := server.Accept()
		if err != nil {
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
	for {
		l, err := conn.Read(buffer)
		if err != nil {
			log.Println(err.Error())
			if err.Error() == "EOF" {
				break
			}
		}
		if l == 0 {
			//buffer is empty, we keep looping waiting for content
			continue
		}
		//the contents of buffer are not a message: tbd

		//the contents of buffer are a message:
		message := string(buffer[:l])
		messageChannel <- message
	}
}
