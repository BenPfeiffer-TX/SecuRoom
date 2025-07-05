package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type message struct {
	Type string
	//	Sender	conn.Addr
	Name      string
	Content   string
	Timestamp time.Time
}

func sendMessage(input *widget.Entry, chatwindow binding.ExternalStringList, conn net.Conn) {
	text := input.Text
	sendingMessage := message{Type: "chat", Name: "anon", Content: text, Timestamp: time.Now()}
	sendingJSON, err := json.Marshal(sendingMessage)
	if err != nil {
		//json marshalling failed
		log.Println(err.Error())
		return
	}
	_, err = conn.Write(sendingJSON) //[]byte(text))
	if err != nil {
		log.Println(err.Error())
		return
	}
	//for testing: we write our sent message to the chatwindow
	chatwindow.Append(text)
	//blank out text field after sending
	input.SetText("")
	return
}

func receiveConnection(conn net.Conn) {
	//function for receiving data from server
	buffer := make([]byte, 1024)
	for {
		l, err := conn.Read(buffer)
		if err != nil {
			log.Println(err.Error())
			if err.Error() == "EOF" {
				break
				//server closed connection, handle it
			}
		}
		if l == 0 {
			continue
		}
		//	content := string(buffer[:l])
	}
}

func main() {
	//TCP server stuff
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	go receiveConnection(conn)

	chatlog := []string{}
	chatwindow := binding.BindStringList(&chatlog)
	/*
	**
	 */
	//fyne app window stuff
	a := app.New()
	windowName := "SecuRoom chat client"
	w := a.NewWindow(windowName)
	message := widget.NewLabel(windowName) //change this to be IP of connected server?
	input := widget.NewEntry()

	chat := widget.NewListWithData(chatwindow,
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i binding.DataItem, o fyne.CanvasObject) { o.(*widget.Label).Bind(i.(binding.String)) },
	)
	//constantly call chat.UnselectAll() in a goroutine so that no entries in the list can be highlighted
	//probably a better way to handle this
	go func() {
		for {
			fyne.Do(func() { chat.UnselectAll() })
		}
	}()

	input.SetPlaceHolder("type here")
	send := widget.NewButton("send", func() { sendMessage(input, chatwindow, conn) })
	content := container.NewBorder(message, input, nil, send, chat)
	w.SetContent(content)
	w.ShowAndRun()
	//todo: key detection for pressing enter
	//make a toolbar for connecting / disconnecting
	//figure out container layout changes
	tidyUp(conn)
}

func tidyUp(conn net.Conn) {
	fmt.Println("cleaning up after closing application")
	conn.Close()
}
