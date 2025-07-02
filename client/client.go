package main

import (
	"fmt"
	"net"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func sendMessage(input *widget.Entry, chatwindow binding.ExternalStringList, conn net.Conn) {
	//input: pointer to widget.Entry, current connection
	text := input.Text
	_, err := conn.Write([]byte(text))
	if err != nil {
		fmt.Println(err)
		return
	}
	chatwindow.Append(text)
	//blank out text field after sending
	input.SetText("")
	return
}

func main() {
	//TCP server stuff
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	chatlog := []string{}
	chatwindow := binding.BindStringList(&chatlog)

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
	//receive chat log from server
	tidyUp(conn)
}

func tidyUp(conn net.Conn) {
	fmt.Println("cleaning up after closing application")
	conn.Close()
}
