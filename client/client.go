package main

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"strings"
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

func sendMessage(input *widget.Entry, conn net.Conn) {
	text := input.Text
	sendingMessage := message{Type: "chat", Name: "anon", Content: text, Timestamp: time.Now()}
	sendingJSON, err := json.Marshal(sendingMessage)
	if err != nil {
		//json marshalling failed
		log.Println(err.Error())
		return
	}
	_, err = conn.Write(sendingJSON)
	if err != nil {
		log.Println(err.Error())
		return
	}
	//blank out text field after sending
	input.SetText("")
	return
}

func receiveConnection(chatwindow binding.ExternalStringList, conn net.Conn) {
	//function for receiving data from server
	var receivedMessage message
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
		//we unmarshal what the server sent us
		if err = json.Unmarshal(buffer[:l], &receivedMessage); err != nil {
			//issue decoding JSON
			log.Println("error unmarshalling:", err.Error())
			return
		}
		switch receivedMessage.Type {
		case "echo":
			//server is echoing our previous message
			chatwindow.Append(receivedMessage.Content)
		default:
			//server sent us some unexpected message type
			log.Println("received unknown message of type: ", receivedMessage.Type)
		}
	}
}

func runtime() {
	chatlog := []string{}
	chatwindow := binding.BindStringList(&chatlog)
	/*
	**
	 */
	//fyne app window stuff
	a := app.New()
	windowName := "SecuRoom chat client"
	w := a.NewWindow(windowName)

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
	var conn net.Conn
	var err error

	input := widget.NewEntry()
	input.SetPlaceHolder("type here")
	send := widget.NewButton("send", func() { sendMessage(input, conn) })

	connect := widget.NewButton("connect", func() {
		conn, err = connectTCP(input)
		if err != nil {
			//connection failed, handle it
			log.Println(err.Error())
			return
		}
		message := widget.NewLabel(conn.RemoteAddr().String())
		contentConnected := container.NewBorder(message, input, nil, send, chat)
		go receiveConnection(chatwindow, conn)

		//connection succeeded, now we display contentConnected
		w.SetContent(contentConnected)
		w.Show()
	})

	contentStart := container.NewVBox(input, connect)
	w.SetContent(contentStart)
	w.Show()
	a.Run()
	defer conn.Close()
}

func connectTCP(input *widget.Entry) (net.Conn, error) {
	//TCP server stuff
	inputs := strings.Split(input.Text, ":")
	switch {
	//check inputs for errors
	case len(inputs) == 1:
		return nil, errors.New("missing port")
	case inputs[1] == "":
		return nil, errors.New("missing port")
	default:
		//input doesnt have any obvious errors, attempting connection
		conn, err := net.Dial("tcp", input.Text)
		if err != nil {
			//connecting failed
			log.Println(err.Error())
			return nil, err
		}
		return conn, nil
	}
}

func main() {

	runtime()
	//todo: key detection for pressing enter
	//make a toolbar for connecting / disconnecting
	tidyUp()
}

func tidyUp() {
	log.Println("cleaning up after closing application")
}
