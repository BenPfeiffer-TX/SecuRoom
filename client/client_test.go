package main

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func TestBadHostnames(t *testing.T) {
	app.New()
	//w := a.NewWindow("")
	input := widget.NewEntry()
	input.SetPlaceHolder("")

	var tests = []struct {
		host string
		resp string
	}{
		{"127.0.0.1", "missing port"},
		{"127.0.0.1:", "missing port"},
		{"offlineserver:8080", "dial tcp: lookup offlineserver: no such host"},
	}

	for _, test := range tests {
		testname := fmt.Sprintf("%s", test.host)
		t.Run(testname, func(t *testing.T) {
			input.SetText(test.host)
			_, err := connectTCP(input)
			if err.Error() != test.resp {
				t.Errorf("got %s, expect %s", err.Error(), test.resp)
			}
		})
	}
}
