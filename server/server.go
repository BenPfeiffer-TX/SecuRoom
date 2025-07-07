package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"
)

type message struct {
	Type string //supported types so far: echo, message
	//	Sender    conn.Addr
	Name      string
	Content   string
	Timestamp time.Time
}

var chatChannel chan message
var chatLog []message

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

	chatChannel = make(chan message)

	server, err := net.Listen("tcp", ":8080")
	defer server.Close()
	if err != nil {
		panic(err)
	}

	//separate thread for printing out message channel

	go func() {
		for msg := range chatChannel {
			chatLog = append(chatLog, msg)
			fmt.Println(msg.Content, msg.Timestamp.Format(time.TimeOnly))
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
	var receivedMessage message
	//we received a connection, now we handle it
	/*
		todo:
		check banned IPs?
		establish encryption key?
		send contents of chat channel
	*/
	fmt.Println("connection received from: ", conn.LocalAddr())
	//spawn separate thread for relaying any received messages to open connections
	go relayMessages(conn)

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
		//unmarshal contents of received message from client
		//content := string(buffer[:l])
		if err = json.Unmarshal(buffer[:l], &receivedMessage); err != nil {
			//issue decoding JSON
			log.Println("error unmarshalling:", err.Error())
			continue
		}
		//message is now unmarshaled

		/*
			response := message{Type: "echo",
				//	Sender:    receivedMessage.Sender,
				Name:      receivedMessage.Name,
				Content:   receivedMessage.Content,
				Timestamp: time.Now()}
			jsonResponse, err := json.Marshal(response)
			if err != nil {
				//issue marshalling response
				log.Println(err.Error())
				continue
			}
			conn.Write(jsonResponse)
		*/

		//switch on message type to determine action
		switch receivedMessage.Type {
		case "message":
			chatChannel <- receivedMessage
		default:
			log.Println("received unknown message of type: ", receivedMessage.Type)
		}
	}
}

func relayMessages(conn net.Conn) {
	i := len(chatLog) //length of chat at the time of connecting
	for {
		if len(chatLog) > i {
			//there are new messages in the chat
			msg, _ := json.Marshal(chatLog[i])
			conn.Write(msg)
			i++
		}
	}
}
