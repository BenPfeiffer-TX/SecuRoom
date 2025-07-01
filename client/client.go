package main

import (
	"fmt"
	"net"

	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

func sendMessage(input *widget.Entry, conn net.Conn) {
	//input: pointer to widget.Entry, current connection
	text := input.Text
	_, err := conn.Write([]byte(text))
	if err != nil {
		fmt.Println(err)
		return
	}
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

	//fyne app window stuff
	a := app.New()
	windowName := "SecuRoom chat client"
	w := a.NewWindow(windowName)
	message := widget.NewLabel(windowName) //change this to be IP of connected server?
	input := widget.NewEntry()
	input.SetPlaceHolder("type here")
	send := widget.NewButton("send", func() { sendMessage(input, conn) })
	content := container.NewVBox(message, input, send)
	w.SetContent(content)
	w.ShowAndRun()

	tidyUp(conn)
}

func tidyUp(conn net.Conn) {
	fmt.Println("cleaning up after closing application")
	conn.Close()
}
